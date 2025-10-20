package service

import (
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"idv/chris/MemoNest/model"
	"idv/chris/MemoNest/utils"
)

// NewNatsIOImpl 建立 NATS 連線
func NewNatsIOImpl(cfg *model.APPConfig) (*nats.Conn, error) {
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
