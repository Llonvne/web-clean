package log

import (
	"go.uber.org/zap"

	"web-clean/domain"
)

type _zap struct {
	*zap.SugaredLogger
}

func (z *_zap) DPanic(args ...interface{}) {
	z.SugaredLogger.DPanic(args...)
}

func (z *_zap) DPanicf(template string, args ...interface{}) {
	z.SugaredLogger.DPanicf(template, args...)
}

func (z *_zap) DPanicw(msg string, keysAndValues ...interface{}) {
	z.SugaredLogger.DPanicw(msg, keysAndValues...)
}

func Zap() domain.Log {
	log, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	sugar := log.Sugar()

	return &_zap{
		SugaredLogger: sugar,
	}
}
