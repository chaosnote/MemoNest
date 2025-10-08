package utils

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

var (
	serial_id_mu  sync.Mutex
	serial_id_pre int64 = time.Now().UTC().Unix()
	serial_id_sub       = map[string]int64{}
)

// GenSerial 產生流水號
//
// 參數:
//   - key: 關鍵字<例:GameID/UID>
func GenSerial(key string) string {
	serial_id_mu.Lock()
	defer serial_id_mu.Unlock()

	var current = time.Now().UTC().Unix()
	var tmp_sub int64
	if current != serial_id_pre {
		serial_id_pre = current
		serial_id_sub[key] = 0
	} else {
		tmp_sub = serial_id_sub[key]
		tmp_sub++
		serial_id_sub[key] = tmp_sub
	}
	pre := decimal.NewFromInt(serial_id_pre).String()
	sub := decimal.NewFromInt(tmp_sub).String()
	if len(sub) > 0 {
		sub = "." + sub
	}

	return key + pre + sub
}
