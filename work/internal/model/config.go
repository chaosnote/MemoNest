package model

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
