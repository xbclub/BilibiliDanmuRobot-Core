package welcome

import (
	"strings"
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// TestWelcomeTemplateNoSplit 真实模板「欢迎 {user}, 来看刀刀酱~」的拆分边界：
// \n 超长回退仅在渲染后长度 > DanmuLen 时触发，B 站 40 上限下任何正常
// 长度用户名（B 站 ≤16 字）都不应再被拆成两条
func TestWelcomeTemplateNoSplit(t *testing.T) {
	const tpl = "欢迎 {user}, 来看刀刀酱~"
	newCtx := func(dl int) *svc.ServiceContext {
		cfg := config.Config{DBPath: t.TempDir(), DBName: "t.db", DanmuLen: dl,
			WelcomeDanmu: []string{tpl}, InteractWord: true}
		return svc.NewServiceContext(cfg)
	}
	rendered := strings.ReplaceAll(tpl, "{user}", "超级无敌酷炫大蟑螂")

	t.Run("40_原样单条", func(t *testing.T) {
		if got := handleInterract(12345, "超级无敌酷炫大蟑螂", newCtx(40)); got != rendered {
			t.Fatalf("got=%q want=%q", got, rendered)
		}
	})
	t.Run("20_恰好等于上限_仍单条", func(t *testing.T) {
		if got := handleInterract(12345, "超级无敌酷炫大蟑螂", newCtx(20)); got != rendered {
			t.Fatalf("got=%q want=%q", got, rendered)
		}
	})
	t.Run("19_触发回退_双换行且逗号被吃", func(t *testing.T) {
		got := handleInterract(12345, "超级无敌酷炫大蟑螂", newCtx(19))
		if !strings.Contains(got, "\n") || strings.Contains(got, ",") {
			t.Fatalf("超长回退应含 \\n 去逗号: got=%q", got)
		}
	})
	t.Run("10字名_40_不触发回退", func(t *testing.T) {
		got := handleInterract(12345, "月下白糖_棘背哈氣版", newCtx(40))
		want := strings.ReplaceAll(tpl, "{user}", "月下白糖_棘背哈氣版")
		if got != want {
			t.Fatalf("got=%q want=%q", got, want)
		}
	})
}
