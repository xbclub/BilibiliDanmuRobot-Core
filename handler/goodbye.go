package handler

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
)

// 下播输出
func (w *wsHandler) sayGoodbyeByWs() {
	// 下播输出
	w.client.RegisterCustomEventHandler("PREPARING", func(s string) {
		if len(w.svc.Config.GoodbyeInfo) > 0 {
			logic.PushToBulletSender(w.svc.Config.GoodbyeInfo)
		}
	})
}
