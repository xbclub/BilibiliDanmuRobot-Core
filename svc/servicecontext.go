package svc

import (
	"fmt"
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/model"
	"gorm.io/gorm"
)

// SafeUidSet 跨 goroutine 共享的 uid 集合，必须通过方法访问，不能直接操作内部 map
type SafeUidSet struct {
	mu sync.RWMutex
	m  map[int64]bool
}

func NewSafeUidSet() *SafeUidSet {
	return &SafeUidSet{m: make(map[int64]bool)}
}

func (s *SafeUidSet) Has(uid int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.m[uid]
}

func (s *SafeUidSet) Add(uid int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[uid] = true
}

func (s *SafeUidSet) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k := range s.m {
		delete(s.m, k)
	}
}

// PauseWelcome 天选/红包活动期间临时关闭欢迎弹幕
func (s *ServiceContext) PauseWelcome() {
	s.Config.InteractWord = false
	s.Config.EntryEffect = false
	s.Config.WelcomeHighWealthy = false
}

// ResumeWelcome 活动结束后把欢迎弹幕恢复为主播设定的默认状态
func (s *ServiceContext) ResumeWelcome() {
	if s.Config.InteractWord != s.Autointerract.InteractWord {
		s.Config.InteractWord = s.Autointerract.InteractWord
	}
	if s.Config.EntryEffect != s.Autointerract.EntryEffect {
		s.Config.EntryEffect = s.Autointerract.EntryEffect
	}
	if s.Config.WelcomeHighWealthy != s.Autointerract.WelcomeHighWealthy {
		s.Config.WelcomeHighWealthy = s.Autointerract.WelcomeHighWealthy
	}
}

// SetWelcomeDefault 主播指令设置欢迎弹幕开关，同时更新恢复用的默认值
func (s *ServiceContext) SetWelcomeDefault(on bool) {
	s.Config.InteractWord = on
	s.Config.EntryEffect = on
	s.Config.WelcomeHighWealthy = on
	s.Autointerract.InteractWord = on
	s.Autointerract.EntryEffect = on
	s.Autointerract.WelcomeHighWealthy = on
}

type ServiceContext struct {
	Config            *config.Config
	OtherSideUid      *SafeUidSet
	SignInModel       model.SignInModel
	DanmuCntModel     model.DanmuCntModel
	BlindBoxStatModel model.BlindBoxStatModel
	UserID            int64 //主播id
	Autointerract     struct {
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
	return &ServiceContext{
		OtherSideUid:      NewSafeUidSet(),
		SignInModel:       model.NewSignInModel(db, int64(c.RoomId)),
		DanmuCntModel:     model.NewDanmuCntModel(db, int64(c.RoomId)),
		BlindBoxStatModel: model.NewBlindBoxStatModel(db, int64(c.RoomId)),
		Config:            &c,
		UserID:            0,
	}
}
