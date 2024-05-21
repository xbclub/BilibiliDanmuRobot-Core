package handler

import (
	"encoding/json"

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
			logic.PushToGuardChan(send)
		}
	})
}
