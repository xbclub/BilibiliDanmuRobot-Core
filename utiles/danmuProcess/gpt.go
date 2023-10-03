package danmuProcess

import (
	"fmt"
	"github.com/Akegarasu/blivedm-go/message"
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
	danmuContent *string
	fromUser     *message.User
	svcCtx       *svc.ServiceContext
}

func (gpt *Gpt) Create() DanmuProcess {
	return new(Gpt)
}

func (gpt *Gpt) DoDanmuProcess() {
	// @帮助 打出来关键词
	if strings.Compare("@帮助", *gpt.danmuContent) == 0 {
		s := fmt.Sprintf("发送带有 %s 的弹幕和我互动", gpt.svcCtx.Config.TalkRobotCmd)
		logx.Info(s)
		logic.PushToBulletSender(" ")
		logic.PushToBulletSender(s)
		logic.PushToBulletSender("请尽情调戏我吧!")
	}

	result := checkIsAtMe(gpt.danmuContent, gpt.svcCtx)
	if result == none {
		return
	}
	content := ""
	if result == contained {
		content = strings.ReplaceAll(*gpt.danmuContent, gpt.svcCtx.Config.TalkRobotCmd, "")
	} else if result == hasPrefix {
		content = strings.TrimPrefix(*gpt.danmuContent, gpt.svcCtx.Config.TalkRobotCmd)
	}
	//如果发现弹幕在@我，那么调用机器人进行回复
	if len(content) > 0 && *gpt.danmuContent != gpt.svcCtx.Config.EntryMsg {
		logic.PushToBulletRobot(content)
	}
}

func (gpt *Gpt) SetConfig(svcCtx *svc.ServiceContext) {
	gpt.svcCtx = svcCtx
}

func (gpt *Gpt) SetDanmu(content *string, user *message.User) {
	gpt.danmuContent = content
	gpt.fromUser = user
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
