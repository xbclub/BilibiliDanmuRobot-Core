package handler

import "github.com/xbclub/BilibiliDanmuRobot-Core/logic"

// 天选
func (w *wsHandler) anchorLot() {
	// 天选启动
	w.client.RegisterCustomEventHandler("ANCHOR_LOT_START", func(s string) {
		if w.svc.Config.InteractWord || w.svc.Config.EntryEffect || w.svc.Config.WelcomeHighWealthy {
			w.svc.Config.InteractWord = false
			w.svc.Config.EntryEffect = false
			w.svc.Config.WelcomeHighWealthy = false
		}
		logic.PushToBulletSender("识别到天选，欢迎弹幕已临时关闭")	})
	// 天选中奖
	w.client.RegisterCustomEventHandler("ANCHOR_LOT_AWARD", func(s string) {
		if w.svc.Config.InteractWord != w.svc.Autointerract.InteractWord {
			w.svc.Config.InteractWord = w.svc.Autointerract.InteractWord
		}
		if w.svc.Config.EntryEffect != w.svc.Autointerract.EntryEffect {
			w.svc.Config.EntryEffect = w.svc.Autointerract.EntryEffect
		}
		if w.svc.Config.WelcomeHighWealthy != w.svc.Autointerract.WelcomeHighWealthy {
			w.svc.Config.WelcomeHighWealthy = w.svc.Autointerract.WelcomeHighWealthy
		}
		logic.PushToBulletSender("天选结束，欢迎弹幕已恢复默认")
	})
}
