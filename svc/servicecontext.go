package svc

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/model"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	OtherSideUid map[int64]bool
	TopUid       map[int64]bool
	SininModel   model.SingInModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	dbFile := fmt.Sprintf("%s/%s?_pragma=busy_timeout(5000)", c.DBPath, c.DBName)
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		OtherSideUid: make(map[int64]bool),
		TopUid:       make(map[int64]bool),
		SininModel:   model.NewSingInModel(db, int64(c.RoomId)),
		Config:       c,
	}
}
