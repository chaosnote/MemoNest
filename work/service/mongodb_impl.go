package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"idv/chris/MemoNest/model"
)

// NewMongoDB 建立 MongoDB 連線
func NewMongoDB(cfg *model.APPConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, e := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongodb.URI))
	if e != nil {
		return nil, e
	}

	return client, nil
}
