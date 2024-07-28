package svc

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/model"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        *config.Config
	OtherSideUid  map[int64]bool
	TopUid        map[int64]bool
	SignInModel   model.SignInModel
	DanmuCntModel model.DanmuCntModel
	UserID        int64 //主播id
	Autointerract struct {
		EntryEffect        bool
		WelcomeHighWealthy bool
		InteractWord       bool
	}
	RobotID string //机器人uid
}

func NewServiceContext(c config.Config) *ServiceContext {
	dbFile := fmt.Sprintf("%s/%s?_pragma=busy_timeout(5000)", c.DBPath, c.DBName)
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DanmuCntdbFile := fmt.Sprintf("%s/%s?_pragma=busy_timeout(5000)", c.DanmuCntDBPath, c.DanmuCntDBName)
	DanmuCntDB, DCerr := gorm.Open(sqlite.Open(DanmuCntdbFile), &gorm.Config{})
	if DCerr != nil {
		panic(DCerr)
	}

	return &ServiceContext{
		OtherSideUid:  make(map[int64]bool),
		TopUid:        make(map[int64]bool),
		SignInModel:   model.NewSignInModel(db, int64(c.RoomId)),
		DanmuCntModel: model.NewDanmuCntModel(DanmuCntDB, int64(c.RoomId)),
		Config:        &c,
		UserID:        0,
	}
}
