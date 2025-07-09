package infra

import (
	"context"

	"web-clean/domain"
	"web-clean/infra/conf"
	"web-clean/infra/loader"
	"web-clean/infra/log"
)

// Context 为应用程序提供核心功能组件，支持直接嵌入业务结构体。
//
// 嵌入说明：
//   - 通过匿名嵌入方式将 Context 集成到业务结构体中
//   - 嵌入后可直接访问 Log 和 Conf 字段
//
// 冲突解决：
//   - 若业务结构体存在同名字段，需修改业务结构体字段名
//   - Context 的内置字段名（Log/Conf）具有保留优先级
type Context struct {
	Log  domain.Log
	Conf *conf.Conf
	Ctx  context.Context
}

type PrepareConfig struct {
	Loader loader.Loader
	config *loader.LoadConfig
}

func Prepare(prepare PrepareConfig) (*Context, error) {

	logger := log.Zap()

	if prepare.config == nil {
		prepare.config = loader.Default()
	}

	config, err := prepare.Loader.Load(&loader.Context{
		Config: prepare.config,
		Log:    logger,
	})
	if err != nil {
		return nil, err
	}

	c := &Context{
		Log:  logger,
		Ctx:  context.Background(),
		Conf: config,
	}

	return c, nil
}
