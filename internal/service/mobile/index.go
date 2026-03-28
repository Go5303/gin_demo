package mobile

import (
	"encoding/json"
	"fmt"
	"time"

	"git.woda.ink/Woda_OA/internal/model"
	"git.woda.ink/Woda_OA/pkg/cache"
	"git.woda.ink/Woda_OA/pkg/errcode"
)

// IndexReq is the index request body
type IndexReq struct {
	UserID int `json:"user_id" binding:"required"`
}

// IndexResp is the index response data
type IndexResp struct {
	User     *model.User `json:"user"`
	CacheHit bool        `json:"cache_hit"`
	CacheKey string      `json:"cache_key"`
}

// Index queries oa_user and demonstrates Redis set/get
func Index(req *IndexReq) (*IndexResp, error) {
	cacheKey := fmt.Sprintf("mobile:user:%d", req.UserID)

	// Try to get from Redis first
	val, err := cache.Get(cacheKey)
	if err == nil && val != "" {
		var user model.User
		if json.Unmarshal([]byte(val), &user) == nil {
			return &IndexResp{
				User:     &user,
				CacheHit: true,
				CacheKey: cacheKey,
			}, nil
		}
	}

	// Cache miss, query from DB
	user, err := model.GetUserByID(req.UserID)
	if err != nil {
		return nil, errcode.ErrUserNotFound
	}

	// Set to Redis, expire 5 minutes
	userJSON, _ := json.Marshal(user)
	_ = cache.Set(cacheKey, string(userJSON), 5*time.Minute)

	return &IndexResp{
		User:     user,
		CacheHit: false,
		CacheKey: cacheKey,
	}, nil
}
