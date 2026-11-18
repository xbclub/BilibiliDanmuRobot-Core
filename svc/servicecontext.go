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

// EventRLock / EventRUnlock 事件处理入口的读锁：
// 持有期间 reload 与运行时配置变更（写锁）不会插入，单个事件内
// 读到的 Config/models 是一致的组合。约定：持锁期间不得调用任何
// 配置变更器（SetWelcomeDefault/PauseWelcome/ResumeWelcome/MutateConfig），
// 否则同 goroutine 死锁；变更器调用点须在锁外或独立 goroutine 中执行
func (s *ServiceContext) EventRLock() { s.mu.RLock() }

func (s *ServiceContext) EventRUnlock() { s.mu.RUnlock() }

// WithEvent 在事件读锁内执行 fn，供事件入口一行包裹使用
func (s *ServiceContext) WithEvent(fn func()) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn()
}

// MutateConfig 以写时复制方式修改配置：fn 收到当前配置的副本，
// 修改后作为新快照整体替换，不在旧快照上原地写（读方持有的旧指针不受影响）
func (s *ServiceContext) MutateConfig(fn func(*config.Config)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nc := *s.Config
	fn(&nc)
	s.Config = &nc
}

// ApplyConfigHotSwap 配置热更：整体替换 Config 指针。
// 与 Autointerract 默认值不同的三个欢迎开关视为运行时覆盖
// （天选暂停/主播指令），热更后予以保留，不被配置文件覆盖。
func (s *ServiceContext) ApplyConfigHotSwap(c *config.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.Config
	if old.InteractWord != s.Autointerract.InteractWord {
		c.InteractWord = old.InteractWord
	}
	if old.EntryEffect != s.Autointerract.EntryEffect {
		c.EntryEffect = old.EntryEffect
	}
	if old.WelcomeHighWealthy != s.Autointerract.WelcomeHighWealthy {
		c.WelcomeHighWealthy = old.WelcomeHighWealthy
	}
	s.Config = c
}

// ApplyRoomSwitch 换房时把新房 ctx 的字段按序应用到本实例（写侧加锁），
// 返回旧数据库连接由调用方延迟关闭。Autointerract 属于本实例状态，
// 不随 ctx 覆盖；欢迎开关以配置文件为准（新房间全新开始）。
func (s *ServiceContext) ApplyRoomSwitch(ctx *ServiceContext) *gorm.DB {
	s.mu.Lock()
	defer s.mu.Unlock()
	oldDb := s.db
	s.Config = ctx.Config
	s.UserID = ctx.UserID
	s.AnchorName = ctx.AnchorName
	s.RobotID = ctx.RobotID
	s.SignInModel = ctx.SignInModel
	s.DanmuCntModel = ctx.DanmuCntModel
	s.BlindBoxStatModel = ctx.BlindBoxStatModel
	s.VisitModel = ctx.VisitModel
	s.db = ctx.db
	return oldDb
}

// Snapshot 加读锁拷贝当前上下文（GetSvc 用，避免撕裂读）。
// 逐字段构造而非整体拷贝，避免复制锁值（go vet copylocks）。
func (s *ServiceContext) Snapshot() ServiceContext {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return ServiceContext{
		Config:            s.Config,
		db:                s.db,
		OtherSideUid:      s.OtherSideUid,
		SignInModel:       s.SignInModel,
		DanmuCntModel:     s.DanmuCntModel,
		BlindBoxStatModel: s.BlindBoxStatModel,
		VisitModel:        s.VisitModel,
		UserID:            s.UserID,
		AnchorName:        s.AnchorName,
		Autointerract:     s.Autointerract,
		RobotID:           s.RobotID,
	}
}

// PauseWelcome 天选/红包活动期间临时关闭欢迎弹幕
func (s *ServiceContext) PauseWelcome() {
	s.MutateConfig(func(c *config.Config) {
		c.InteractWord = false
		c.EntryEffect = false
		c.WelcomeHighWealthy = false
	})
}

// ResumeWelcome 活动结束后把欢迎弹幕恢复为主播设定的默认状态
func (s *ServiceContext) ResumeWelcome() {
	s.MutateConfig(func(c *config.Config) {
		if c.InteractWord != s.Autointerract.InteractWord {
			c.InteractWord = s.Autointerract.InteractWord
		}
		if c.EntryEffect != s.Autointerract.EntryEffect {
			c.EntryEffect = s.Autointerract.EntryEffect
		}
		if c.WelcomeHighWealthy != s.Autointerract.WelcomeHighWealthy {
			c.WelcomeHighWealthy = s.Autointerract.WelcomeHighWealthy
		}
	})
}

// SetWelcomeDefault 主播指令设置欢迎弹幕开关，同时更新恢复用的默认值
func (s *ServiceContext) SetWelcomeDefault(on bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nc := *s.Config
	nc.InteractWord = on
	nc.EntryEffect = on
	nc.WelcomeHighWealthy = on
	s.Config = &nc
	s.Autointerract.InteractWord = on
	s.Autointerract.EntryEffect = on
	s.Autointerract.WelcomeHighWealthy = on
}

type ServiceContext struct {
	// mu 序列化 reload 的字段替换与运行时配置变更（写侧互斥）。
	// 事件读取方（features 的 handler）不加锁：写侧采用「新 *Config 快照 +
	// 字段顺序替换」，64 位平台上指针/标量对齐写为原子，读侧最坏读到
	// 新旧混合的相邻值（如 Config 新 + Model 旧各一事件），不会撕裂崩溃；
	// 若需要严格 -race 干净，后续可在各 feature 事件入口包 RLock（约 15 处）。
	mu                sync.RWMutex
	Config            *config.Config
	db                *gorm.DB
	OtherSideUid      *SafeUidSet
	SignInModel       model.SignInModel
	DanmuCntModel     model.DanmuCntModel
	BlindBoxStatModel model.BlindBoxStatModel
	VisitModel        model.VisitModel
	UserID            int64  //主播id
	AnchorName        string //主播昵称
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
	s := &ServiceContext{
		OtherSideUid:      NewSafeUidSet(),
		SignInModel:       model.NewSignInModel(db, int64(c.RoomId)),
		DanmuCntModel:     model.NewDanmuCntModel(db, int64(c.RoomId)),
		BlindBoxStatModel: model.NewBlindBoxStatModel(db, int64(c.RoomId)),
		VisitModel:        model.NewVisitModel(db, int64(c.RoomId)),
		db:                db,
		Config:            &c,
		UserID:            0,
	}
	// 恢复基准取配置文件的初始值：天选/红包结束后 ResumeWelcome 回到主播
	// 配置的默认状态（临时关语义）。历史版本从未初始化该字段，导致恢复到
	// 零值 false、欢迎在首次天选后静默失效，属于存量缺陷，此处一并修复；
	// 主播用指令开关时 SetWelcomeDefault 会同步更新该基准
	s.Autointerract.InteractWord = c.InteractWord
	s.Autointerract.EntryEffect = c.EntryEffect
	s.Autointerract.WelcomeHighWealthy = c.WelcomeHighWealthy
	return s
}

// Close 释放持有的数据库连接；换房时对旧 context 调用，避免连接泄漏
func (s *ServiceContext) Close() error {
	if s.db == nil {
		return nil
	}
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
