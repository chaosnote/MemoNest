package http

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/config"
	"idv/chris/MemoNest/server/middleware"
	"idv/chris/MemoNest/utils"
)

// NewServerRoute 建立 Gin Engine
func NewServerRoute(cfg *config.APPConfig, store redis.Store) *gin.Engine {
	logger := utils.NewFileLogger("./dist/logs/server", "console", 1)

	gin.SetMode(cfg.Gin.Mode)

	r := gin.New()
	r.Use(sessions.Sessions("custom_session", store))
	r.Use(middleware.NewLoggerMiddleware(logger))
	r.Use(middleware.NewRecoveryMiddleware(logger))

	return r
}
