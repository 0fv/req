package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewSimplelog() *zap.Logger {
	var cores []zapcore.Core
	cores = append(cores, newConsole())
	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.AddCaller())
	return logger
}

func newConsole() zapcore.Core {
	consoleWrite := zapcore.AddSync(io.MultiWriter(os.Stdout))
	consoleConfig := zap.NewProductionEncoderConfig()
	// 控制台输出颜色。
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	/** 定义日志控制台输出核心。*/
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleConfig),
		consoleWrite,
		zap.DebugLevel,
	)
	return consoleCore

}
