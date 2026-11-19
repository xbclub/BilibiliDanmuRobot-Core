package danmuproc

import (
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// TestBlindBoxDayPrefix 日/月盲盒查询的称呼前缀
func TestBlindBoxDayPrefix(t *testing.T) {
	cfg := config.Config{UserNicknames: map[string]string{"42": "小坏"}}
	ctx := &svc.ServiceContext{Config: &cfg, UserID: 90001}

	// 熟人昵称：昵称是称呼不是 @，有无 reply 都保留
	if got := blindBoxDayPrefix("42", "断劍無我", ctx, true); got != "小坏 " {
		t.Fatalf("nickname prefix = %q", got)
	}
	// 普通观众 + 有 reply：真 @ 由 reply_mid 承担，文本里不再拼 @，
	// 否则客户端渲染出两个 @（真 @ + 纯文本 @）
	if got := blindBoxDayPrefix("43", "路人甲", ctx, true); got != "" {
		t.Fatalf("normal prefix with reply = %q", got)
	}
	// 普通观众 + 无 reply（无真 @ 通道）：落文本 @
	if got := blindBoxDayPrefix("43", "路人甲", ctx, false); got != "@路人甲 " {
		t.Fatalf("normal prefix without reply = %q", got)
	}
	// 主播自查：走全房间统计，不带前缀（避免全房盈亏归于主播名下）
	if got := blindBoxDayPrefix("90001", "主播", ctx, true); got != "" {
		t.Fatalf("anchor prefix = %q", got)
	}
}
