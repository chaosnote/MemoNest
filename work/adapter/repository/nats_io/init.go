package nats_io

import (
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"idv/chris/MemoNest/config"
	"idv/chris/MemoNest/utils"
)

// NewNatsIO 建立 NATS 連線
func NewNatsIO(cfg *config.APPConfig) (*nats.Conn, error) {
	logger := utils.NewFileLogger("./dist/logs/natsio", "console", 1)
	c, e := nats.Connect(cfg.Natsio.URL,
		nats.PingInterval(time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Error("nats disconnect", zap.Error(err))
		}),
	)
	if e != nil {
		return nil, e
	}

	return c, nil
}
