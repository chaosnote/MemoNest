package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	// "go.uber.org/zap"

	"idv/chris/MemoNest/adapter/dependency"
	"idv/chris/MemoNest/adapter/repository/mongo"
	"idv/chris/MemoNest/adapter/repository/nats_io"
	"idv/chris/MemoNest/adapter/repository/redis"
	"idv/chris/MemoNest/api/http"
	"idv/chris/MemoNest/config"
	"idv/chris/MemoNest/service"
	"idv/chris/MemoNest/utils"
)

func main() {
	app := fx.New(
		dependency.Module,

		// 提供設定檔與各種 Client
		fx.Provide(
			// fx.Annotate(
			// 	func() *zap.Logger {
			// 		return logger.Named("system")
			// 	},
			// 	fx.ResultTags(`name:"system"`),
			// ),
			// func() *zap.Logger {
			//  logger := utils.NewFileLogger("./dist/logs/system", "console", 1)
			// 	logger := utils.NewConsoleLogger("console", 1)
			// 	return logger.Named("system")
			// },
			config.NewAPPConfig,
			redis.NewRedisDB,
			mongo.NewMongoDB,
			nats_io.NewNatsIO,
			service.NewFlagImpl,
			// fx.Annotate(
			// 	service.NewNatsIOClient,
			// 	fx.ParamTags(``, ``, `name:"system"`), // `` 為預設值、留意注入參數順序，需對應函式參數
			// ),
			service.NewTPAImpl,

			// Gin Server
			http.NewServerRoute,
			// Router
		),

		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{Logger: utils.NewFileLogger("./dist/logs/fx", "console", 1)}
		}),

		// 啟動 HTTP Server
		fx.Invoke(http.RegisterRoutes),
	)

	app.Run()
}
