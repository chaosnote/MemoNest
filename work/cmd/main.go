package main

import (
	"idv/chris/MemoNest/internal/server"
	"idv/chris/MemoNest/internal/service"
	"idv/chris/MemoNest/utils"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	// "go.uber.org/zap"
)

func main() {
	app := fx.New(
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
			service.NewAPPConfig,
			service.NewRedisDB,
			service.NewMariaDB,
			service.NewMongoDB,
			service.NewNatsIOImpl,
			service.NewFlagImpl,
			// fx.Annotate(
			// 	service.NewNatsIOClient,
			// 	fx.ParamTags(``, ``, `name:"system"`), // `` 為預設值、留意注入參數順序，需對應函式參數
			// ),
			service.NewTPAImpl,
			server.NewGinEngine,
		),

		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{Logger: utils.NewFileLogger("./dist/logs/fx", "console", 1)}
		}),

		// 啟動 HTTP Server
		fx.Invoke(server.Register),
	)

	app.Run()
}
