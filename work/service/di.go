package service

import (
	"database/sql"

	"github.com/gin-contrib/sessions/redis"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"

	"idv/chris/MemoNest/config"
)

// DI 注入所有服務依賴 (Fx.In)
type DI struct {
	fx.In

	Config  *config.APPConfig
	RedisDB redis.Store
	MariaDB *sql.DB
	MongoDB *mongo.Client
	NatsIO  *nats.Conn
	TPA     *TPAImpl
	Flag    *FlagImpl
}
