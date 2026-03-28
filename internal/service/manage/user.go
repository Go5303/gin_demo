package manage

import (
	"encoding/json"
	"fmt"
	"time"

	"git.woda.ink/Woda_OA/config"
	"git.woda.ink/Woda_OA/internal/model"
	"git.woda.ink/Woda_OA/pkg/cache"
	"git.woda.ink/Woda_OA/pkg/crypto"
	"git.woda.ink/Woda_OA/pkg/errcode"
)

// LoginReq is the login request body
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResp is the login response data
type LoginResp struct {
	Token string `json:"token"`
	Admin any    `json:"admin"`
}

// Login handles admin login logic
func Login(req *LoginReq) (*LoginResp, error) {
	// Query admin user
	var admin struct {
		ID   int    `gorm:"column:gid"`
		Name string `gorm:"column:gname"`
		Pwd  string `gorm:"column:gpwd"`
		SP   int    `gorm:"column:super"`
	}

	result := model.GetDB().Table("oa_gadmin").
		Where("gname = ?", req.Username).
		First(&admin)
	if result.Error != nil {
		return nil, errcode.ErrLoginFailed
	}

	// Verify password
	encryptedPwd := crypto.MD6(req.Password, "", "agg_")
	if admin.Pwd != encryptedPwd {
		return nil, errcode.ErrLoginFailed
	}

	// Generate token
	token := crypto.MD5(fmt.Sprintf("%d_%s_%d", admin.ID, admin.Name, time.Now().UnixNano()))

	adminJSON, _ := json.Marshal(map[string]any{
		"id":   admin.ID,
		"gid":  admin.ID,
		"name": admin.Name,
		"sp":   admin.SP,
	})

	cfg := config.Get()
	expire := time.Duration(cfg.Session.Expire) * time.Second
	if err := cache.Set("manage:token:"+token, string(adminJSON), expire); err != nil {
		return nil, errcode.ErrSystemError
	}

	return &LoginResp{
		Token: token,
		Admin: map[string]any{
			"id":   admin.ID,
			"name": admin.Name,
			"sp":   admin.SP,
		},
	}, nil
}
