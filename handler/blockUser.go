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
			if info.Data.Operator == 2 {
				op = "解开禁言"
			} else {
				op = "禁言"
			}
			s := fmt.Sprintf("用户 %s 被管理员 %s!", info.Data.UName, op)
			logic.PushToBulletSender(s)
		}
	})
}
