package service

import (
	"go.uber.org/fx"

	"idv/chris/MemoNest/internal/model"
)

// Deps 注入所有服務依賴 (Fx.In)
// Deps 為 Dependencies 的常見縮寫
type Deps struct {
	fx.In

	Config  *model.APPConfig
	Redis   *RedisDBImpl
	Mariadb *MariaDBImpl
	MongoDB *MongoDBImpl
	NatsIO  *NatsIOImpl
	API     *APIImpl
	Flag    *FlagImpl
}
