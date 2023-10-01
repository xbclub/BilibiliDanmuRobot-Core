package handler

import (
	"fmt"
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/xbclub/BilibiliDanmuRobot-Core/http"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

// 欢迎弹幕处理
func (w *wsHandler) welcomedanmu() {
	w.client.OnDanmaku(func(danmaku *message.Danmaku) {
		// @帮助 打出来关键词
		if strings.Compare("@帮助", danmaku.Content) == 0 {
			s := fmt.Sprintf("发送带有 %s 的弹幕和我互动", w.svc.Config.TalkRobotCmd)
			logx.Info(s)
			logic.PushToBulletSender(" ")
			logic.PushToBulletSender(s)
			logic.PushToBulletSender("请尽情调戏我吧!")
		}
		// 如果发现弹幕在@我，那么调用机器人进行回复
		y, content := checkIsAtMe(danmaku.Content, string(rune(danmaku.Sender.Uid)), w.svc)
		if y && len(content) > 0 && danmaku.Content != w.svc.Config.EntryMsg {
			logic.PushToBulletRobot(content)
		}
		// 实时输出弹幕消息
		logx.Infof("%d %s:%s", danmaku.Sender.Uid, danmaku.Sender.Uname, danmaku.Content)
	})
}

// 检查弹幕是否在@我，返回bool和@我要说的内容
func checkIsAtMe(msg, u string, svcCtx *svc.ServiceContext) (bool, string) {
	// 自己发的包含关键字 不与理会 避免递归

	userId, ok := http.CookieList["DedeUserID"]

	if svcCtx.Config.FuzzyMatchCmd {
		if ok && userId != u && strings.Contains(msg, svcCtx.Config.TalkRobotCmd) {
			return true, strings.ReplaceAll(msg, svcCtx.Config.TalkRobotCmd, "")
		} else {
			return false, ""
		}
	} else {
		if ok && userId != u && strings.HasPrefix(msg, svcCtx.Config.TalkRobotCmd) {
			return true, strings.TrimPrefix(msg, svcCtx.Config.TalkRobotCmd)
		} else {
			return false, ""
		}
	}
}
