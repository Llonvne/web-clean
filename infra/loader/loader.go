package loader

import (
	"web-clean/domain"
	"web-clean/infra/conf"
)

type Context struct {
	Config *LoadConfig
	Log    domain.Log
}

type LoadConfig struct {
	Paths []string
	Files []string
}

type Loader interface {
	Load(ctx *Context) (*conf.Conf, error)
}

func Default() *LoadConfig {
	return &LoadConfig{
		Paths: []string{".", "./config"},
		Files: []string{"config.json", "app.json"},
	}
}
