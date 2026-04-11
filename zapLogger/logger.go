package zapLogger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Create() (*zap.Logger, func() error, error) {
	if err := os.MkdirAll("logs", 755); err != nil {
		return &zap.Logger{}, nil, err
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000000")
	logFilePath := filepath.Join("logs", fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return &zap.Logger{}, nil, err
	}

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	level := zapcore.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), level),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger, logFile.Close, nil
}
