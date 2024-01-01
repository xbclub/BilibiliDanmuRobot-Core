package handler

import "github.com/xbclub/BilibiliDanmuRobot-Core/svc"

func (w *wsHandler) pkBattleEnd() {
	w.client.RegisterCustomEventHandler("PK_BATTLE_END", func(s string) {
		cleanOtherSide(w.svc)
	})
	w.client.RegisterCustomEventHandler("PK_END", func(s string) {
		cleanOtherSide(w.svc)
	})
	w.client.RegisterCustomEventHandler("PK_BATTLE_CRIT", func(s string) {
		cleanOtherSide(w.svc)
	})
	w.client.RegisterCustomEventHandler("PK_BATTLE_SETTLE_NEW", func(s string) {
		cleanOtherSide(w.svc)
	})
}
func cleanOtherSide(svcCtx *svc.ServiceContext) {
	for k := range svcCtx.OtherSideUid {
		delete(svcCtx.OtherSideUid, k)
	}
}
