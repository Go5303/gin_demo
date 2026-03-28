package logger

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.SugaredLogger

// Init initializes the global logger
func Init(level string, logPath string) {
	var lv zapcore.Level
	switch level {
	case "debug":
		lv = zapcore.DebugLevel
	case "info":
		lv = zapcore.InfoLevel
	case "warn":
		lv = zapcore.WarnLevel
	case "error":
		lv = zapcore.ErrorLevel
	default:
		lv = zapcore.InfoLevel
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		CallerKey:    "caller",
		MessageKey:   "msg",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	// 终端输出（Console 格式，方便开发看）
	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		lv,
	))

	// 文件输出（JSON 格式，方便 ELK 采集）
	if logPath != "" {
		if err := os.MkdirAll(filepath.Dir(logPath), 0755); err == nil {
			file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err == nil {
				cores = append(cores, zapcore.NewCore(
					zapcore.NewJSONEncoder(encoderCfg),
					zapcore.AddSync(file),
					lv,
				))
			}
		}
	}

	core := zapcore.NewTee(cores...)
	L = zap.New(core, zap.AddCaller()).Sugar()
}

// WithCtx returns a logger with request_id from gin.Context
func WithCtx(c *gin.Context) *zap.SugaredLogger {
	if requestID, exists := c.Get("request_id"); exists {
		return L.With("request_id", requestID)
	}
	return L
}
