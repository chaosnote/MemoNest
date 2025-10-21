package redis

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"

	"idv/chris/MemoNest/config"
)

// NewRedisDB 建立 Redis 連線
func NewRedisDB(cfg *config.APPConfig) (redis.Store, error) {
	store, err := redis.NewStoreWithDB(
		10,                              // 連線池大小
		"tcp",                           // 網路類型
		cfg.Redis.Addr,                  // Redis 地址
		"",                              // user name
		cfg.Redis.Password,              // 密碼 (如果沒有則留空)
		fmt.Sprintf("%v", cfg.Redis.DB), // 指定 DB
		[]byte("secret-key"),            // Session 密鑰
	)
	if err != nil {
		return nil, err
	}

	store.Options(sessions.Options{
		MaxAge:   300,                     // 存活時間（秒）
		Path:     "/",                     // 全站有效
		HttpOnly: true,                    // 禁止 JS 存取 cookie
		SameSite: http.SameSiteStrictMode, // 防止 CSRF
	})
	// Secure:true,    // 只允許 HTTPS 傳輸

	return store, nil
}
