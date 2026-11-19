package welcome

import (
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// TestWelcomePathsConsistent 主欢迎、分时段欢迎、特效欢迎三条路径
// 必须产出相同的格式化结果——它们共用 renderEnterWelcome，
// 此前各写一份导致行为漂移（特效无条件拆条、主路径双换行）
func TestWelcomePathsConsistent(t *testing.T) {
	const tpl = "欢迎 {user}, 来看刀刀酱~"
	const name = "超级无敌酷炫大蟑螂"
	want := "欢迎 超级无敌酷炫大蟑螂, 来看刀刀酱~"

	// 比较两条路径的产出（主池路径用 handleInterract，特效路径用 getRandomWelcome）
	base := config.Config{
		DBPath: t.TempDir(), DBName: "t.db", DanmuLen: 40,
		WelcomeDanmu: []string{tpl},
	}
	ctx := svc.NewServiceContext(base)
	main := handleInterract(1, name, ctx)
	effect := getRandomWelcome(name, ctx)
	if main != want || effect != want {
		t.Fatalf("主路径=%q 特效路径=%q, 均应为 %q", main, effect, want)
	}
	if main != effect {
		t.Fatalf("两条路径不一致: main=%q effect=%q", main, effect)
	}
}

// TestRenderEnterWelcomeShared 统一渲染函数的边界行为
func TestRenderEnterWelcomeShared(t *testing.T) {
	mk := func(dl int, useAt bool) *svc.ServiceContext {
		cfg := config.Config{DBPath: t.TempDir(), DBName: "t.db", DanmuLen: dl, WelcomeUseAt: useAt}
		return svc.NewServiceContext(cfg)
	}
	const tpl = "欢迎 {user}, 来看刀刀酱~"
	const name = "小蟑螂"

	// 未超长：原样单条，标点完整
	if got := renderEnterWelcome(tpl, name, mk(40, false)); got != "欢迎 小蟑螂, 来看刀刀酱~" {
		t.Fatalf("未超长 = %q", got)
	}
	// 带 @ 档：用户名不进正文，逗号收敛
	if got := renderEnterWelcome(tpl, name, mk(40, true)); got != "欢迎，来看刀刀酱~" {
		t.Fatalf("带@档 = %q", got)
	}
	// 超长：在用户名后换行，且只有一个换行（旧实现会插两个）
	long := "超级无敌酷炫大蟑螂啊啊啊啊啊啊"
	got := renderEnterWelcome(tpl, long, mk(10, false))
	if n := len(splitByNewline(got)); n != 2 {
		t.Fatalf("超长应恰好 2 段（单个换行），got %d 段: %q", n, got)
	}
}

func splitByNewline(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == '\n' {
			out = append(out, cur)
			cur = ""
			continue
		}
		cur += string(r)
	}
	return append(out, cur)
}
