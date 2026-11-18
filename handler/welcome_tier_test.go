package handler

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

func newTierTestHandler(first bool, cfg config.Config) *wsHandler {
	return &wsHandler{svc: &svc.ServiceContext{
		Config:     &cfg,
		UserID:     90001,
		AnchorName: "家里云",
		VisitModel: &fakeVisitModel{first: first},
	}}
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

	if msg, ok := newTierTestHandler(true, cfg).tierWelcome(1, "小明", fan); !ok || msg != "老粉" {
		t.Fatalf("fan should win over first visit, got %q %v", msg, ok)
	}
	if msg, ok := newTierTestHandler(true, cfg).tierWelcome(2, "小明", nil); !ok || msg != "新人" {
		t.Fatalf("first visit should use new template, got %q %v", msg, ok)
	}
	if msg, ok := newTierTestHandler(false, cfg).tierWelcome(3, "小明", nil); !ok || msg != "回头" {
		t.Fatalf("repeat visit should use back template, got %q %v", msg, ok)
	}
	otherRoomFan := &entity.FansMedalInfo{TargetID: 123, MedalName: "别家", MedalLevel: 10}
	if msg, ok := newTierTestHandler(false, cfg).tierWelcome(4, "小明", otherRoomFan); !ok || msg != "回头" {
		t.Fatalf("a medal for another room should not count as fan, got %q %v", msg, ok)
	}
}

func TestTierWelcomeFallsBackWhenPoolsEmpty(t *testing.T) {
	cfg := config.Config{WelcomeTier: true, DanmuLen: 40}
	if _, ok := newTierTestHandler(true, cfg).tierWelcome(5, "小明", nil); ok {
		t.Fatal("empty tier pools should fall back to the legacy welcome path")
	}
}

func TestIsFan(t *testing.T) {
	w := newTierTestHandler(true, config.Config{WelcomeFanMedalLevel: 5})
	if w.isFan(nil) {
		t.Fatal("nil medal should not be a fan")
	}
	if w.isFan(&entity.FansMedalInfo{TargetID: 90001, MedalName: "家里云", MedalLevel: 4}) {
		t.Fatal("medal below the configured level should not be a fan")
	}
	if w.isFan(&entity.FansMedalInfo{TargetID: 90001, MedalName: "", MedalLevel: 10}) {
		t.Fatal("empty medal name should not be a fan")
	}
	if w.isFan(&entity.FansMedalInfo{TargetID: 90002, MedalName: "别家", MedalLevel: 10}) {
		t.Fatal("medal for another room should not be a fan")
	}
	if !w.isFan(&entity.FansMedalInfo{TargetID: 90001, MedalName: "家里云", MedalLevel: 5}) {
		t.Fatal("valid medal at the configured level should be a fan")
	}
}

func TestRenderWelcomeAtMode(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		Config:     &config.Config{WelcomeUseAt: true, DanmuLen: 40},
		AnchorName: "家里云",
	}
	got := renderWelcome("欢迎呀~{at}恭喜你发现了{anchor}", "测试", nil, svcCtx)
	if want := "欢迎呀~恭喜你发现了家里云"; got != want {
		t.Fatalf("renderWelcome() = %q, want %q", got, want)
	}
}

func TestRenderWelcomeTextMode(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		Config:     &config.Config{WelcomeUseAt: false, DanmuLen: 40},
		AnchorName: "家里云",
	}
	got := renderWelcome("哈喽{at}今天有没有更想{anchor}呀~", "测试", nil, svcCtx)
	if want := "哈喽@测试 今天有没有更想家里云呀~"; got != want {
		t.Fatalf("renderWelcome() = %q, want %q", got, want)
	}
}
