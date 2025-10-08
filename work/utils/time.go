package utils

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

/*
時間戳記
*/

/*
UTCUnix ex. PHP
*/
func UTCUnix() int64 {
	return time.Now().UTC().Unix()
}

/*
UTCUnixString ex. PHP
*/
func UTCUnixString() string {
	return decimal.NewFromInt(UTCUnix()).String()
}

/*
UTCUnixNano timestamp ex. JS
*/
func UTCUnixNano() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

/*
TimeFromUnixString unix 字串轉換為 time
*/
func TimeFromUnixString(s string) (time.Time, error) {
	i, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return time.Now(), e
	}
	t := time.Unix(i, 0)
	return t, nil
}
