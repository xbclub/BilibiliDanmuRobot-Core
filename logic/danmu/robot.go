package danmu

import (
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"strings"
)

const (
	none int = iota
	contained
	hasPrefix
)

func DoDanmuProcess(msg string, svcCtx *svc.ServiceContext, reply ...*entity.DanmuMsgTextReplyInfo) {
	// @帮助 打出来关键词
	if strings.Compare("@帮助", msg) == 0 {
		s := ""
		if len(svcCtx.Config.TalkRobotCmd) > 0 {
			s = fmt.Sprintf("发送带有 %s 的弹幕和我互动", svcCtx.Config.TalkRobotCmd)
			logic.PushToBulletSender(s)
			logic.PushToBulletSender("请尽情调戏我吧!")
		} else {
			s = "互动聊天已禁用..."
			logic.PushToBulletSender(s)
		}
		//logic.PushToBulletSender(" ")
		// logx.Info(s)
		logic.PushToBulletSender("发送「签到/打卡」即可签到")
		logic.PushToBulletSender("发送「查询弹幕」查询自己近三天的弹幕数")
		logic.PushToBulletSender("发送「X月盲盒」查询在本直播间的盲盒盈亏")
		logic.PushToBulletSender("发送「抽签」即可抽签")
		logic.PushToBulletSender("主播发送「关闭欢迎弹幕」即可关闭欢迎弹幕")
		logic.PushToBulletSender("主播发送「开启欢迎弹幕」即可开启欢迎弹幕")
		logic.PushToBulletSender("本软件为永久免费软件")
	}
	if strings.Compare("@作者", msg) == 0 {
		logic.PushToBulletSender("作者为@超凶一只花酱酱")
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
		logic.PushToBulletRobot(content, reply...)
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
