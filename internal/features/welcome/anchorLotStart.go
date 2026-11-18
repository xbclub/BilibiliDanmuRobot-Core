package welcome

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// 天选期间临时关闭欢迎弹幕，结束后恢复
func registerAnchorLot(c *client.Client, svcCtx *svc.ServiceContext) {
	// 天选启动
	c.RegisterCustomEventHandler("ANCHOR_LOT_START", func(s string) {
		svcCtx.PauseWelcome()
		logic.PushToBulletSender("识别到天选，欢迎弹幕已临时关闭")
	})
	// 天选中奖
	c.RegisterCustomEventHandler("ANCHOR_LOT_AWARD", func(s string) {
		svcCtx.ResumeWelcome()
		logic.PushToBulletSender("天选结束，欢迎弹幕已恢复默认")
	})
}
