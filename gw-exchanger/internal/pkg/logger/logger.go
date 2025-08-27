package logger

import (
	"go.uber.org/zap"
)

var L *zap.SugaredLogger

func Init() {
	var baseLogger *zap.Logger

	baseLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	L = baseLogger.Sugar()
}
