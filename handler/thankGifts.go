package handler

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
)

// 礼物感谢
func (w *wsHandler) thankGifts() {
	w.client.RegisterCustomEventHandler("SEND_GIFT", func(s string) {
		if w.svc.Config.ThanksGift {
			send := &entity.SendGiftText{}
			_ = json.Unmarshal([]byte(s), send)
			logic.PushToGiftChan(send)
		}
	})
	w.client.RegisterCustomEventHandler("GUARD_BUY", func(s string) {
		if w.svc.Config.ThanksGift {
			send := &entity.GuardBuyText{}
			_ = json.Unmarshal([]byte(s), send)
			if w.svc.Config.ThanksGiftUseAt {
				logic.PushToGuardChan(send, &entity.DanmuMsgTextReplyInfo{
					ReplyUid: strconv.Itoa(send.Data.Uid),
				})
			} else {
				logic.PushToGuardChan(send)
			}
		}
	})

	w.client.RegisterCustomEventHandler("COMMON_NOTICE_DANMAKU", func(s string) {
		if w.svc.Config.ThanksGift {
			data := &entity.CommonNoticeDanmaku{}
			_ = json.Unmarshal([]byte(s), data)
			if len(data.Data.ContentSegments) == 5 &&
				data.Data.ContentSegments[1].Text == "投喂" &&
				data.Data.ContentSegments[2].Text == "大航海盲盒" {

				logic.PushToBulletSender(fmt.Sprintf("感谢 %s 的 %s", data.Data.ContentSegments[0].Text, data.Data.ContentSegments[4].Text))

			}
		}
	})
}
