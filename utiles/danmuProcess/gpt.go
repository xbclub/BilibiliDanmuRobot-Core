package danmuProcess

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

type Gpt struct {
	DanmuContent *string
	SvcCtx       *svc.ServiceContext
}

func (gpt *Gpt) DoDanmuProcess() {
	// @帮助 打出来关键词
	if strings.Compare("@帮助", *gpt.DanmuContent) == 0 {
		s := fmt.Sprintf("发送带有 %s 的弹幕和我互动", gpt.SvcCtx.Config.TalkRobotCmd)
		logx.Info(s)
		logic.PushToBulletSender(" ")
		logic.PushToBulletSender(s)
		logic.PushToBulletSender("请尽情调戏我吧!")
	}

	result := checkIsAtMe(gpt.DanmuContent, gpt.SvcCtx)
	if result == none {
		return
	}
	content := ""
	if result == contained {
		content = strings.ReplaceAll(*gpt.DanmuContent, gpt.SvcCtx.Config.TalkRobotCmd, "")
	} else if result == hasPrefix {
		content = strings.TrimPrefix(*gpt.DanmuContent, gpt.SvcCtx.Config.TalkRobotCmd)
	}
	//如果发现弹幕在@我，那么调用机器人进行回复
	if len(content) > 0 && *gpt.DanmuContent != gpt.SvcCtx.Config.EntryMsg {
		logic.PushToBulletRobot(content)
	}
}

// 检查弹幕是否在@我，返回bool和@我要说的内容
func checkIsAtMe(msg *string, svcCtx *svc.ServiceContext) int {
	if strings.Contains(*msg, svcCtx.Config.TalkRobotCmd) {
		return contained
	} else if strings.HasPrefix(*msg, svcCtx.Config.TalkRobotCmd) {
		return hasPrefix
	} else {
		return none
	}
}
