package dependency

import (
	"idv/chris/MemoNest/adapter/infra"
	"idv/chris/MemoNest/adapter/repository/mysql"
	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/domain/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		mysql.NewMariaDB,
		//
		fx.Annotate(
			mysql.NewNodeRepo,
			fx.As(new(repo.NodeRepository)),
		),
		fx.Annotate(
			mysql.NewArticleRepo,
			fx.As(new(repo.ArticleRepository)),
		),
		fx.Annotate(
			infra.NewMenuProvider,
			fx.As(new(service.MenuProvider)),
		),
		fx.Annotate(
			infra.NewNodeTree,
			fx.As(new(service.NodeTree)),
		),
		fx.Annotate(
			infra.NewImageProcessor,
			fx.As(new(service.ImageProcessor)),
		),
	),
)
