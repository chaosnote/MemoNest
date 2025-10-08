package service

import (
	"github.com/go-resty/resty/v2"

	"idv/chris/MemoNest/internal/model"
)

// APIImpl REST API 客戶端結構
type APIImpl struct {
	Client *resty.Client
}

// NewAPIImpl 建立 REST API 客戶端
func NewAPIImpl(cfg *model.APPConfig) (*APIImpl, error) {
	client := resty.New().SetBaseURL(cfg.API.BaseURL)
	return &APIImpl{Client: client}, nil
}
