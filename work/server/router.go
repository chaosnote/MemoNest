package server

import (
	"github.com/gin-gonic/gin"

	"idv/chris/MemoNest/config"
	"idv/chris/MemoNest/domain/repo"
	zzz "idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/server/controllers"
)

// RegisterRoutes 註冊所有路由
func RegisterRoutes(
	engine *gin.Engine,
	cfg *config.APPConfig,
	repo_node repo.NodeRepository,
	repo_article repo.ArticleRepository,
	service_menu zzz.MenuProvider,
	service_tree zzz.NodeTree,
	service_img zzz.ImageProcessor,
) {
	controllers.NewIndexController(engine, cfg)
	controllers.NewAssetController(engine, di, service_img)

	const prefix = "/api/v1"
	g := engine.Group(prefix)
	controllers.NewMemberController(g, di)
	controllers.NewToolsController(g)
	controllers.NewNodeController(g, di, repo_node, service_menu, service_tree)
	controllers.NewArticleController(g, di, repo_article, service_menu, service_tree, service_img)
}
