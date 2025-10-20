package dependency

import (
	"idv/chris/MemoNest/adapter/repository/mysql"
	"idv/chris/MemoNest/domain/repo"

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
	),
)
