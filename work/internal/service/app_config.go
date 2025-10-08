package service

import (
	"encoding/json"
	"os"

	"idv/chris/MemoNest/internal/model"
)

// NewAPPConfig 讀取設定檔並返回 APPConfig
func NewAPPConfig() (*model.APPConfig, error) {
	data, err := os.ReadFile("./assets/config.json")
	if err != nil {
		return nil, err
	}

	var cfg model.APPConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
