package handler

import (
	"encoding/json"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/zeromicro/go-zero/core/logx"
)

// 禁言提醒
func (w *wsHandler) blockUser() {
	// 禁言提醒
	w.client.RegisterCustomEventHandler("ROOM_BLOCK_MSG", func(s string) {
		if w.svc.Config.ShowBlockMsg {
			info := &entity.RoomBlockMsg{}
			err := json.Unmarshal([]byte(s), info)
			if err != nil {
				logx.Error(err)
				logx.Errorf("禁言数据解析失败:%s", s)
				return
			}
			op := ""
			oper := ""
			// FIXME: 抓包量少, 还需要更多测试
			// 2 主播禁言 1 管理员禁言
			if info.Data.Operator == 2 {
				oper = "主播"
				op = "禁言"
			} else if info.Data.Operator == 1 {
				oper = "房管"
				op = "禁言"
			} else { // FIXME: 没有抓到数据 待定
				op = "解开禁言"
			}
			s := fmt.Sprintf("用户 %s 被%s %s!", info.Data.UName, oper, op)
			logic.PushToBulletSender(s)
		}
	})
}
