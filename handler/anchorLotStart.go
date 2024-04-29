package handler

import "github.com/xbclub/BilibiliDanmuRobot-Core/logic"

// 天选
func (w *wsHandler) anchorLot() {
	// 天选启动
	w.client.RegisterCustomEventHandler("ANCHOR_LOT_START", func(s string) {
		// 将 LotteryEnable 设置为 false
		w.svc.Config.LotteryEnable = false
		if w.svc.Config.InteractWord || w.svc.Config.EntryEffect {
			w.svc.Config.InteractWord = false
			w.svc.Config.EntryEffect = false
			logic.PushToBulletSender("识别到天选，欢迎、抽奖已临时关闭")
		}
	})
	// 天选中奖
	w.client.RegisterCustomEventHandler("ANCHOR_LOT_AWARD", func(s string) {
		if w.svc.Config.InteractWord != w.svc.Autointerract.InteractWord || w.svc.Config.EntryEffect != w.svc.Autointerract.EntryEffect {
			w.svc.Config.InteractWord = w.svc.Autointerract.InteractWord
		}

		w.svc.Config.EntryEffect = w.svc.Autointerract.EntryEffect
		w.svc.Config.LotteryEnable = true
		logic.PushToBulletSender("天选结束，欢迎、抽奖功能已恢复默认")
	})
}
