package log

import (
	"go.uber.org/zap"

	"web-clean/domain"
)

type _zap struct {
	*zap.SugaredLogger
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
