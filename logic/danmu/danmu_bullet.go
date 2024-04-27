package danmu

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var danmuHandler *DanmuLogic

type DanmuLogic struct {
	danmuChan chan string
}

func PushToBDanmuLogic(bullet string) {
	danmuHandler.danmuChan <- bullet
}

func StartDanmuLogic(ctx context.Context, svcCtx *svc.ServiceContext) {
	var err error

	danmuHandler = &DanmuLogic{
		danmuChan: make(chan string, 1000),
	}

	var msg string
	for {
		select {
		case <-ctx.Done():
			goto END
		case msg = <-danmuHandler.danmuChan:
			danmu := &entity.DanmuMsgText{}
			err = json.Unmarshal([]byte(msg), danmu)
			if err != nil {
				logx.Error(err)
			}
			danmumsg := danmu.Info[1].(string)
			from := danmu.Info[2].([]interface{})
			uid := fmt.Sprintf("%.0f", from[0].(float64))
			re := regexp.MustCompile("\\[(.*?)\\]")
			danmumsg = re.ReplaceAllString(danmumsg, "")
			if len(danmumsg) > 0 && uid != svcCtx.RobotID {
				// 机器人相关
				go DoDanmuProcess(danmumsg, svcCtx)
				// 签到
				if svcCtx.Config.SignInEnable {
					go DosignInProcess(danmumsg, uid, from[1].(string), svcCtx)
				}
				// 抽签
				if svcCtx.Config.DrawByLot {
					go DodrawByLotProcess(danmumsg, from[1].(string), svcCtx)

				}
				// 抽奖
				if svcCtx.Config.LotteryEnable {
					go DoLotteryProcess(danmumsg, uid, from[1].(string), svcCtx)
				}
				// 主播指令控制
				go DoCMDProcess(danmumsg, uid, svcCtx)
				// 关键词回复
				if svcCtx.Config.KeywordReply {
					go KeywordReply(danmumsg, svcCtx)
				}
			}
			// 实时输出弹幕消息
			logx.Infof("%v %s:%s", uid, from[1], danmumsg)
		}

	}
END:
}
