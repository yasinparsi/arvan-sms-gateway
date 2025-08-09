package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func SetupLogger() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Log = l.Sugar()
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
