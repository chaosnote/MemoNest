package server

import (
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/internal/server/middleware"
	"idv/chris/MemoNest/utils"
)

// NewGinEngine 建立 Gin Engine 並加入中介層
func NewGinEngine(cfg *model.APPConfig) *gin.Engine {
	logger := utils.NewFileLogger("./dist/server", "console", 1)

	gin.SetMode(cfg.Gin.Mode)

	r := gin.New()
	r.Use(middleware.NewLoggerMiddleware(logger))
	r.Use(middleware.NewRecoveryMiddleware(logger))

	return r
}
