package thanksgift

import (
	"encoding/json"
	"fmt"
	"strconv"

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
		send := &entity.SendGiftText{}
		if err := json.Unmarshal([]byte(s), send); err != nil {
			logx.Errorf("SEND_GIFT 解析失败: %v", err)
			return
		}
		if svcCtx.Config.ThanksGift {
			logx.Infof("收到礼物 uname=%s gift=%s num=%d price=%d total=%d",
				send.Data.Uname, send.Data.GiftName, send.Data.Num, send.Data.Price, send.Data.TotalCoin)
			PushToGiftChan(send)
		}
		danmuproc.SaveBlindBoxStat(send, svcCtx)
	})
	c.RegisterCustomEventHandler("SEND_GIFT_V2", func(s string) {
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
			if svcCtx.Config.ThanksGift {
				logx.Infof("收到礼物(V2) uname=%s gift=%s num=%d price=%d total=%d",
					send.Data.Uname, send.Data.GiftName, send.Data.Num, send.Data.Price, send.Data.TotalCoin)
				PushToGiftChan(send)
			}
			danmuproc.SaveBlindBoxStat(send, svcCtx)
		}
	})
	c.RegisterCustomEventHandler("GUARD_BUY", func(s string) {
		if svcCtx.Config.ThanksGift {
			send := &entity.GuardBuyText{}
			_ = json.Unmarshal([]byte(s), send)
			if svcCtx.Config.ThanksGiftUseAt {
				PushToGuardChan(send, &entity.DanmuMsgTextReplyInfo{
					ReplyUid: strconv.Itoa(send.Data.Uid),
				})
			} else {
				PushToGuardChan(send)
			}
		}
	})

	c.RegisterCustomEventHandler("COMMON_NOTICE_DANMAKU", func(s string) {
		if svcCtx.Config.ThanksGift {
			data := &entity.CommonNoticeDanmaku{}
			_ = json.Unmarshal([]byte(s), data)
			if len(data.Data.ContentSegments) == 5 &&
				data.Data.ContentSegments[1].Text == "投喂" &&
				data.Data.ContentSegments[2].Text == "大航海盲盒" {

				logic.PushToBulletSender(fmt.Sprintf("感谢 %s 的 %s", data.Data.ContentSegments[0].Text, data.Data.ContentSegments[4].Text))
			} else if len(data.Data.ContentSegments) == 6 &&
				data.Data.ContentSegments[2].Text == "投喂" &&
				data.Data.ContentSegments[3].Text == "大航海盲盒" {

				logic.PushToBulletSender(fmt.Sprintf("感谢 %s 的 %s", data.Data.ContentSegments[1].Text, data.Data.ContentSegments[5].Text))
			}
		}
	})
}
