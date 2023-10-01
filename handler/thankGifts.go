package handler

import (
	"encoding/json"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
)

// 礼物感谢
func (w *wsHandler) thankGifts() {
	w.client.RegisterCustomEventHandler("SEND_GIFT", func(s string) {
		send := &entity.SendGiftText{}
		_ = json.Unmarshal([]byte(s), send)
		logic.PushToGiftChan(send)
	})
}
