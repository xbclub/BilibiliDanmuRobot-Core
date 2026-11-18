package danmuproc

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"strconv"
)

// DoCMDProcess 主播指令控制欢迎弹幕开关
func DoCMDProcess(msg, uid string, svcCtx *svc.ServiceContext) {
	if uid == strconv.FormatInt(svcCtx.UserID, 10) {
		switch msg {
		case "关闭欢迎弹幕":
			svcCtx.SetWelcomeDefault(false)
			logic.PushToBulletSender("已临时关闭欢迎弹幕")
		case "开启欢迎弹幕":
			svcCtx.SetWelcomeDefault(true)
			logic.PushToBulletSender("已临时开启欢迎弹幕")
		}
	}
}
