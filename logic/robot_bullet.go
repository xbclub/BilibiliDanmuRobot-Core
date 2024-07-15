package logic

import (
	"context"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/http"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strings"
)

var robot *BulletRobot

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

func StartBulletRobot(ctx context.Context, svcCtx *svc.ServiceContext) {
	robot = &BulletRobot{
		bulletRobotChan: make(chan entity.Bullet, 1000),
	}

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
	if svcCtx.Config.RobotMode == "ChatGPT" {
		if reply, err = http.RequestChatgptRobot(content.Msg, svcCtx); err != nil {
			logx.Errorf("请求机器人失败：%v", err)
			PushToBulletSender("不好意思，机器人坏掉了...", content.Reply...)
			return
		}
	} else {
		if reply, err = http.RequestQingyunkeRobot(content.Msg); err != nil {
			logx.Errorf("请求机器人失败：%v", err)
			PushToBulletSender("不好意思，机器人坏掉了...", content.Reply...)
			return
		}
		bulltes := splitRobotReply(reply, svcCtx)
		for _, v := range bulltes {
			PushToBulletSender(v, content.Reply...)
		}
		return
	}
	PushToBulletSender(reply, content.Reply...)
	logx.Infof("机器人回复：%s", reply)

}

// 将机器人回复语句中的 {br} 进行分割
// b站弹幕一次只能发20个字符，需要切分
func splitRobotReply(content string, svcCtx *svc.ServiceContext) []string {

	// 将机器人回复中的菲菲替换为橘子
	content = strings.ReplaceAll(content, "菲菲", svcCtx.Config.RobotName)
	re, _ := regexp.Compile(`\{face\:.*\}`)
	content = re.ReplaceAllString(content, "")

	//var res []string
	reply := strings.Split(content, "{br}")

	//for _, r := range reply {
	//	// 长度大于20再分割
	//	zh := []rune(r)
	//	if len(zh) > 20 {
	//		i := 0
	//		for i < len(zh) {
	//			if i+20 > len(zh) {
	//				res = append(res, string(zh[i:]))
	//			} else {
	//				res = append(res, string(zh[i:i+20]))
	//			}
	//			i += 20
	//		}
	//	} else {
	//		res = append(res, string(zh))
	//	}
	//}
	return reply
}
