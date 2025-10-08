package utils

import (
	"fmt"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// encodeTimeExcludingTimezone 自定義時間編碼器，排除時區資訊
func encodeTimeExcludingTimezone(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("01-02 15:04:05.000"))
}

// 自定義日誌等級編碼器 (已註解)
// func levelEncoder(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
// 	encoder.AppendString("[" + strings.ToUpper(level.String()[:1]) + "]") // 只取首字母，並加上中括號
// }

func NewEmptyLogger() *zap.Logger {
	return zap.NewNop()
}

// NewConsoleLogger 建立一個用於控制台輸出的日誌記錄器。
//
// 參數:
//
//	encoding: 日誌編碼格式，可以是 "json" 或 "console"。
//	callerSkip: 呼叫堆疊中要跳過的層數，通常用於調整日誌中顯示的呼叫者位置。
//
// 回傳:
//
//	*zap.Logger
func NewConsoleLogger(encoding string, callerSkip int) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.LevelKey = zapcore.OmitKey
	// config.EncoderConfig.MessageKey = zapcore.OmitKey
	config.EncoderConfig.EncodeTime = encodeTimeExcludingTimezone
	// config.EncoderConfig.TimeKey = zapcore.OmitKey
	config.EncoderConfig.CallerKey = zapcore.OmitKey
	config.Encoding = encoding

	logger, err := config.Build(zap.AddCallerSkip(callerSkip))
	if err != nil {
		panic(err)
	}
	return logger
}

// NewFileLogger 建立並回傳一個寫入檔案的日誌記錄器。
//
// 參數:
//
//	dir: 日誌檔案存放的目錄路徑。
//	encoding: 日誌的編碼格式，例如 "json" 或 "console"。
//	callerSkip: 呼叫者日誌層級的跳過層數，用於正確地記錄呼叫者的資訊。
//
// 回傳:
//
//	*zap.Logger
func NewFileLogger(dir string, encoding string, callerSkip int) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   path.Join(dir, ".log"), // 文件輸出路徑
		MaxSize:    10,                     // 文件最大大小 (MB)
		LocalTime:  true,                   // 使用本地時間
		Compress:   false,                  // 是否壓縮檔案
		MaxAge:     30,                     // 舊檔案保留天數
		MaxBackups: 50,                     // 最多備份檔案數量
	}
	writeSyncer := zapcore.AddSync(&hook)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = zapcore.OmitKey
	// encoderConfig.MessageKey = zapcore.OmitKey
	encoderConfig.EncodeTime = encodeTimeExcludingTimezone
	encoderConfig.CallerKey = zapcore.OmitKey
	// encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 短檔案路徑

	var encoder zapcore.Encoder
	switch encoding {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		panic(fmt.Sprintf("unknow encoding %s", encoder))
	}

	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		zapcore.InfoLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(callerSkip))
}
