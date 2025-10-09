package service

import (
	"context"

	"github.com/go-redis/redis/v8"

	"idv/chris/MemoNest/internal/model"
)

// NewRedisDB 建立 Redis 連線
func NewRedisDB(cfg *model.APPConfig) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 測試連線
	if e := r.Ping(context.Background()).Err(); e != nil {
		return nil, e
	}

	return r, nil
}
