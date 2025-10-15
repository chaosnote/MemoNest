package server

import (
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/internal/server/controllers"
	"idv/chris/MemoNest/internal/service"
)

// RegisterRoutes 註冊所有路由
func RegisterRoutes(engine *gin.Engine, di service.DI) {
	const prefix = "/api/v1"
	controllers.NewIndexController(engine, di)

	g := engine.Group(prefix)
	controllers.NewMemberController(g, di)
	controllers.NewToolsController(g, di)
	controllers.NewNodeController(g, di)
	controllers.NewArticleController(g, di)
}
