package robot

import (
	"context"
	"github.com/xbclub/BilibiliDanmuRobot-Core/biliapi"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strings"
)

var robot = &BulletRobot{
	bulletRobotChan: make(chan entity.Bullet, 1000),
}

type BulletRobot struct {
	bulletRobotChan chan entity.Bullet
}

func PushToBulletRobot(content string, reply ...*entity.DanmuMsgTextReplyInfo) {
	logx.Infof("PushToBulletRobot成功：%s", content)
	buttle := entity.Bullet{
		Msg:   content,
		Reply: reply,
	}
	robot.bulletRobotChan <- buttle
}

// Run 启动聊天机器人 worker，把弹幕转发给 ChatGPT/青云客并回复
func Run(ctx context.Context, svcCtx *svc.ServiceContext) {

	var content entity.Bullet

	for {
		select {
		case <-ctx.Done():
			goto END
		case content = <-robot.bulletRobotChan:
			handleRobotBullet(content, svcCtx)
		}
	}
END:
}

func handleRobotBullet(content entity.Bullet, svcCtx *svc.ServiceContext) {
	var err error
	var reply string
	// 快照模式：机器人请求是网络 IO，锁内只取模式配置
	svcCtx.EventRLock()
	mode := svcCtx.Config.RobotMode
	svcCtx.EventRUnlock()
	if mode == "ChatGPT" {
		if reply, err = biliapi.RequestChatgptRobot(content.Msg, svcCtx); err != nil {
			logx.Errorf("请求机器人失败：%v", err)
			logic.PushToBulletSender("不好意思，机器人坏掉了...", content.Reply...)
			return
		}
	} else {
		if reply, err = biliapi.RequestQingyunkeRobot(content.Msg); err != nil {
			logx.Errorf("请求机器人失败：%v", err)
			logic.PushToBulletSender("不好意思，机器人坏掉了...", content.Reply...)
			return
		}
		bulltes := splitRobotReply(reply, svcCtx)
		for _, v := range bulltes {
			logic.PushToBulletSender(v, content.Reply...)
		}
		return
	}
	logic.PushToBulletSender(reply, content.Reply...)
	logx.Infof("机器人回复：%s", reply)

}

// 将机器人回复语句中的 {br} 进行分割
func splitRobotReply(content string, svcCtx *svc.ServiceContext) []string {

	// 将机器人回复中的菲菲替换为橘子
	content = strings.ReplaceAll(content, "菲菲", svcCtx.Config.RobotName)
	re, _ := regexp.Compile(`\{face\:.*\}`)
	content = re.ReplaceAllString(content, "")

	reply := strings.Split(content, "{br}")
	return reply
}
