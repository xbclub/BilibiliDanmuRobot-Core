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
)

// registerGifts 礼物感谢
func registerGifts(c *client.Client, svcCtx *svc.ServiceContext) {
	c.RegisterCustomEventHandler("SEND_GIFT", func(s string) {
		send := &entity.SendGiftText{}
		_ = json.Unmarshal([]byte(s), send)
		if svcCtx.Config.ThanksGift {
			PushToGiftChan(send)
		}
		danmuproc.SaveBlindBoxStat(send, svcCtx)
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
