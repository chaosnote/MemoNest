package service

import (
	"encoding/json"
	"fmt"
	"os"

	"idv/chris/MemoNest/internal/model"
	"idv/chris/MemoNest/utils"

	"go.uber.org/zap"
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

	logger := utils.NewConsoleLogger("console", 1)
	logger.Debug("server", zap.Any("addr", fmt.Sprintf("http://localhost:%s", cfg.Gin.Port)))

	return &cfg, nil
}
