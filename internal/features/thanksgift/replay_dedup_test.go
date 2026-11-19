package thanksgift

import (
	"encoding/json"
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// newReplayCtx 构造回放用 ServiceContext（独立临时 DB）
func newReplayCtx(t *testing.T, statEnabled bool) *svc.ServiceContext {
	t.Helper()
	cfg := config.Config{
		DBPath:                 t.TempDir(),
		DBName:                 "replay.db",
		ThanksGift:             true,
		DanmuLen:               40,
		BlindBoxProfitLossStat: statEnabled,
	}
	return svc.NewServiceContext(cfg)
}

// decodeReplayV2 解析一条录制的 SEND_GIFT_V2 原始报文（回放测试共用）
func decodeReplayV2(t *testing.T, rawLine string) *entity.SendGiftV2 {
	t.Helper()
	var msg entity.SendGiftV2Text
	if err := json.Unmarshal([]byte(rawLine), &msg); err != nil {
		t.Fatalf("V2 envelope 解析失败: %v", err)
	}
	pb := msg.Data.Pb
	if pb == "" {
		pb = msg.Data.Data.Pb
	}
	g, err := entity.ParseSendGiftV2(pb)
	if err != nil {
		t.Fatalf("V2 pb 解析失败: %v", err)
	}
	return g
}

// TestReplayV1V2TidDedup 过渡期同一礼物经 V1(JSON) 与 V2(pb) 双下发、
// 交易号 tid 相同：无论到达顺序，聚合感谢只计一次（giftChan 单条）。
// V1 报文由真实 V2 样本字段合成（tid 等字段即双事件共享的交易号）。
func TestReplayV1V2TidDedup(t *testing.T) {
	t.Cleanup(resetGiftTidSeen)
	lines := loadReplayLines(t, "testdata/send_gift_v2_room254.jsonl")

	syntheticV1 := func(rawLine string) string {
		g := decodeReplayV2(t, rawLine)
		send := g.ToSendGiftText(g.Gifts[0])
		b, err := json.Marshal(send)
		if err != nil {
			t.Fatalf("合成 V1 报文失败: %v", err)
		}
		return string(b)
	}

	assertCount := func(t *testing.T, want int) {
		t.Helper()
		got := drainGiftChan()
		if len(got) != want {
			t.Fatalf("giftChan 应 %d 条, got %d", want, len(got))
		}
	}

	t.Run("V1先到_V2后到_只计一次", func(t *testing.T) {
		t.Cleanup(resetGiftTidSeen)
		s := newReplayCtx(t, true)
		onSendGiftV1Raw(s, syntheticV1(lines[0]))
		onSendGiftV2Raw(s, lines[0])
		assertCount(t, 1)
	})

	t.Run("V2先到_V1后到_只计一次", func(t *testing.T) {
		t.Cleanup(resetGiftTidSeen)
		s := newReplayCtx(t, true)
		onSendGiftV2Raw(s, lines[0])
		onSendGiftV1Raw(s, syntheticV1(lines[0]))
		assertCount(t, 1)
	})

	t.Run("不同tid互不误杀", func(t *testing.T) {
		t.Cleanup(resetGiftTidSeen)
		s := newReplayCtx(t, true)
		onSendGiftV2Raw(s, lines[0])
		onSendGiftV2Raw(s, lines[1])
		assertCount(t, 2)
	})

	t.Run("batch四礼物逐项去重_合成V1只补零", func(t *testing.T) {
		t.Cleanup(resetGiftTidSeen)
		s := newReplayCtx(t, true)
		onSendGiftV2Raw(s, lines[3]) // 真实 4 礼物 batch，4 个不同 tid
		assertCount(t, 4)
		// batch 首个礼物的同 tid V1 追加到达：不再计入
		onSendGiftV1Raw(s, syntheticV1(lines[3]))
		assertCount(t, 0)
	})
}

// TestReplayV2BlindBoxReal 真实 V2 盲盒样本（阿狗盲盒→温柔出击，房间 25470981）：
// 字段解码、聚合/即时两档行为、盈亏话术
func TestReplayV2BlindBoxReal(t *testing.T) {
	t.Cleanup(resetGiftTidSeen)
	lines := loadReplayLines(t, "testdata/send_gift_v2_room254.jsonl")

	t.Run("字段与默认聚合档", func(t *testing.T) {
		t.Cleanup(resetGiftTidSeen)
		s := newReplayCtx(t, true) // 盲盒盈亏统计开（默认）
		onSendGiftV2Raw(s, lines[2])
		got := drainGiftChan()
		if len(got) != 1 {
			t.Fatalf("应 1 条, got %d", len(got))
		}
		d := got[0].Data
		if d.Uname != "あすひM路" || d.GiftName != "温柔出击" || d.Num != 1 || d.Price != 9000 {
			t.Fatalf("盲盒字段与真实样本不符: %+v", d)
		}
		if d.BlindGift.OriginalGiftName != "阿狗盲盒" || d.BlindGift.OriginalGiftPrice != 9000 {
			t.Fatalf("盲盒原始礼物字段不符: %+v", d.BlindGift)
		}
		if got[0].SkipThanks {
			t.Fatal("默认聚合档不应跳过聚合感谢")
		}
		// 同价开出：没赚没亏
		if got0 := blindBoxGiftThanks(got[0]); got0 != "感谢あすひM路赠送的阿狗盲盒，爆出温柔出击，没赚没亏" {
			t.Fatalf("盲盒话术不符: %q", got0)
		}
	})

	t.Run("统计关闭走即时话术", func(t *testing.T) {
		t.Cleanup(resetGiftTidSeen)
		s := newReplayCtx(t, false)
		onSendGiftV2Raw(s, lines[2])
		got := drainGiftChan()
		if len(got) != 1 {
			t.Fatalf("应 1 条, got %d", len(got))
		}
		if !got[0].SkipThanks {
			t.Fatal("即时话术档应置 SkipThanks 跳过聚合")
		}
	})
}
