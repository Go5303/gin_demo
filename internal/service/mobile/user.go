package mobile

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Go5303/gin_demo/config"
	"github.com/Go5303/gin_demo/internal/model"
	"github.com/Go5303/gin_demo/pkg/cache"
	"github.com/Go5303/gin_demo/pkg/crypto"
	"github.com/Go5303/gin_demo/pkg/errcode"
)

// LoginReq is the mobile login request body
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResp is the mobile login response data
type LoginResp struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

// Login handles mobile user login logic
func Login(req *LoginReq) (*LoginResp, error) {
	// Query user from oa_user
	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errcode.ErrLoginFailed
	}

	// Verify password
	encryptedPwd := crypto.MD6(req.Password, "", "agg_")
	if user.Password != encryptedPwd {
		return nil, errcode.ErrLoginFailed
	}

	// Generate token
	token := crypto.MD5(fmt.Sprintf("%d_%s_%d", user.ID, user.Username, time.Now().UnixNano()))

	userJSON, _ := json.Marshal(map[string]any{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"phone":    user.Phone,
	})

	cfg := config.Get()
	expire := time.Duration(cfg.Session.Expire) * time.Second
	if err := cache.Set("mobile:token:"+token, string(userJSON), expire); err != nil {
		return nil, errcode.ErrSystemError
	}

	return &LoginResp{
		Token: token,
		User: map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"phone":    user.Phone,
		},
	}, nil
}
