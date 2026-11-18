package pk

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// Register 注册 PK 相关的直播间事件回调
func Register(c *client.Client, svcCtx *svc.ServiceContext) {
	registerBattleStart(c, svcCtx)
	registerBattleEnd(c, svcCtx)
}
