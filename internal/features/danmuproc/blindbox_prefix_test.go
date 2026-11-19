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

	// 熟人昵称：昵称替代 @，不带 @
	if got := blindBoxDayPrefix("42", "断劍無我", ctx); got != "小坏 " {
		t.Fatalf("nickname prefix = %q", got)
	}
	// 普通观众：@用户名
	if got := blindBoxDayPrefix("43", "路人甲", ctx); got != "@路人甲 " {
		t.Fatalf("normal prefix = %q", got)
	}
	// 主播自查：走全房间统计，不带前缀（避免全房盈亏归于主播名下）
	if got := blindBoxDayPrefix("90001", "主播", ctx); got != "" {
		t.Fatalf("anchor prefix = %q", got)
	}
}
