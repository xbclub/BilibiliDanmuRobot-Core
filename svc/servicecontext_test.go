package svc

import (
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
)

func newWelcomeTestCtx(t *testing.T) *ServiceContext {
	t.Helper()
	c := config.Config{
		DBPath:             t.TempDir(),
		DBName:             "test.db",
		InteractWord:       true,
		EntryEffect:        true,
		WelcomeHighWealthy: true,
		LotteryEnable:      true,
	}
	return NewServiceContext(c)
}

// TestTempOffSemantics 锁定「欢迎弹幕临时关」全链路语义：
// 天选/红包暂停 → 恢复到配置默认；主播指令覆盖默认；热更保留临时关
func TestTempOffSemantics(t *testing.T) {
	s := newWelcomeTestCtx(t)

	// 1. 暂停：三开关全关
	s.PauseWelcome()
	if s.Config.InteractWord || s.Config.EntryEffect || s.Config.WelcomeHighWealthy {
		t.Fatal("PauseWelcome 后三开关应全为 false")
	}
	// 2. 恢复：回到配置默认 true（修复存量缺陷：此前恢复到未初始化的 false）
	s.ResumeWelcome()
	if !s.Config.InteractWord || !s.Config.EntryEffect || !s.Config.WelcomeHighWealthy {
		t.Fatal("ResumeWelcome 后应恢复配置默认 true")
	}

	// 3. 红包组合：PauseWelcome + LotteryEnable=false（等价旧版原地连写）
	s.PauseWelcome()
	s.MutateConfig(func(c *config.Config) { c.LotteryEnable = false })
	if s.Config.LotteryEnable || s.Config.InteractWord {
		t.Fatal("红包暂停态应为 LotteryEnable=false 且欢迎保持关闭")
	}
	s.ResumeWelcome()
	if !s.Config.InteractWord {
		t.Fatal("红包结束应恢复欢迎")
	}
	if s.Config.LotteryEnable {
		t.Fatal("ResumeWelcome 不应恢复 LotteryEnable（历史行为：仅 reload 恢复）")
	}

	// 4. 主播指令覆盖默认：关闭后，暂停→恢复应停在 false
	s.SetWelcomeDefault(false)
	s.PauseWelcome()
	s.ResumeWelcome()
	if s.Config.InteractWord {
		t.Fatal("指令关闭后恢复应保持 false")
	}
	// 指令重新开启
	s.SetWelcomeDefault(true)
	s.PauseWelcome()
	s.ResumeWelcome()
	if !s.Config.InteractWord {
		t.Fatal("指令开启后恢复应为 true")
	}

	// 5. 热更保留临时关；LotteryEnable 随热更回到文件值（与历史一致）
	s.PauseWelcome()
	fresh := config.Config{
		InteractWord:       true,
		EntryEffect:        true,
		WelcomeHighWealthy: true,
		LotteryEnable:      true,
	}
	s.ApplyConfigHotSwap(&fresh)
	if s.Config.InteractWord {
		t.Fatal("热更应保留天选/红包期间的临时关状态")
	}
	if !s.Config.LotteryEnable {
		t.Fatal("LotteryEnable 应随热更恢复文件值")
	}
	s.ResumeWelcome()
	if !s.Config.InteractWord {
		t.Fatal("热更后仍应能正确恢复欢迎")
	}
}
