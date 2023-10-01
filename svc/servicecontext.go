package svc

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
