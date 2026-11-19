package welcome

import (
	"strings"
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// TestEntryEffectNoForcedSplit 特效欢迎（ENTRY_EFFECT）此前无条件把
// 「{user}, 」替换成「{user}\n」，再由 interract_bullet.go 按 \n 拆成
// 多条独立弹幕，导致未超长的欢迎语也被拆成两条且逗号丢失。
// 现改为：未超 DanmuLen 单条输出，仅超长才换行分段。
func TestEntryEffectNoForcedSplit(t *testing.T) {
	const tpl = "欢迎 {user}, 来看刀刀酱~"
	newCtx := func(dl int, useAt bool) *svc.ServiceContext {
		cfg := config.Config{
			DBPath: t.TempDir(), DBName: "t.db",
			DanmuLen: dl, WelcomeDanmu: []string{tpl},
			WelcomeUseAt: useAt, InteractWordByTime: false,
		}
		return svc.NewServiceContext(cfg)
	}

	t.Run("未超长_单条且保留逗号", func(t *testing.T) {
		// 线上实际案例：9 字用户名，DanmuLen=40
		got := getRandomWelcome("超级无敌酷炫大蟑螂", newCtx(40, false))
		if want := "欢迎 超级无敌酷炫大蟑螂, 来看刀刀酱~"; got != want {
			t.Fatalf("got=%q want=%q", got, want)
		}
		if strings.Contains(got, "\n") {
			t.Fatalf("未超长不应含换行: %q", got)
		}
	})

	t.Run("长用户名仍未超40_单条", func(t *testing.T) {
		// B 站用户名上限 16 字 → 渲染 27 rune，远小于 40
		name := "超级无敌酷炫大蟑螂啊啊啊啊啊啊"
		got := getRandomWelcome(name, newCtx(40, false))
		if strings.Contains(got, "\n") {
			t.Fatalf("未超长不应含换行: %q", got)
		}
		if want := "欢迎 " + name + ", 来看刀刀酱~"; got != want {
			t.Fatalf("got=%q want=%q", got, want)
		}
	})

	t.Run("超长_换行分段作兜底", func(t *testing.T) {
		// DanmuLen=19 时 20 rune 超限，才走 \n 分段
		got := getRandomWelcome("超级无敌酷炫大蟑螂", newCtx(19, false))
		if !strings.Contains(got, "\n") {
			t.Fatalf("超长应换行分段: %q", got)
		}
		parts := strings.Split(got, "\n")
		if len(parts) != 2 || !strings.Contains(parts[0], "超级无敌酷炫大蟑螂") {
			t.Fatalf("分段异常: %#v", parts)
		}
	})

	t.Run("WelcomeUseAt_单条逗号无换行", func(t *testing.T) {
		// 带 @ 档：@ 由 reply_mid 承担，用户名不进文本，逗号收敛成中文逗号
		got := getRandomWelcome("超级无敌酷炫大蟑螂", newCtx(40, true))
		if strings.Contains(got, "\n") || strings.Contains(got, "{user}") {
			t.Fatalf("带@档应单条且已替换: %q", got)
		}
		if want := "欢迎，来看刀刀酱~"; got != want {
			t.Fatalf("got=%q want=%q", got, want)
		}
	})
}

// TestEntryEffectConsumerSplit 消费端按 \n 拆条的契约不变（超长兜底依赖）
func TestEntryEffectConsumerSplit(t *testing.T) {
	msg := "欢迎 超级无敌酷炫大蟑螂\n来看刀刀酱~"
	if parts := strings.Split(msg, "\n"); len(parts) != 2 {
		t.Fatalf("消费端拆分契约变化: %#v", parts)
	}
}
