package service

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"idv/chris/MemoNest/internal/model"
)

// RedisDBImpl Redis 客戶端結構
type RedisDBImpl struct {
	client *redis.Client
}

func (rds *RedisDBImpl) SetToken(ctx context.Context, key, value string) (e error) {
	cmd := rds.client.Set(ctx, key, value, 45*time.Minute)
	e = cmd.Err()
	if e != nil {
		return
	}
	return
}

func (rds *RedisDBImpl) Close() error {
	return rds.client.Close()
}

// NewRedisDBImpl 建立 Redis 連線
func NewRedisDBImpl(cfg *model.APPConfig) (*RedisDBImpl, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 測試連線
	if e := rdb.Ping(context.Background()).Err(); e != nil {
		return nil, e
	}

	return &RedisDBImpl{client: rdb}, nil
}
