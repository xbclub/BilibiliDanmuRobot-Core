package danmuproc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// TestReplayDanmuMsg 用真实录制的 DANMU_MSG 报文回放解析分发：
// 全量过 processDanmuMessage 不 panic，且 uid/文本结构可正常提取
func TestReplayDanmuMsg(t *testing.T) {
	cfg := config.Config{
		DBPath: t.TempDir(),
		DBName: "replay.db",
		// 功能开关全关：回放只验证解析与分发路径，不触网
	}
	s := svc.NewServiceContext(cfg)

	f, err := os.Open("testdata/danmu_msg.jsonl")
	if err != nil {
		t.Fatalf("打开回放数据失败: %v", err)
	}
	defer f.Close()

	total, parsed := 0, 0
	firstUid, firstText := "", ""
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		total++
		// 真实报文喂完整处理链（解析/分发/日志），任何 panic 都会让测试失败
		processDanmuMessage(line, s)

		// 结构断言：Info[1] 文本与 Info[2] uid 可提取
		var dm entity.DanmuMsgText
		if err := json.Unmarshal([]byte(line), &dm); err != nil {
			continue
		}
		if text, ok := dm.Info[1].(string); ok && text != "" {
			parsed++
			if firstText == "" {
				firstText = text
				if from, ok := dm.Info[2].([]interface{}); ok && len(from) > 0 {
					firstUid = fmt.Sprintf("%.0f", from[0].(float64))
				}
			}
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("读取回放数据失败: %v", err)
	}
	if total == 0 {
		t.Fatal("回放数据为空")
	}
	if parsed != total {
		t.Fatalf("全部报文都应可提取弹幕文本: %d/%d", parsed, total)
	}
	if firstUid == "" || firstText == "" {
		t.Fatalf("首条弹幕应可提取 uid 与文本: uid=%q text=%q", firstUid, firstText)
	}
}
