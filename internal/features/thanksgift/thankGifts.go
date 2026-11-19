package thanksgift

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/internal/features/danmuproc"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

// registerGifts 礼物感谢
func registerGifts(c *client.Client, svcCtx *svc.ServiceContext) {
	c.RegisterCustomEventHandler("SEND_GIFT", func(s string) {
		onSendGiftV1Raw(svcCtx, s)
	})
	c.RegisterCustomEventHandler("SEND_GIFT_V2", func(s string) {
		onSendGiftV2Raw(svcCtx, s)
	})
	c.RegisterCustomEventHandler("GUARD_BUY", func(s string) {
		svcCtx.EventRLock()
		defer svcCtx.EventRUnlock()
		if svcCtx.Config.ThanksGift {
			send := &entity.GuardBuyText{}
			_ = json.Unmarshal([]byte(s), send)
			if svcCtx.Config.ThanksGiftUseAt {
				PushToGuardChan(svcCtx, send, &entity.DanmuMsgTextReplyInfo{
					ReplyUid: strconv.Itoa(send.Data.Uid),
				})
			} else {
				PushToGuardChan(svcCtx, send)
			}
		}
	})

	c.RegisterCustomEventHandler("COMMON_NOTICE_DANMAKU", func(s string) {
		svcCtx.EventRLock()
		defer svcCtx.EventRUnlock()
		if svcCtx.Config.ThanksGift {
			data := &entity.CommonNoticeDanmaku{}
			_ = json.Unmarshal([]byte(s), data)
			if len(data.Data.ContentSegments) == 5 &&
				data.Data.ContentSegments[1].Text == "投喂" &&
				data.Data.ContentSegments[2].Text == "大航海盲盒" {

				PushThanksBlindBox(svcCtx, data.Data.ContentSegments[0].Text, data.Data.ContentSegments[4].Text)
			} else if len(data.Data.ContentSegments) == 6 &&
				data.Data.ContentSegments[2].Text == "投喂" &&
				data.Data.ContentSegments[3].Text == "大航海盲盒" {

				PushThanksBlindBox(svcCtx, data.Data.ContentSegments[1].Text, data.Data.ContentSegments[5].Text)
			}
		}
	})
}

// onSendGiftV1Raw 处理 SEND_GIFT（V1 JSON）原始报文，
// 抽出以便录制数据回放测试（与 V2 路径对称）
func onSendGiftV1Raw(svcCtx *svc.ServiceContext, s string) {
	svcCtx.EventRLock()
	defer svcCtx.EventRUnlock()
	send := &entity.SendGiftText{}
	if err := json.Unmarshal([]byte(s), send); err != nil {
		logx.Errorf("SEND_GIFT 解析失败: %v", err)
		return
	}
	// V1/V2 双事件 tid 去重：同一礼物若已作为 SEND_GIFT_V2 处理过，
	// 跳过聚合与统计，避免双计数
	if giftTidDuplicate(send.Data.Tid) {
		logx.Debugf("礼物 tid 短窗去重: tid=%s gift=%s", send.Data.Tid, send.Data.GiftName)
		return
	}
	if svcCtx.Config.ThanksGift {
		logx.Infof("收到礼物 uname=%s gift=%s num=%d price=%d total=%d",
			send.Data.Uname, send.Data.GiftName, send.Data.Num, send.Data.Price, send.Data.TotalCoin)
		// 盲盒盈亏统计开启（默认）走旧聚合感谢；关闭时才用即时话术并
		// 跳过聚合，避免默认配置的聚合播报被静默替换（评审 P1-2）
		if send.Data.BlindGift.OriginalGiftName != "" && !svcCtx.Config.BlindBoxProfitLossStat {
			pushBlindGiftThanks(svcCtx, send)
			send.SkipThanks = true
		}
		PushToGiftChan(send)
	}
	danmuproc.SaveBlindBoxStat(send, svcCtx)
}

// onSendGiftV2Raw 处理 SEND_GIFT_V2 原始报文（pb 双形态回退 + 解码 + 入队），
// 抽出以便录制数据回放测试
func onSendGiftV2Raw(svcCtx *svc.ServiceContext, s string) {
	svcCtx.EventRLock()
	defer svcCtx.EventRUnlock()
	var msg entity.SendGiftV2Text
	if err := json.Unmarshal([]byte(s), &msg); err != nil {
		logx.Errorf("SEND_GIFT_V2 解析失败: %v", err)
		return
	}
	pb := msg.Data.Pb
	if pb == "" {
		pb = msg.Data.Data.Pb
	}
	if pb == "" {
		logx.Info("SEND_GIFT_V2 没有 pb 字段")
		return
	}
	gift, err := entity.ParseSendGiftV2(pb)
	if err != nil {
		logx.Errorf("SEND_GIFT_V2 pb 解析失败: %v", err)
		return
	}
	if len(gift.Gifts) == 0 {
		logx.Infof("SEND_GIFT_V2 未解析到礼物 uid=%d uname=%s", gift.UID, gift.Uname)
		return
	}
	for _, item := range gift.Gifts {
		send := gift.ToSendGiftText(item)
		// V1/V2 双事件 tid 去重：同一礼物若已作为 SEND_GIFT 处理过，
		// 跳过聚合与统计，避免双计数
		if giftTidDuplicate(send.Data.Tid) {
			logx.Debugf("礼物(V2) tid 短窗去重: tid=%s gift=%s", send.Data.Tid, send.Data.GiftName)
			continue
		}
		if svcCtx.Config.ThanksGift {
			logx.Infof("收到礼物(V2) uname=%s gift=%s num=%d price=%d total=%d",
				send.Data.Uname, send.Data.GiftName, send.Data.Num, send.Data.Price, send.Data.TotalCoin)
			// 盲盒盈亏统计开启（默认）走旧聚合感谢；关闭时才用即时话术并
			// 跳过聚合，避免默认配置的聚合播报被静默替换（评审 P1-2）
			if send.Data.BlindGift.OriginalGiftName != "" && !svcCtx.Config.BlindBoxProfitLossStat {
				pushBlindGiftThanks(svcCtx, send)
				send.SkipThanks = true
			}
			PushToGiftChan(send)
		}
		danmuproc.SaveBlindBoxStat(send, svcCtx)
	}
}

func pushBlindGiftThanks(svcCtx *svc.ServiceContext, send *entity.SendGiftText) {
	msg := blindBoxGiftThanks(send)
	if svcCtx.Config.ThanksGiftUseAt {
		logic.PushToBulletSender(msg, &entity.DanmuMsgTextReplyInfo{
			ReplyUid: strconv.Itoa(int(send.Data.UID)),
		})
	} else {
		logic.PushToBulletSender(msg)
	}
}

func blindBoxGiftThanks(send *entity.SendGiftText) string {
	box := send.Data.BlindGift.OriginalGiftName
	if !strings.Contains(box, "盲盒") {
		box += "盲盒"
	}
	profit := (send.Data.Price - send.Data.BlindGift.OriginalGiftPrice) * send.Data.Num
	if profit > 0 {
		return fmt.Sprintf("感谢%s赠送的%s，爆出%s，血赚了%d电池", send.Data.Uname, box, send.Data.GiftName, profit)
	}
	if profit < 0 {
		return fmt.Sprintf("感谢%s赠送的%s，爆出%s，血亏了%d电池", send.Data.Uname, box, send.Data.GiftName, -profit)
	}
	return fmt.Sprintf("感谢%s赠送的%s，爆出%s，没赚没亏", send.Data.Uname, box, send.Data.GiftName)
}
