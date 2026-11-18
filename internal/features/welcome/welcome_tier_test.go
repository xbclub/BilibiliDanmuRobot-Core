package welcome

import (
	"context"
	"testing"

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

	if got := renderWelcome(tpl, "小明", medal, ctx); got != "@小明 小明的云宝已7级，最爱家里云" {
		t.Fatalf("non-at render mismatch: %q", got)
	}

	ctx.Config.WelcomeUseAt = true
	if got := renderWelcome(tpl, "小明", medal, ctx); got != "小明的云宝已7级，最爱家里云" {
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
