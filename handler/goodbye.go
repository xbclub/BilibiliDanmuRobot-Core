package handler

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
)

// 下播输出
// 注意：sayGoodbyeByWs 目前未在 registerHandler 中注册，
// 下播告别走 GUI 手动调用 SayGoodbye() 的路径；接线属行为变化，需确认后再启用
func (w *wsHandler) sayGoodbyeByWs() {
	// 下播输出
	w.client.RegisterCustomEventHandler("PREPARING", func(s string) {
		if len(w.svc.Config.GoodbyeInfo) > 0 {
			logic.PushToBulletSender(w.svc.Config.GoodbyeInfo)
		}
	})
}
