package handler

import (
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/xbclub/BilibiliDanmuRobot-Core/utiles/danmuProcess"
	"github.com/zeromicro/go-zero/core/logx"
)

// 弹幕处理
func (w *wsHandler) receiveDanmu() {
	//弹幕处理的功能类接口
	var danmuProcessFunc danmuProcess.DanmuProcess
	//弹幕处理的功能类的集合
	danmuProcessClass := new(danmuProcess.DanmuProcessClass)
	//判断启用机器人聊天
	if w.svc.Config.FuzzyMatchCmd {
		danmuProcessClass.GptClass = new(danmuProcess.Gpt)
		danmuProcessClass.GptClass.SvcCtx = w.svc
	}
	w.client.OnDanmaku(func(danmaku *message.Danmaku) {
		if danmaku.Sender.Uid != w.userId {
			//启用机器人聊天
			if danmuProcessClass.GptClass != nil {
				danmuProcessClass.GptClass.DanmuContent = &danmaku.Content
				danmuProcessFunc = danmuProcessClass.GptClass
				danmuProcessFunc.DoDanmuProcess()
			}
		}
		// 实时输出弹幕消息
		logx.Infof("%d %s:%s", danmaku.Sender.Uid, danmaku.Sender.Uname, danmaku.Content)
	})
}
