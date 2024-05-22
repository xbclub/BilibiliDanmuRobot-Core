package danmu

import (
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

const (
	none int = iota
	contained
	hasPrefix
)

func DoDanmuProcess(msg string, svcCtx *svc.ServiceContext) {
	// @帮助 打出来关键词
	if strings.Compare("@帮助", msg) == 0 {
		s := fmt.Sprintf("发送带有 %s 的弹幕和我互动", svcCtx.Config.TalkRobotCmd)
		logx.Info(s)
		logic.PushToBulletSender(" ")
		logic.PushToBulletSender(s)
		logic.PushToBulletSender("发送 签到 即可签到")
		logic.PushToBulletSender("发送 抽签 即可抽签")
		logic.PushToBulletSender("主播发送 关闭欢迎弹幕 即可关闭欢迎弹幕")
		logic.PushToBulletSender("主播发送 开启欢迎弹幕 即可开启欢迎弹幕")
		logic.PushToBulletSender("请尽情调戏我吧!")
	}

	result := checkIsAtMe(&msg, svcCtx)
	if result == none {
		return
	}
	content := ""
	if result == contained {
		content = strings.ReplaceAll(msg, svcCtx.Config.TalkRobotCmd, "")
	} else if result == hasPrefix {
		content = strings.TrimPrefix(msg, svcCtx.Config.TalkRobotCmd)
	}
	//如果发现弹幕在@我，那么调用机器人进行回复
	if len(content) > 0 && len(svcCtx.Config.TalkRobotCmd) > 0 && msg != svcCtx.Config.EntryMsg {
		logic.PushToBulletRobot(content)
	}
}

// 检查弹幕是否在@我，返回bool和@我要说的内容
func checkIsAtMe(msg *string, svcCtx *svc.ServiceContext) int {
	if strings.Contains(*msg, svcCtx.Config.TalkRobotCmd) && svcCtx.Config.FuzzyMatchCmd {
		return contained
	} else if strings.HasPrefix(*msg, svcCtx.Config.TalkRobotCmd) {
		return hasPrefix
	} else {
		return none
	}
}
