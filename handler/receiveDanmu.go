package handler

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic/danmu"
	"sync"
)

// 弹幕处理
func (w *wsHandler) receiveDanmu() {
	//弹幕处理的功能类接口
	locked := new(sync.Mutex)
	w.client.RegisterCustomEventHandler("DANMU_MSG", func(s string) {
		locked.Lock()
		danmu.PushToBDanmuLogic(s)
		locked.Unlock()
	})
}
