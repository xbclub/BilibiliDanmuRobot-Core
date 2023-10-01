package handler

import (
	"encoding/json"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func (w *wsHandler) pkBattleStart() {
	w.client.RegisterCustomEventHandler("PK_BATTLE_START_NEW", func(s string) {
		pkbattlestartfunc(w.svc, s)
	})
	w.client.RegisterCustomEventHandler("PK_BATTLE_START", func(s string) {
		pkbattlestartfunc(w.svc, s)
	})
}
func pkbattlestartfunc(svcCtx *svc.ServiceContext, s string) {
	if svcCtx.Config.PKNotice {
		info := &entity.PKStartInfo{}
		roomid := 0
		err := json.Unmarshal([]byte(s), info)
		if err != nil {
			logx.Error(err)
			logx.Errorf("pk数据解析失败:%s", string(s))
			return
		}
		if info.Data.InitInfo.RoomId == svcCtx.Config.RoomId {
			roomid = info.Data.MatchInfo.RoomId
		} else {
			roomid = info.Data.InitInfo.RoomId
		}
		logx.Debug("开始pk")
		//go handlerPK(svcCtx, body)
		if roomid == 0 {
			logx.Error("未获取的pk对手信息")
		} else {
			logic.PushToPKChan(&roomid)
		}

	}
}
