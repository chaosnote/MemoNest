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
	"go.uber.org/zap"

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
	cli *config.CLIFlags,
	redis_store redis.Store,
	maria_db *sql.DB,
	mongo_db *mongo.Client,
	nats_io *nats.Conn,
) *gin.Engine {
	utils.RSAInit("./dist/logs/crypt/rsa.txt", 1024, true)

	var logger *zap.Logger
	if cli.Debug {
		logger = utils.NewConsoleLogger("console", 1)
	} else {
		logger = utils.NewFileLogger("./dist/logs/server", "console", 1)
	}

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
	session service.Session,
) {
	h := &handle.MemberHandler{
		UC:      uc,
		Session: session,
	}
	r := engine.Group(filepath.Join(API_VER, "/member"))
	r.POST("/login", h.Login)
	r.GET("/logout", h.Logout)
}

func NewNodeHandler(
	engine *gin.Engine,
	uc *usecase.NodeUsecase,
	session service.Session,
) {
	h := &handle.NodeHandler{
		UC:      uc,
		Session: session,
	}
	r := engine.Group(filepath.Join(API_VER, "/node"))
	r.Use(middleware.Auth(session))
	r.GET("/list", h.List)
	r.POST("/add", h.Add)
	r.POST("/del", h.Del)
	r.POST("/edit", h.Edit)
	r.POST("/move", h.Move)
}

func NewAssetHandler(
	engine *gin.Engine,
	uc *usecase.AssetUsecase,
	session service.Session,
) {
	h := &handle.AssetHandler{
		UC:      uc,
		Session: session,
	}

	r := engine.Group("/asset/article")
	r.Use(middleware.Auth(session))
	r.GET("/image/:id/:name", h.Image)
}

func NewArticleHandler(
	engine *gin.Engine,
	uc *usecase.ArticleUsecase,
	session service.Session,
) {
	h := &handle.ArticleHandler{
		UC:      uc,
		Session: session,
	}
	r := engine.Group(filepath.Join(API_VER, "/article"))
	r.Use(middleware.Auth(session))
	r.GET("/fresh", h.Fresh)
	r.GET("/list", h.List)
	r.POST("/add", h.Add)
	r.POST("/del", h.Del)
	r.POST("/renew", h.Renew)
	r.GET("/edit/:id", h.Edit)
}
