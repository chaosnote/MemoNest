package server

import (
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/server/controllers"
	"idv/chris/MemoNest/internal/service"
)

// RegisterRoutes 註冊所有路由
func RegisterRoutes(engine *gin.Engine, deps service.DI) {
	const prefix = "/api/v1"
	controllers.NewIndexController(engine, deps)
	controllers.NewToolsController(engine.Group(prefix), deps)
	controllers.NewMemberController(engine.Group(prefix), deps)
}
