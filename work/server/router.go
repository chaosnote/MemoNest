package server

import (
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/server/controllers"
	"idv/chris/MemoNest/service"
)

// RegisterRoutes 註冊所有路由
func RegisterRoutes(engine *gin.Engine, di service.DI, repo_node repo.NodeRepository, repo_article repo.ArticleRepository) {
	controllers.NewIndexController(engine, di)
	controllers.NewAssetController(engine, di)

	const prefix = "/api/v1"
	g := engine.Group(prefix)
	controllers.NewMemberController(g, di)
	controllers.NewToolsController(g, di)
	controllers.NewNodeController(g, di, repo_node)
	controllers.NewArticleController(g, di, repo_article)
}
