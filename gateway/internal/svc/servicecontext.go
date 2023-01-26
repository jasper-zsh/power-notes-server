package svc

import (
	"powernotes-server/gateway/internal/config"
	"powernotes-server/gateway/internal/model"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	model.InitDB(c.Database)
	return &ServiceContext{
		Config: c,
	}
}
