package config

import (
	"encoding/json"
	"os"

	"idv/chris/MemoNest/utils"

	"go.uber.org/zap"
)

// APPConfig 設定檔結構
type APPConfig struct {
	Gin struct {
		Mode string `json:"mode"`
		Port string `json:"port"`
	} `json:"gin"`

	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`

	Mariadb struct {
		DSN string `json:"dsn"`
	} `json:"mariadb"`

	Mongodb struct {
		URI string `json:"uri"`
	} `json:"mongodb"`

	Natsio struct {
		URL string `json:"url"`
	} `json:"natsio"`

	API struct {
		BaseURL string `json:"baseURL"`
	} `json:"api"`
}

func Load(file_path string) (cfg *APPConfig, err error) {
	data, err := os.ReadFile(file_path)
	if err != nil {
		return
	}

	var tmp_config APPConfig
	err = json.Unmarshal(data, &tmp_config)
	if err != nil {
		return
	}
	cfg = &tmp_config
	return
}

// NewAPPConfig 讀取設定檔並返回 APPConfig
func NewAPPConfig() (cfg *APPConfig, err error) {
	cfg, err = Load("./assets/config.json")

	logger := utils.NewConsoleLogger("console", 1)
	logger.Debug("server", zap.Any("addr", "http://172.31.235.34:8080/"))

	return
}
