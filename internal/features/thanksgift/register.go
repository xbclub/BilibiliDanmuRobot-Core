package thanksgift

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// Register 注册礼物答谢相关的直播间事件回调
func Register(c *client.Client, svcCtx *svc.ServiceContext) {
	registerGifts(c, svcCtx)
	// 红包
	registerRedPocket(c, svcCtx)
}
