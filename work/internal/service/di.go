package service

import (
	"go.uber.org/fx"

	"idv/chris/MemoNest/internal/model"
)

// DI 注入所有服務依賴 (Fx.In)
type DI struct {
	fx.In

	Config  *model.APPConfig
	Redis   *RedisDBImpl
	Mariadb *MariaDBImpl
	MongoDB *MongoDBImpl
	NatsIO  *NatsIOImpl
	TPA     *TPAImpl
	Flag    *FlagImpl
}
