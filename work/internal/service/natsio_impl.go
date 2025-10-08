package service

import (
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"idv/chris/MemoNest/internal/model"
)

// NatsIOImpl NATS 客戶端結構
type NatsIOImpl struct {
	conn *nats.Conn
}

func (nis *NatsIOImpl) Close() {
	nis.conn.Close()
}

// NewNatsIOImpl 建立 NATS 連線
func NewNatsIOImpl(cfg *model.APPConfig, logger *zap.Logger) (*NatsIOImpl, error) {
	c, e := nats.Connect(cfg.Natsio.URL,
		nats.PingInterval(time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Error("nats disconnect", zap.Error(err))
		}),
	)
	if e != nil {
		return nil, e
	}

	return &NatsIOImpl{conn: c}, nil
}
