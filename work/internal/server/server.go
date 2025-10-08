package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/service"
	"idv/chris/MemoNest/utils"
)

// RegisterServer 啟動 HTTP 服務
func RegisterServer(lc fx.Lifecycle, engine *gin.Engine, logger *zap.Logger, deps service.Deps) {
	utils.RSAInit("./dist/crypt/rsa.txt", 1024, true)

	// 註冊路由
	RegisterRoutes(engine, logger, deps)

	addr := fmt.Sprintf(":%s", deps.Config.Gin.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			deps.Mariadb.Close()
			deps.MongoDB.Close(ctx)
			deps.NatsIO.Close()
			deps.Redis.Close()
			return srv.Shutdown(ctx) // 可在此關閉資料庫連線、釋放資源
		},
	})
}
