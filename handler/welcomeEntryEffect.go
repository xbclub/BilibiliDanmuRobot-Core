package handler

import (
	"encoding/json"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 进场特效欢迎
func (w *wsHandler) welcomeEntryEffect() {
	w.client.RegisterCustomEventHandler("ENTRY_EFFECT", func(s string) {
		entry := &entity.EntryEffectText{}
		_ = json.Unmarshal([]byte(s), entry)

		if !w.svc.Config.InteractSelf && strconv.Itoa(int(entry.Data.Uid)) == w.svc.RobotID {
			return
		}
		if !w.svc.Config.InteractAnchor && entry.Data.Uid == w.svc.UserID {
			return
		}

		if v, ok := w.svc.Config.WelcomeString[fmt.Sprint(entry.Data.Uid)]; w.svc.Config.WelcomeSwitch && ok && w.svc.Config.EntryEffect {
			//logic.PushToBulletSender(v)
			logic.PushToInterractChan(&logic.InterractData{
				Uid: entry.Data.Uid,
				Msg: v,
			})
		} else if w.svc.Config.EntryEffect {
			logx.Info("特效欢迎")

			level := ""
			switch entry.Data.Uinfo.Guard.Level {
			case 1:
				level = "总督"
			case 2:
				level = "提督"
			case 3:
				level = "舰长"
			default:
				level = ""
			}

			msg := ""
			if len(level) > 0 {
				msg = fmt.Sprintf("%s %s", level, entry.Data.Uinfo.Base.Name)
			} else if w.svc.Config.WelcomeHighWealthy {
				if entry.Data.Uinfo.Wealth.Level >= w.svc.Config.WelcomeHighWealthyLevel {
					msg = entry.Data.Uinfo.Base.Name
				}
			}

			logx.Info(msg)

			if len(msg) > 0 {
				logic.PushToInterractChan(&logic.InterractData{
					Uid: entry.Data.Uid,
					Msg: getRandomWelcome(msg, w.svc),
				})
			}
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
	content := msg
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

	// // 定义正则表达式
	// re := regexp.MustCompile(`<%(.*?)%>`)

	// // 提取匹配的部分
	// matches := re.FindAllStringSubmatch(msg, -1)
	// // 遍历匹配结果
	// for _, match := range matches {
	// 	content = match[1]
	// 	break
	// }
	// if strings.Contains(msg, "舰长") {
	// 	content = "舰长 " + content
	// }
	r := "{user}"
	s = strings.ReplaceAll(s, r+", ", r+"\n")
	s = strings.ReplaceAll(s, r+",", r+"\n")
	s = strings.ReplaceAll(s, r+"，", r+"\n")
	s = strings.ReplaceAll(s, r, content)
	return s
}
