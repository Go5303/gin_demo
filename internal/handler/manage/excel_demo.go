package manage

import (
	"fmt"
	"time"

	"github.com/Go5303/gin_demo/pkg/excel"
	"github.com/gin-gonic/gin"
)

// ExcelDemo demonstrates the excel builder with merge, styles, and table
func ExcelDemo(c *gin.Context) {
	eb := excel.New()
	defer eb.Close()

	eb.RenameSheet("员工考勤汇总")

	// ============ Title: merged across A1:H1 ============
	eb.Merge("A1", "H1").
		Set("A1", "我大OA - 2026年3月员工考勤汇总表").
		SetRangeStyle("A1", "H1", excel.TitleStyle()).
		RowHeight(1, 42)

	// ============ Sub info: merged rows ============
	eb.Merge("A2", "D2").Set("A2", "部门：技术部").
		SetRangeStyle("A2", "D2", excel.NewStyle().FontSize(11).Bold().HAlign(excel.AlignLeft).VAlign(excel.VAlignMiddle))
	eb.Merge("E2", "H2").Set("E2", fmt.Sprintf("导出时间：%s", time.Now().Format("2006-01-02 15:04:05"))).
		SetRangeStyle("E2", "H2", excel.NewStyle().FontSize(11).HAlign(excel.AlignRight).VAlign(excel.VAlignMiddle))
	eb.RowHeight(2, 28)

	// ============ Table header with two-level merge ============
	// First level header
	eb.Merge("A3", "A4").Set("A3", "序号")
	eb.Merge("B3", "B4").Set("B3", "姓名")
	eb.Merge("C3", "C4").Set("C3", "工号")
	eb.Merge("D3", "D4").Set("D3", "部门")
	eb.Merge("E3", "F3").Set("E3", "出勤情况")     // merged two cols
	eb.Merge("G3", "H3").Set("G3", "异常统计")      // merged two cols

	// Second level header
	eb.Set("E4", "应出勤(天)").Set("F4", "实出勤(天)")
	eb.Set("G4", "迟到(次)").Set("H4", "早退(次)")

	// Style all headers
	eb.SetRangeStyle("A3", "H4", excel.HeaderStyle())
	eb.RowHeight(3, 30).RowHeight(4, 26)

	// Column widths
	eb.ColWidth("A", 8).ColWidth("B", 12).ColWidth("C", 14).ColWidth("D", 16).
		ColWidthRange("E", "H", 14)

	// ============ Data rows ============
	employees := [][]interface{}{
		{1, "张三", "EMP001", "技术部", 22, 22, 0, 0},
		{2, "李四", "EMP002", "技术部", 22, 21, 1, 0},
		{3, "王五", "EMP003", "技术部", 22, 20, 2, 1},
		{4, "赵六", "EMP004", "技术部", 22, 22, 0, 0},
		{5, "孙七", "EMP005", "技术部", 22, 19, 3, 2},
	}

	normalCell := excel.CenterCellStyle()
	warnCell := excel.NewStyle().FontSize(10).Center().Border(excel.BorderAll).FontColor(excel.ColorRed).Bold()

	for i, row := range employees {
		r := 5 + i
		for j, v := range row {
			eb.SetRC(r, j+1, v)
		}
		// Style the whole row first
		from := fmt.Sprintf("A%d", r)
		to := fmt.Sprintf("H%d", r)
		eb.SetRangeStyle(from, to, normalCell)

		// Highlight anomalies in red
		late, _ := row[6].(int)
		early, _ := row[7].(int)
		if late > 0 {
			cell := fmt.Sprintf("G%d", r)
			eb.SetStyle(cell, warnCell)
		}
		if early > 0 {
			cell := fmt.Sprintf("H%d", r)
			eb.SetStyle(cell, warnCell)
		}
	}

	// ============ Summary row: merged ============
	sumRow := 5 + len(employees)
	eb.Merge(fmt.Sprintf("A%d", sumRow), fmt.Sprintf("D%d", sumRow)).
		SetRC(sumRow, 1, "合 计").
		SetRC(sumRow, 5, 110).
		SetRC(sumRow, 6, 104).
		SetRC(sumRow, 7, 6).
		SetRC(sumRow, 8, 3)

	sumStyle := excel.NewStyle().FontSize(11).Bold().Center().BgColor(excel.ColorLightBg).Border(excel.BorderAll)
	eb.SetRangeStyle(fmt.Sprintf("A%d", sumRow), fmt.Sprintf("H%d", sumRow), sumStyle)
	eb.RowHeight(sumRow, 28)

	// ============ Second sheet: simple table demo ============
	eb.SetSheet("请假统计")

	headers := []string{"序号", "姓名", "请假类型", "开始日期", "结束日期", "天数", "审批状态"}
	data := [][]interface{}{
		{1, "李四", "事假", "2026-03-05", "2026-03-05", 1, "已通过"},
		{2, "王五", "病假", "2026-03-10", "2026-03-11", 2, "已通过"},
		{3, "孙七", "年假", "2026-03-15", "2026-03-17", 3, "已通过"},
	}

	eb.WriteTable("A1", headers, data, &excel.TableOptions{
		Title:       "2026年3月请假统计表",
		TitleStyle:  excel.TitleStyle(),
		HeaderStyle: excel.HeaderStyle(),
		CellStyle:   excel.CenterCellStyle(),
		AutoWidth:   true,
	})

	// Download
	eb.Download(c, "考勤汇总_2026年3月.xlsx")
}
