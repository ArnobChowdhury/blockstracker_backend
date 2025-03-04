package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func init() {
	rawLogger, _ := zap.NewProduction()
	Log = rawLogger.Sugar()
}

func LoggerProvider() *zap.SugaredLogger {
	return Log
}
