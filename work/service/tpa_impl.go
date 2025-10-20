package service

import (
	"github.com/go-resty/resty/v2"

	"idv/chris/MemoNest/model"
)

// TPAImpl Third Party API Implement
type TPAImpl struct {
	Client *resty.Client
}

// NewTPAImpl 建立 REST API 客戶端
func NewTPAImpl(cfg *model.APPConfig) (*TPAImpl, error) {
	client := resty.New().SetBaseURL(cfg.API.BaseURL)
	return &TPAImpl{Client: client}, nil
}
