package welcome

import (
	"context"
	"testing"
	"time"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

type fakeVisitModel struct {
	first bool
}

func (f *fakeVisitModel) Touch(context.Context, int64, int64) (bool, error) {
	return f.first, nil
}

func newTierTestCtx(first bool, cfg config.Config) *svc.ServiceContext {
	return &svc.ServiceContext{
		Config:     &cfg,
		UserID:     90001,
		AnchorName: "家里云",
		VisitModel: &fakeVisitModel{first: first},
	}
}

func TestTierWelcomePriority(t *testing.T) {
	cfg := config.Config{
		WelcomeTier:          true,
		DanmuLen:             40,
		WelcomeNew:           []string{"新人"},
		WelcomeBack:          []string{"回头"},
		WelcomeFan:           []string{"老粉"},
		WelcomeFanMedalLevel: 0,
	}
	fan := &entity.FansMedalInfo{TargetID: 90001, MedalName: "家里云", MedalLevel: 10}

	if msg, ok := tierWelcome(newTierTestCtx(true, cfg), 1, "小明", fan); !ok || msg != "老粉" {
		t.Fatalf("fan should win over first visit, got %q %v", msg, ok)
	}
	if msg, ok := tierWelcome(newTierTestCtx(true, cfg), 2, "小明", nil); !ok || msg != "新人" {
		t.Fatalf("first visit should use new template, got %q %v", msg, ok)
	}
	if msg, ok := tierWelcome(newTierTestCtx(false, cfg), 3, "小明", nil); !ok || msg != "回头" {
		t.Fatalf("repeat visit should use back template, got %q %v", msg, ok)
	}
	otherRoomFan := &entity.FansMedalInfo{TargetID: 123, MedalName: "别家", MedalLevel: 10}
	if msg, ok := tierWelcome(newTierTestCtx(false, cfg), 4, "小明", otherRoomFan); !ok || msg != "回头" {
		t.Fatalf("a medal for another room should not count as fan, got %q %v", msg, ok)
	}
}

type errVisitModel struct{}

func (errVisitModel) Touch(context.Context, int64, int64) (bool, error) {
	return false, context.DeadlineExceeded
}

func TestTierWelcomeFallsBackWhenPoolsEmpty(t *testing.T) {
	cfg := config.Config{WelcomeTier: true, DanmuLen: 40}
	if _, ok := tierWelcome(newTierTestCtx(true, cfg), 5, "小明", nil); ok {
		t.Fatal("empty tier pools should fall back to the legacy welcome path")
	}
}

func TestTierWelcomeFanPoolEmptyFallsToNew(t *testing.T) {
	cfg := config.Config{
		WelcomeTier: true,
		DanmuLen:    40,
		WelcomeNew:  []string{"新人"},
	}
	fan := &entity.FansMedalInfo{TargetID: 90001, MedalName: "家里云", MedalLevel: 10}
	if msg, ok := tierWelcome(newTierTestCtx(true, cfg), 10, "小明", fan); !ok || msg != "新人" {
		t.Fatalf("empty fan pool should fall to the new tier, got %q %v", msg, ok)
	}
}

func TestTierWelcomeTouchErrorFallsBack(t *testing.T) {
	cfg := config.Config{WelcomeTier: true, DanmuLen: 40, WelcomeNew: []string{"新人"}}
	ctx := &svc.ServiceContext{
		Config:     &cfg,
		UserID:     90001,
		VisitModel: errVisitModel{},
	}
	if _, ok := tierWelcome(ctx, 11, "小明", nil); ok {
		t.Fatal("touch error should fall back to the legacy welcome path")
	}
}

func TestTierWelcomeNilVisitModelFallsBack(t *testing.T) {
	cfg := config.Config{WelcomeTier: true, DanmuLen: 40, WelcomeNew: []string{"新人"}}
	ctx := &svc.ServiceContext{Config: &cfg, UserID: 90001}
	if _, ok := tierWelcome(ctx, 12, "小明", nil); ok {
		t.Fatal("nil VisitModel should fall back to the legacy welcome path")
	}
}

func TestRenderWelcomePlaceholders(t *testing.T) {
	cfg := config.Config{DanmuLen: 40}
	ctx := &svc.ServiceContext{Config: &cfg, AnchorName: "家里云"}
	medal := &entity.FansMedalInfo{MedalName: "云宝", MedalLevel: 7}
	tpl := "{at}{user}的{medal}已{medal_level}级，最爱{anchor}"

	if got := renderWelcome(tpl, "小明", medal, ctx, 0); got != "@小明 小明的云宝已7级，最爱家里云" {
		t.Fatalf("non-at render mismatch: %q", got)
	}

	ctx.Config.WelcomeUseAt = true
	if got := renderWelcome(tpl, "小明", medal, ctx, 0); got != "小明的云宝已7级，最爱家里云" {
		t.Fatalf("at-mode render mismatch: %q", got)
	}
}

func TestDedupEnter(t *testing.T) {
	if !dedupEnter(42000001) {
		t.Fatal("first enter should pass")
	}
	if dedupEnter(42000001) {
		t.Fatal("duplicate within window should be dropped")
	}
	if !dedupEnter(42000002) {
		t.Fatal("different uid should pass")
	}
}

func TestDedupEnterExpiry(t *testing.T) {
	t0 := time.Now()
	if !dedupEnterAt(43000001, t0) {
		t.Fatal("first enter should pass")
	}
	if dedupEnterAt(43000001, t0.Add(9*time.Second)) {
		t.Fatal("9s duplicate should be dropped")
	}
	// 窗口过后同 uid 可再次通过，且过期条目被顺带清理
	if !dedupEnterAt(43000001, t0.Add(11*time.Second)) {
		t.Fatal("enter after window should pass")
	}
	enterSeenMu.Lock()
	if _, ok := enterSeen[43000002]; ok {
		enterSeenMu.Unlock()
		t.Fatal("expired entries should be cleaned")
	}
	enterSeenMu.Unlock()
}

func TestRenderWelcomeSmallBudget(t *testing.T) {
	cfg := config.Config{DanmuLen: 2}
	ctx := &svc.ServiceContext{Config: &cfg, AnchorName: "家里云云云云云"}
	medal := &entity.FansMedalInfo{MedalName: "云宝牌名很长", MedalLevel: 7}
	// DanmuLen=2 时预算保底 4：{anchor}/{medal} 至多 3 rune + 省略号
	got := renderWelcome("{anchor}|{medal}|{medal_level}", "小明", medal, ctx, 0)
	if got != "家里云…|云宝牌…|7" {
		t.Fatalf("small budget render mismatch: %q", got)
	}

	// nil medal：占位符渲染为空串（new/back 档模板不应包含 medal 占位符）
	got = renderWelcome("{user}来了", "小明", nil, ctx, 0)
	if got != "小明来了" {
		t.Fatalf("nil medal render mismatch: %q", got)
	}
}

func TestRenderWelcomeAlias(t *testing.T) {
	cfg := config.Config{
		DanmuLen:      40,
		UserNicknames: map[string]string{"90001": "小坏"},
	}
	ctx := &svc.ServiceContext{Config: &cfg, AnchorName: "家里云"}
	got := renderWelcome("哈喽{at}今天有没有更想{anchor}呀~", "小明", nil, ctx, 90001)
	if want := "哈喽小坏今天有没有更想家里云呀~"; got != want {
		t.Fatalf("renderWelcome() = %q, want %q", got, want)
	}

	// 超长昵称同样走长度预算截断：DanmuLen=8 → 预算保底 4，至多 3 rune + 省略号
	small := config.Config{DanmuLen: 8, UserNicknames: map[string]string{"90002": "这是一个特别长的昵称"}}
	smallCtx := &svc.ServiceContext{Config: &small, AnchorName: "家里云"}
	if got := renderWelcome("{at}来啦", "小明", nil, smallCtx, 90002); got != "这是一…来啦" {
		t.Fatalf("long alias should be clipped by budget, got %q", got)
	}
}
