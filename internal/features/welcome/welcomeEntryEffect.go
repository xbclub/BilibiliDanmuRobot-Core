package welcome

import (
	"encoding/json"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// registerEntryEffect 进场特效欢迎
func registerEntryEffect(c *client.Client, svcCtx *svc.ServiceContext) {
	c.RegisterCustomEventHandler("ENTRY_EFFECT", func(s string) {
		entry := &entity.EntryEffectText{}
		_ = json.Unmarshal([]byte(s), entry)

		if !svcCtx.Config.InteractSelf && strconv.Itoa(int(entry.Data.Uid)) == svcCtx.RobotID {
			return
		}
		if !svcCtx.Config.InteractAnchor && entry.Data.Uid == svcCtx.UserID {
			return
		}

		if v, ok := svcCtx.Config.WelcomeString[fmt.Sprint(entry.Data.Uid)]; svcCtx.Config.WelcomeSwitch && ok && svcCtx.Config.EntryEffect {
			PushToInterractChan(&InterractData{
				Uid: entry.Data.Uid,
				Msg: v,
			})
		} else if svcCtx.Config.EntryEffect {
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
			} else if svcCtx.Config.WelcomeHighWealthy {
				if entry.Data.Uinfo.Wealth.Level >= svcCtx.Config.WelcomeHighWealthyLevel {
					msg = entry.Data.Uinfo.Base.Name
				}
			}

			logx.Info(msg)

			if len(msg) > 0 {
				PushToInterractChan(&InterractData{
					Uid: entry.Data.Uid,
					Msg: getRandomWelcome(msg, svcCtx),
				})
			}
		}
	})
}

// getRandomDanmuKeyByTime 从分时段弹幕列表里获取当前时间对应的时间KEY
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

	r := "{user}"
	replace := r + "\n"
	if svcCtx.Config.WelcomeUseAt {
		replace = "，"
		r = " {user}"
	}
	s = strings.ReplaceAll(s, r+", ", replace)
	s = strings.ReplaceAll(s, r+",", replace)
	s = strings.ReplaceAll(s, r+"，", replace)

	s = strings.ReplaceAll(s, r, content)
	return s
}
