package danmuproc

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/internal/features/robot"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var danmuHandler = &DanmuLogic{
	danmuChan: make(chan string, 1000),
}

type DanmuLogic struct {
	danmuChan chan string
}

func PushToBDanmuLogic(bullet string) {
	danmuHandler.danmuChan <- bullet
}

// Run 启动弹幕解析分发 worker，把 DANMU_MSG 拆给各功能协程处理
func Run(ctx context.Context, svcCtx *svc.ServiceContext) {
	for {
		select {
		case <-ctx.Done():
			goto END
		case msg := <-danmuHandler.danmuChan:
			processDanmuMessage(msg, svcCtx)
		}

	}
END:
}

// processDanmuMessage 处理单条 DANMU_MSG 原始报文：解析、按配置分发到
// 各功能协程、输出日志。抽出独立函数以便用录制报文做回放测试
func processDanmuMessage(msg string, svcCtx *svc.ServiceContext) {
	danmu := &entity.DanmuMsgText{}
	if err := json.Unmarshal([]byte(msg), danmu); err != nil {
		logx.Error(err)
		return
	}

	danmumsg := danmu.Info[1].(string)
	from := danmu.Info[2].([]interface{})
	uid := fmt.Sprintf("%.0f", from[0].(float64))
	re := regexp.MustCompile("\\[(.*?)\\]")
	danmumsg = re.ReplaceAllString(danmumsg, "")

	// 解析Info[0]
	extras := danmu.Info[0].([]interface{})
	extra := extras[15].(map[string]interface{})
	jsonExtra := strings.ReplaceAll(extra["extra"].(string), `\"`, `"`)
	tagExtra := &entity.DanmuMsgTextInfo0Extra{}
	if e := json.Unmarshal([]byte(jsonExtra), tagExtra); e != nil {
		logx.Alert(e.Error())
	}

	reply := &entity.DanmuMsgTextReplyInfo{
		ReplyUid:   uid,
		ReplyMsgId: tagExtra.IdStr,
	}

	cardLv := "0"
	card := "无信仰"
	if len(danmu.Info) > 3 {
		cardInfo := danmu.Info[3].([]interface{})
		if len(cardInfo) > 1 {
			cardLv = fmt.Sprintf("%.0f", cardInfo[0].(float64))
			card = cardInfo[1].(string)
		}
	}
	// 事件读锁内做功能开关判定与分发（spawn-only，不持锁做 IO）；
	// DoCMDProcess 在独立 goroutine 中取写锁，无自升级死锁
	svcCtx.EventRLock()
	if len(danmumsg) > 0 && uid != svcCtx.RobotID {
		// 机器人相关
		go robot.DoDanmuProcess(danmumsg, svcCtx, reply)
		// 弹幕统计
		if svcCtx.Config.DanmuCntEnable {
			go BadgeActiveCheckProcess(danmumsg, uid, from[1].(string), svcCtx, reply)
		}
		// 关键词回复
		if svcCtx.Config.KeywordReply {
			go KeywordReply(danmumsg, svcCtx, reply)
		}
	}
	// 签到
	if svcCtx.Config.SignInEnable {
		go DosignInProcess(danmumsg, uid, from[1].(string), svcCtx, reply)
	}
	// 抽签
	if svcCtx.Config.DrawByLot {
		go DodrawByLotProcess(danmumsg, from[1].(string), svcCtx, reply)
	}
	// 盲盒统计
	if svcCtx.Config.BlindBoxStat {
		go DoBlindBoxStat(danmumsg, uid, from[1].(string), svcCtx, reply)
	}
	if len(danmumsg) > 0 && uid == strconv.FormatInt(svcCtx.UserID, 10) {
		// 主播指令控制
		go DoCMDProcess(danmumsg, uid, svcCtx)
	}
	svcCtx.EventRUnlock()
	// 实时输出弹幕消息
	if tagExtra.ReplyMid > 0 {
		danmumsg = fmt.Sprintf("@%s %s", tagExtra.ReplyUname, danmumsg)
	}
	logx.Infof("%v 「%s %s」%s:%s", uid, cardLv, card, from[1], danmumsg)
}
