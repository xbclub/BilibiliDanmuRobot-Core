package biliapi

import (
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// TestDryRunSend dry-run 开启时 Send 不触网、直接成功返回（空凭证下
// 真实发送必然失败，返回 nil 即证明未走 HTTP 路径）
func TestDryRunSend(t *testing.T) {
	EnableDryRun()
	defer DisableDryRun()
	if !DryRunEnabled() {
		t.Fatal("EnableDryRun 后应处于 dry-run 模式")
	}

	s := &svc.ServiceContext{Config: &config.Config{RoomId: 23058}}
	if err := Send("测试弹幕", s); err != nil {
		t.Fatalf("dry-run 下 Send 应返回 nil: %v", err)
	}
	if err := Send("测试@", s, &entity.DanmuMsgTextReplyInfo{ReplyUid: "42"}); err != nil {
		t.Fatalf("dry-run 下带回复的 Send 应返回 nil: %v", err)
	}

	DisableDryRun()
	if DryRunEnabled() {
		t.Fatal("DisableDryRun 后应退出 dry-run 模式")
	}
}
