package handler

import (
	"encoding/json"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// 进场特效欢迎
func (w *wsHandler) welcomeEntryEffect() {
	// 限制器方案待商议
	//limiter := rate.NewLimiter(1, w.svc.Config.WelcomeTimeLimiter)
	//w.client.RegisterCustomEventHandler("ENTRY_EFFECT", func(s string) {
	//	if limiter.AllowN(time.Now(), w.svc.Config.WelcomeTimeLimiter) {
	//		entry := &entity.EntryEffectText{}
	//		_ = json.Unmarshal([]byte(s), entry)
	//		if v, ok := w.svc.Config.WelcomeString[fmt.Sprint(entry.Data.Uid)]; w.svc.Config.WelcomeSwitch && ok && w.svc.Config.EntryEffect {
	//			logic.PushToBulletSender(v)
	//		} else if w.svc.Config.EntryEffect {
	//			logx.Info("特效欢迎")
	//			logic.PushToInterractChan(&logic.InterractData{
	//				Uid: entry.Data.Uid,
	//				Msg: getRandomWelcome(entry.Data.CopyWriting, w.svc),
	//			})
	//		}
	//	}
	//})
	w.client.RegisterCustomEventHandler("ENTRY_EFFECT", func(s string) {
		entry := &entity.EntryEffectText{}
		_ = json.Unmarshal([]byte(s), entry)
		if v, ok := w.svc.Config.WelcomeString[fmt.Sprint(entry.Data.Uid)]; w.svc.Config.WelcomeSwitch && ok && w.svc.Config.EntryEffect {
			//logic.PushToBulletSender(v)
			logic.PushToInterractChan(&logic.InterractData{
				Uid: entry.Data.Uid,
				Msg: v,
			})
		} else if w.svc.Config.EntryEffect {
			logx.Info("特效欢迎")
			logic.PushToInterractChan(&logic.InterractData{
				Uid: entry.Data.Uid,
				Msg: getRandomWelcome(entry.Data.CopyWriting, w.svc),
			})
		}
	})
}

// 从分时段弹幕列表里获取当前时间对应的时间KEY
func getRandomDanmuKeyByTime() (key string) {
	now := time.Now().Hour()
	switch now {
	case 0, 1:
		// 午夜
		key = "midnight"

	case 2, 3, 4:
		// 凌晨
		key = "earlymorning"

	case 5, 6, 7, 8:
		// 早上
		key = "morning"

	case 9, 10:
		// 上午
		key = "latemorning"

	case 11, 12, 13:
		// 中午
		key = "noon"

	case 14, 15, 16, 17, 18, 19:
		// 下午
		key = "afternoon"

	case 20, 21, 22, 23:
		// 晚上
		key = "night"
	}
	return key
}
func getRandomWelcome(msg string, svcCtx *svc.ServiceContext) string {
	s := ""
	content := ""
	if svcCtx.Config.InteractWordByTime &&
		svcCtx.Config.WelcomeDanmuByTime != nil &&
		len(svcCtx.Config.WelcomeDanmuByTime) > 0 {

		key := getRandomDanmuKeyByTime()

		for _, danmuCfg := range svcCtx.Config.WelcomeDanmuByTime {
			if danmuCfg.Key == key {
				if danmuCfg.Enabled && len(danmuCfg.Danmu) > 0 {
					s = danmuCfg.Danmu[rand.Intn(len(danmuCfg.Danmu))]
				} else {
					s = svcCtx.Config.WelcomeDanmu[rand.Intn(len(svcCtx.Config.WelcomeDanmu))]
				}
				break
			}
		}
	} else {
		s = svcCtx.Config.WelcomeDanmu[rand.Intn(len(svcCtx.Config.WelcomeDanmu))]
	}
	if len(s) == 0 {
		s = svcCtx.Config.WelcomeDanmu[rand.Intn(len(svcCtx.Config.WelcomeDanmu))]
	}

	// 定义正则表达式
	re := regexp.MustCompile(`<%(.*?)%>`)

	// 提取匹配的部分
	matches := re.FindAllStringSubmatch(msg, -1)
	// 遍历匹配结果
	for _, match := range matches {
		content = match[1]
		break
	}
	if strings.Contains(msg, "舰长") {
		content = "舰长 " + content
	}
	r := "{user}"
	s = strings.ReplaceAll(s, r+", ", r+"\n")
	s = strings.ReplaceAll(s, r+",", r+"\n")
	s = strings.ReplaceAll(s, r+"，", r+"\n")
	s = strings.ReplaceAll(s, r, content)
	return s
}
