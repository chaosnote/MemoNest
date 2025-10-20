package service

import (
	"database/sql"

	"github.com/gin-contrib/sessions/redis"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"

	"idv/chris/MemoNest/model"
)

// DI 注入所有服務依賴 (Fx.In)
type DI struct {
	fx.In

	Config  *model.APPConfig
	RedisDB redis.Store
	MariaDB *sql.DB
	MongoDB *mongo.Client
	NatsIO  *nats.Conn
	TPA     *TPAImpl
	Flag    *FlagImpl
}
