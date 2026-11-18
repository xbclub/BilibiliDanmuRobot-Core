package thanksgift

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

func loadReplayLines(t *testing.T, path string) []string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("打开回放数据失败: %v", err)
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		if s := strings.TrimSpace(sc.Text()); s != "" {
			lines = append(lines, s)
		}
	}
	return lines
}

func drainGiftChan() []*entity.SendGiftText {
	var out []*entity.SendGiftText
	for {
		select {
		case g := <-thanksGiver.giftChan:
			out = append(out, g)
		default:
			return out
		}
	}
}

// TestReplaySendGiftV2 用真实录制的 SEND_GIFT_V2 报文回放：
// pb 双形态回退 → 解码 → ToSendGiftText → 礼物感谢入队
func TestReplaySendGiftV2(t *testing.T) {
	cfg := config.Config{
		DBPath:     t.TempDir(),
		DBName:     "replay.db",
		ThanksGift: true,
		DanmuLen:   40,
	}
	s := svc.NewServiceContext(cfg)

	lines := loadReplayLines(t, "testdata/send_gift_v2.jsonl")
	if len(lines) == 0 {
		t.Fatal("回放数据为空")
	}

	for _, l := range lines {
		onSendGiftV2Raw(s, l)
	}
	got := drainGiftChan()
	if len(got) != 1 {
		t.Fatalf("应收到 1 份礼物, got %d", len(got))
	}
	g := got[0]
	if g.Data.Uname != "桃花依旧笑春风ペ" || g.Data.GiftName != "元气满格" ||
		g.Data.Num != 1 || g.Data.Price != 9000 || g.Data.TotalCoin != 9000 || g.Data.CoinType != "gold" {
		t.Fatalf("礼物字段与真实样本不符: %+v", g.Data)
	}
}
