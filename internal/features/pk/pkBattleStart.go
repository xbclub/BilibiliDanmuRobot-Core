package pk

import (
	"encoding/json"

	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func registerBattleStart(c *client.Client, svcCtx *svc.ServiceContext) {
	c.RegisterCustomEventHandler("PK_BATTLE_START_NEW", func(s string) {
		svcCtx.EventRLock()
		defer svcCtx.EventRUnlock()
		pkbattlestartfunc(svcCtx, s)
	})
	c.RegisterCustomEventHandler("PK_BATTLE_START", func(s string) {
		svcCtx.EventRLock()
		defer svcCtx.EventRUnlock()
		pkbattlestartfunc(svcCtx, s)
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
		if roomid == 0 {
			logx.Error("未获取的pk对手信息")
		} else {
			PushToPKChan(&roomid)
		}

	}
}
