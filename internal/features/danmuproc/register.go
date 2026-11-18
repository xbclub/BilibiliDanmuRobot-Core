package danmuproc

import (
	"sync"

	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// Register 注册弹幕消息回调
func Register(c *client.Client, svcCtx *svc.ServiceContext) {
	locked := new(sync.Mutex)
	c.RegisterCustomEventHandler("DANMU_MSG", func(s string) {
		locked.Lock()
		PushToBDanmuLogic(s)
		locked.Unlock()
	})
}
