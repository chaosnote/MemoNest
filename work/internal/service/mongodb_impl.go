package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"idv/chris/MemoNest/internal/model"
)

// MongoDBImpl MongoDB 客戶端結構
type MongoDBImpl struct {
	client *mongo.Client
}

func (mds *MongoDBImpl) Close(ctx context.Context) error {
	return mds.client.Disconnect(ctx)
}

// NewMongoDBImpl 建立 MongoDB 連線
func NewMongoDBImpl(cfg *model.APPConfig) (*MongoDBImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, e := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongodb.URI))
	if e != nil {
		return nil, e
	}

	return &MongoDBImpl{client: client}, nil
}
