package svc

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
)

type ServiceContext struct {
	Config       config.Config
	OtherSideUid map[int64]bool
	TopUid       map[int64]bool
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		OtherSideUid: make(map[int64]bool),
		TopUid:       make(map[int64]bool),
		Config:       c,
	}
}
