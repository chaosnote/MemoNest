package dependency

import (
	"idv/chris/MemoNest/adapter/http"
	"idv/chris/MemoNest/adapter/infra"
	"idv/chris/MemoNest/adapter/repository/mongo"
	"idv/chris/MemoNest/adapter/repository/mysql"
	"idv/chris/MemoNest/adapter/repository/nats_io"
	"idv/chris/MemoNest/adapter/repository/redis"
	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/domain/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		// 資料來源
		mysql.NewMariaDB,
		redis.NewRedisDB,
		mongo.NewMongoDB,
		nats_io.NewNatsIO,
		// repository（保留 fx.As，因為是 interface 實作）
		fx.Annotate(mysql.NewNodeRepo, fx.As(new(repo.NodeRepository))),
		fx.Annotate(mysql.NewArticleRepo, fx.As(new(repo.ArticleRepository))),
		// domain service adapter 實作（保留 fx.As）
		fx.Annotate(http.NewGinSession, fx.As(new(service.Session))),
		fx.Annotate(infra.NewMenuProvider, fx.As(new(service.MenuProvider))),
		fx.Annotate(infra.NewNodeTree, fx.As(new(service.NodeTree))),
		fx.Annotate(infra.NewImageProcessor, fx.As(new(service.ImageProcessor))),
		// usecase
		usecase.NewIndexUsecase,
		usecase.NewToolUsecase,
	),
)
