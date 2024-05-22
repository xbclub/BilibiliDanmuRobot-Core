package danmu

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"strconv"
)

func DoCMDProcess(msg, uid string, svcCtx *svc.ServiceContext) {
	if uid == strconv.FormatInt(svcCtx.UserID, 10) {
		switch msg {
		case "关闭欢迎弹幕":
			svcCtx.Config.InteractWord = false
			svcCtx.Config.EntryEffect = false
			svcCtx.Config.WelcomeHighWealthy = false
			svcCtx.Autointerract.InteractWord = false
			svcCtx.Autointerract.EntryEffect = false
			svcCtx.Autointerract.WelcomeHighWealthy = false
			logic.PushToBulletSender("已临时关闭欢迎弹幕")
		case "开启欢迎弹幕":
			svcCtx.Config.InteractWord = true
			svcCtx.Config.EntryEffect = true
			svcCtx.Config.WelcomeHighWealthy = true
			svcCtx.Autointerract.InteractWord = true
			svcCtx.Autointerract.EntryEffect = true
			svcCtx.Autointerract.WelcomeHighWealthy = true
			logic.PushToBulletSender("已临时开启欢迎弹幕")
		}
	}
}
