package welcome

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// loadReplayLines 读取录制报文（danmu-recoder recoder-logs 真实线上事件）
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
	if err := sc.Err(); err != nil {
		t.Fatalf("读取回放数据失败: %v", err)
	}
	return lines
}

func drainInteractChan() []*InterractData {
	var out []*InterractData
	for {
		select {
		case g := <-interractGiver.interractChan:
			out = append(out, g)
		default:
			return out
		}
	}
}

func resetEnterSeen() {
	enterSeenMu.Lock()
	enterSeen = map[int64]time.Time{}
	enterSeenMu.Unlock()
}

// TestReplayInteractWordV2 用真实录制的 INTERACT_WORD_V2 报文回放全链路：
// pb 解码 → 三档分流 → visit 表首访/回访 → 事件层去重
func TestReplayInteractWordV2(t *testing.T) {
	cfg := config.Config{
		DBPath:       t.TempDir(),
		DBName:       "replay.db",
		InteractWord: true,
		WelcomeTier:  true,
		DanmuLen:     40,
		WelcomeNew:   []string{"新人驾到{at}"},
		WelcomeBack:  []string{"又来啦{at}"},
		WelcomeFan:   []string{"老粉{at}"},
	}
	s := svc.NewServiceContext(cfg)

	lines := loadReplayLines(t, "testdata/interact_word_v2.jsonl")
	if len(lines) < 2 {
		t.Fatalf("回放数据不足: %d 条", len(lines))
	}

	// 第一遍：两位用户均为首访（录制样本无粉丝牌）→ 新人档
	for _, l := range lines {
		onInteractWordV2Raw(s, l)
	}
	msgs := drainInteractChan()
	if len(msgs) != len(lines) {
		t.Fatalf("首访应产生 %d 条欢迎, got %d", len(lines), len(msgs))
	}
	for _, m := range msgs {
		if !strings.Contains(m.Msg, "新人驾到") {
			t.Fatalf("首访欢迎应为新人档, got %q", m.Msg)
		}
	}

	// 原样重放：事件层去重应全部拦截，不产生新欢迎
	for _, l := range lines {
		onInteractWordV2Raw(s, l)
	}
	if n := len(drainInteractChan()); n != 0 {
		t.Fatalf("10 秒内重复进场应被去重, got %d 条", n)
	}

	// 清除去重戳后重放：visit 表已记录 → 回头客档
	resetEnterSeen()
	for _, l := range lines {
		onInteractWordV2Raw(s, l)
	}
	msgs = drainInteractChan()
	if len(msgs) != len(lines) {
		t.Fatalf("二次进场应产生 %d 条欢迎, got %d", len(lines), len(msgs))
	}
	for _, m := range msgs {
		if !strings.Contains(m.Msg, "又来啦") {
			t.Fatalf("二次进场应为回头客档, got %q", m.Msg)
		}
	}
}
