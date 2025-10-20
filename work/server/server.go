package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"idv/chris/MemoNest/domain/repo"
	"idv/chris/MemoNest/service"
	"idv/chris/MemoNest/utils"
)

// Register 啟動 HTTP 服務
func Register(lc fx.Lifecycle, engine *gin.Engine, deps service.DI, repo repo.NodeRepository) {
	utils.RSAInit("./dist/logs/crypt/rsa.txt", 1024, true)

	// 註冊路由
	RegisterRoutes(engine, deps, repo)

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
			deps.MariaDB.Close()
			deps.MongoDB.Disconnect(ctx)
			deps.NatsIO.Close()
			return srv.Shutdown(ctx) // 可在此關閉資料庫連線、釋放資源
		},
	})
}
