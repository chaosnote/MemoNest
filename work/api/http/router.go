package http

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/adapter/http/middleware"
	"idv/chris/MemoNest/config"
	"idv/chris/MemoNest/utils"

	"idv/chris/MemoNest/domain/repo"
)

func NewServerRoute(cfg *config.APPConfig, store redis.Store) *gin.Engine {
	logger := utils.NewFileLogger("./dist/logs/server", "console", 1)

	gin.SetMode(cfg.Gin.Mode)

	r := gin.New()
	r.Use(sessions.Sessions("custom_session", store))
	r.Use(middleware.NewLoggerMiddleware(logger))
	r.Use(middleware.NewRecoveryMiddleware(logger))

	return r
}

func RegisterRoutes(engine *gin.Engine, nodeRepo repo.NodeRepository) {
	_ = engine.Group("/api/v1")
}
