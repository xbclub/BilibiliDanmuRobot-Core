package welcome

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// Register 注册欢迎相关的直播间事件回调
func Register(c *client.Client, svcCtx *svc.ServiceContext) {
	registerEntryEffect(c, svcCtx)
	registerInteractWord(c, svcCtx)
	// 天选自动关闭欢迎
	registerAnchorLot(c, svcCtx)
}
