package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"

	"idv/chris/MemoNest/adapter/http/middleware"
	"idv/chris/MemoNest/api/http/handle"
	"idv/chris/MemoNest/application/usecase"
	"idv/chris/MemoNest/config"
	"idv/chris/MemoNest/domain/service"
	"idv/chris/MemoNest/utils"
)

const API_VER = "/api/v1"

func NewServerRoute(
	lc fx.Lifecycle,
	cfg *config.APPConfig,
	redis_store redis.Store,
	maria_db *sql.DB,
	mongo_db *mongo.Client,
	nats_io *nats.Conn,
) *gin.Engine {
	utils.RSAInit("./dist/logs/crypt/rsa.txt", 1024, true)

	logger := utils.NewFileLogger("./dist/logs/server", "console", 1)

	gin.SetMode(cfg.Gin.Mode)

	engine := gin.New()
	engine.Use(sessions.Sessions("custom_session", redis_store))
	engine.Use(middleware.GinLogger(logger))
	engine.Use(middleware.GinRecovery(logger))

	addr := fmt.Sprintf(":%s", cfg.Gin.Port)
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
			maria_db.Close()
			mongo_db.Disconnect(ctx)
			nats_io.Close()
			return srv.Shutdown(ctx) // 可在此關閉資料庫連線、釋放資源
		},
	})

	return engine
}

func NewIndexHandler(
	cfg *config.APPConfig,
	engine *gin.Engine,
	uc *usecase.IndexUsecase,
	session service.Session,
) {
	h := &handle.IndexHandler{
		Debug:   cfg.Gin.Mode == "debug",
		UC:      uc,
		Session: session,
	}
	engine.GET("/", h.Entry)
	engine.GET("/health", h.Health)
}

func NewToolHandler(
	engine *gin.Engine,
	uc *usecase.ToolUsecase,
	session service.Session,
) {
	h := &handle.ToolHandler{
		UC: uc,
	}
	r := engine.Group("/tools")
	r.Use(middleware.Auth(session))
	r.GET("/uuid", h.GenUUID)
}

func NewMemberHandler(
	engine *gin.Engine,
	uc *usecase.MemberUsecase,
) {
	c := &handle.MemberHandler{}
	r := engine.Group(filepath.Join(API_VER, "/member"))
	r.GET("/login", c.Login)
	r.GET("/logout", c.Logout)
}
