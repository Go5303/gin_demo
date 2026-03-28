package excel

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// ============================================================
// Color constants
// ============================================================

const (
	ColorWhite   = "#FFFFFF"
	ColorBlack   = "#000000"
	ColorRed     = "#FF0000"
	ColorGreen   = "#00B050"
	ColorBlue    = "#4472C4"
	ColorYellow  = "#FFFF00"
	ColorOrange  = "#FFC000"
	ColorGray    = "#D9D9D9"
	ColorLightBg = "#F2F2F2"
	ColorDarkBg  = "#333F4F"
)

// ============================================================
// Border presets
// ============================================================

type BorderPreset int

const (
	BorderNone BorderPreset = iota
	BorderThin
	BorderMedium
	BorderThick
	BorderAll     // thin all sides
	BorderOutline // medium outline only
)

// ============================================================
// Alignment presets
// ============================================================

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

type VAlign int

const (
	VAlignTop VAlign = iota
	VAlignMiddle
	VAlignBottom
)

// ============================================================
// Style builder - chain style
// ============================================================

// Style is a fluent builder for cell styles
type Style struct {
	font      *excelize.Font
	fill      *excelize.Fill
	alignment *excelize.Alignment
	border    []excelize.Border
	numFmt    string
}

// NewStyle creates a new style builder
func NewStyle() *Style {
	return &Style{
		font:      &excelize.Font{Size: 11, Family: "Microsoft YaHei"},
		alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
	}
}

// Font sets font name
func (s *Style) Font(name string) *Style {
	s.font.Family = name
	return s
}

// FontSize sets font size
func (s *Style) FontSize(size float64) *Style {
	s.font.Size = size
	return s
}

// Bold sets bold
func (s *Style) Bold() *Style {
	s.font.Bold = true
	return s
}

// Italic sets italic
func (s *Style) Italic() *Style {
	s.font.Italic = true
	return s
}

// FontColor sets font color (hex like "#FF0000")
func (s *Style) FontColor(color string) *Style {
	s.font.Color = color
	return s
}

// BgColor sets solid fill background color
func (s *Style) BgColor(color string) *Style {
	s.fill = &excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{color}}
	return s
}

// HAlign sets horizontal alignment
func (s *Style) HAlign(a Align) *Style {
	switch a {
	case AlignLeft:
		s.alignment.Horizontal = "left"
	case AlignCenter:
		s.alignment.Horizontal = "center"
	case AlignRight:
		s.alignment.Horizontal = "right"
	}
	return s
}

// VAlign sets vertical alignment
func (s *Style) VAlign(a VAlign) *Style {
	switch a {
	case VAlignTop:
		s.alignment.Vertical = "top"
	case VAlignMiddle:
		s.alignment.Vertical = "center"
	case VAlignBottom:
		s.alignment.Vertical = "bottom"
	}
	return s
}

// Center sets both horizontal and vertical center
func (s *Style) Center() *Style {
	s.alignment.Horizontal = "center"
	s.alignment.Vertical = "center"
	return s
}

// WrapText enables/disables text wrapping
func (s *Style) WrapText(wrap bool) *Style {
	s.alignment.WrapText = wrap
	return s
}

// Border sets border preset
func (s *Style) Border(preset BorderPreset) *Style {
	sides := []string{"left", "top", "right", "bottom"}
	switch preset {
	case BorderNone:
		s.border = nil
	case BorderThin, BorderAll:
		for _, side := range sides {
			s.border = append(s.border, excelize.Border{Type: side, Color: ColorBlack, Style: 1})
		}
	case BorderMedium:
		for _, side := range sides {
			s.border = append(s.border, excelize.Border{Type: side, Color: ColorBlack, Style: 2})
		}
	case BorderThick:
		for _, side := range sides {
			s.border = append(s.border, excelize.Border{Type: side, Color: ColorBlack, Style: 5})
		}
	case BorderOutline:
		for _, side := range sides {
			s.border = append(s.border, excelize.Border{Type: side, Color: ColorBlack, Style: 2})
		}
	}
	return s
}

// BorderColor sets border with custom color
func (s *Style) BorderColor(preset BorderPreset, color string) *Style {
	s.Border(preset)
	for i := range s.border {
		s.border[i].Color = color
	}
	return s
}

// NumFmt sets number format (e.g. "#,##0.00", "0.00%", "yyyy-mm-dd")
func (s *Style) NumFmt(fmt string) *Style {
	s.numFmt = fmt
	return s
}

// build converts to excelize.Style
func (s *Style) build() *excelize.Style {
	st := &excelize.Style{
		Font:      s.font,
		Alignment: s.alignment,
		Border:    s.border,
	}
	if s.fill != nil {
		st.Fill = *s.fill
	}
	if s.numFmt != "" {
		st.CustomNumFmt = &s.numFmt
	}
	return st
}

// ============================================================
// Preset styles (common OA scenarios)
// ============================================================

// TitleStyle returns a big bold centered title style
func TitleStyle() *Style {
	return NewStyle().FontSize(16).Bold().Center().FontColor(ColorWhite).BgColor(ColorDarkBg)
}

// HeaderStyle returns a standard table header style
func HeaderStyle() *Style {
	return NewStyle().FontSize(11).Bold().Center().BgColor(ColorGray).Border(BorderAll)
}

// CellStyle returns a standard data cell style
func CellStyle() *Style {
	return NewStyle().FontSize(10).HAlign(AlignLeft).VAlign(VAlignMiddle).Border(BorderAll)
}

// CenterCellStyle returns a centered data cell style
func CenterCellStyle() *Style {
	return NewStyle().FontSize(10).Center().Border(BorderAll)
}

// AmountStyle returns a right-aligned amount cell style
func AmountStyle() *Style {
	return NewStyle().FontSize(10).HAlign(AlignRight).VAlign(VAlignMiddle).Border(BorderAll).NumFmt("#,##0.00")
}

// ============================================================
// Excel builder - the main wrapper
// ============================================================

// Builder wraps excelize.File with convenience methods
type Builder struct {
	f     *excelize.File
	sheet string
}

// New creates a new Excel builder with default sheet
func New() *Builder {
	f := excelize.NewFile()
	return &Builder{f: f, sheet: "Sheet1"}
}

// SetSheet switches to a named sheet, creates it if not exists
func (b *Builder) SetSheet(name string) *Builder {
	idx, err := b.f.GetSheetIndex(name)
	if err != nil || idx < 0 {
		b.f.NewSheet(name)
	}
	b.sheet = name
	return b
}

// RenameSheet renames the current sheet
func (b *Builder) RenameSheet(newName string) *Builder {
	b.f.SetSheetName(b.sheet, newName)
	b.sheet = newName
	return b
}

// ---- Cell operations ----

// Set writes a value to a cell
func (b *Builder) Set(cell string, value interface{}) *Builder {
	b.f.SetCellValue(b.sheet, cell, value)
	return b
}

// SetRC writes a value by row/col (1-based)
func (b *Builder) SetRC(row, col int, value interface{}) *Builder {
	cell, _ := excelize.CoordinatesToCellName(col, row)
	b.f.SetCellValue(b.sheet, cell, value)
	return b
}

// SetStyle applies a Style to a cell or range (e.g. "A1" or "A1:D1")
func (b *Builder) SetStyle(cellRange string, style *Style) *Builder {
	styleID, err := b.f.NewStyle(style.build())
	if err == nil {
		b.f.SetCellStyle(b.sheet, cellRange, cellRange, styleID)
	}
	return b
}

// SetRangeStyle applies a Style to a range "A1:D5"
func (b *Builder) SetRangeStyle(from, to string, style *Style) *Builder {
	styleID, err := b.f.NewStyle(style.build())
	if err == nil {
		b.f.SetCellStyle(b.sheet, from, to, styleID)
	}
	return b
}

// Merge merges cells (e.g. "A1", "D1")
func (b *Builder) Merge(from, to string) *Builder {
	b.f.MergeCell(b.sheet, from, to)
	return b
}

// MergeRC merges cells by row/col (1-based)
func (b *Builder) MergeRC(fromRow, fromCol, toRow, toCol int) *Builder {
	from, _ := excelize.CoordinatesToCellName(fromCol, fromRow)
	to, _ := excelize.CoordinatesToCellName(toCol, toRow)
	b.f.MergeCell(b.sheet, from, to)
	return b
}

// ---- Width / Height ----

// ColWidth sets column width (col like "A", "B", "C")
func (b *Builder) ColWidth(col string, width float64) *Builder {
	b.f.SetColWidth(b.sheet, col, col, width)
	return b
}

// ColWidthRange sets width for a range of columns ("A" to "F")
func (b *Builder) ColWidthRange(from, to string, width float64) *Builder {
	b.f.SetColWidth(b.sheet, from, to, width)
	return b
}

// RowHeight sets row height (1-based)
func (b *Builder) RowHeight(row int, height float64) *Builder {
	b.f.SetRowHeight(b.sheet, row, height)
	return b
}

// ---- Batch write helpers ----

// WriteRow writes a row of values starting from a cell (e.g. "A1")
func (b *Builder) WriteRow(startCell string, values ...interface{}) *Builder {
	col, row, _ := excelize.CellNameToCoordinates(startCell)
	for i, v := range values {
		cell, _ := excelize.CoordinatesToCellName(col+i, row)
		b.f.SetCellValue(b.sheet, cell, v)
	}
	return b
}

// WriteRows writes multiple rows starting from startCell
func (b *Builder) WriteRows(startCell string, rows [][]interface{}) *Builder {
	col, row, _ := excelize.CellNameToCoordinates(startCell)
	for r, rowData := range rows {
		for c, v := range rowData {
			cell, _ := excelize.CoordinatesToCellName(col+c, row+r)
			b.f.SetCellValue(b.sheet, cell, v)
		}
	}
	return b
}

// WriteTable writes headers + data rows with auto-styling
// Returns the Builder for chaining; the table starts at startCell
func (b *Builder) WriteTable(startCell string, headers []string, data [][]interface{}, opts *TableOptions) *Builder {
	if opts == nil {
		opts = DefaultTableOptions()
	}

	startCol, startRow, _ := excelize.CellNameToCoordinates(startCell)

	// -- Title row (optional) --
	dataStartRow := startRow
	if opts.Title != "" {
		from, _ := excelize.CoordinatesToCellName(startCol, startRow)
		to, _ := excelize.CoordinatesToCellName(startCol+len(headers)-1, startRow)
		b.Merge(from, to)
		b.Set(from, opts.Title)
		b.SetRangeStyle(from, to, opts.TitleStyle)
		b.RowHeight(startRow, 36)
		dataStartRow++
	}

	// -- Header row --
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(startCol+i, dataStartRow)
		b.Set(cell, h)
	}
	hFrom, _ := excelize.CoordinatesToCellName(startCol, dataStartRow)
	hTo, _ := excelize.CoordinatesToCellName(startCol+len(headers)-1, dataStartRow)
	b.SetRangeStyle(hFrom, hTo, opts.HeaderStyle)
	b.RowHeight(dataStartRow, 28)

	// -- Data rows --
	for r, row := range data {
		rowIdx := dataStartRow + 1 + r
		for c, v := range row {
			cell, _ := excelize.CoordinatesToCellName(startCol+c, rowIdx)
			b.Set(cell, v)
		}
		dFrom, _ := excelize.CoordinatesToCellName(startCol, rowIdx)
		dTo, _ := excelize.CoordinatesToCellName(startCol+len(headers)-1, rowIdx)
		b.SetRangeStyle(dFrom, dTo, opts.CellStyle)
	}

	// -- Auto column width --
	if opts.AutoWidth {
		for i, h := range headers {
			w := float64(len([]rune(h)))*2.5 + 4
			if w < 12 {
				w = 12
			}
			if w > 50 {
				w = 50
			}
			colName, _ := excelize.ColumnNumberToName(startCol + i)
			b.ColWidth(colName, w)
		}
	}

	return b
}

// ---- Table options ----

// TableOptions configures WriteTable behavior
type TableOptions struct {
	Title       string
	TitleStyle  *Style
	HeaderStyle *Style
	CellStyle   *Style
	AutoWidth   bool
}

// DefaultTableOptions returns sensible defaults
func DefaultTableOptions() *TableOptions {
	return &TableOptions{
		TitleStyle:  TitleStyle(),
		HeaderStyle: HeaderStyle(),
		CellStyle:   CellStyle(),
		AutoWidth:   true,
	}
}

// ---- Output ----

// SaveAs saves the Excel file to disk
func (b *Builder) SaveAs(path string) error {
	return b.f.SaveAs(path)
}

// WriteTo writes the Excel to an io.Writer
func (b *Builder) WriteTo(w io.Writer) error {
	return b.f.Write(w)
}

// Download sends the Excel as a download response in Gin
func (b *Builder) Download(c *gin.Context, filename string) {
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.QueryEscape(filename)))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	if err := b.f.Write(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "导出失败: %v", err)
	}
}

// File returns the underlying excelize.File for advanced operations
func (b *Builder) File() *excelize.File {
	return b.f
}

// Close closes the builder and releases resources
func (b *Builder) Close() error {
	return b.f.Close()
}
