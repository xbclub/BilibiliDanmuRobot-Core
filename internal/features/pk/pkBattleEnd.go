package pk

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

func registerBattleEnd(c *client.Client, svcCtx *svc.ServiceContext) {
	c.RegisterCustomEventHandler("PK_BATTLE_END", func(s string) {
		svcCtx.EventRLock()
		defer svcCtx.EventRUnlock()
		cleanOtherSide(svcCtx)
	})
	c.RegisterCustomEventHandler("PK_END", func(s string) {
		svcCtx.EventRLock()
		defer svcCtx.EventRUnlock()
		cleanOtherSide(svcCtx)
	})
	c.RegisterCustomEventHandler("PK_BATTLE_CRIT", func(s string) {
		svcCtx.EventRLock()
		defer svcCtx.EventRUnlock()
		cleanOtherSide(svcCtx)
	})
	c.RegisterCustomEventHandler("PK_BATTLE_SETTLE_NEW", func(s string) {
		svcCtx.EventRLock()
		defer svcCtx.EventRUnlock()
		cleanOtherSide(svcCtx)
	})
}
func cleanOtherSide(svcCtx *svc.ServiceContext) {
	svcCtx.OtherSideUid.Clear()
}
