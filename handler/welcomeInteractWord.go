package handler

import (
	"encoding/json"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixMilli()))

func (w *wsHandler) welcomeInteractWord() {
	w.client.RegisterCustomEventHandler("INTERACT_WORD", func(s string) {
		interact := &entity.InteractWordText{}
		_ = json.Unmarshal([]byte(s), interact)
		// 1 进场 2 关注 3 分享 5(互关)
		if interact.Data.MsgType == 1 {
			if !w.svc.Config.InteractSelf && strconv.Itoa(int(interact.Data.Uid)) == w.svc.RobotID {
				return
			}
			if !w.svc.Config.InteractAnchor && interact.Data.Uid == w.svc.UserID {
				return
			}

			if v, ok := w.svc.Config.WelcomeString[fmt.Sprint(interact.Data.Uid)]; w.svc.Config.WelcomeSwitch && ok {
				logic.PushToInterractChan(&logic.InterractData{
					Uid: interact.Data.Uid,
					Msg: v,
				})
			} else if w.svc.Config.InteractWord {
				// 不在黑名单才欢迎
				if !inWide(interact.Data.Uname, w.svc.Config.WelcomeBlacklistWide) &&
					!in(interact.Data.Uname, w.svc.Config.WelcomeBlacklist) {
					if w.svc.Config.InteractWordByTime {
						msg := handleInterractByTime(interact.Data.Uid, welcomeInteract(interact.Data.Uname), w.svc)
						logx.Debug(msg)
						logic.PushToInterractChan(&logic.InterractData{
							Uid: interact.Data.Uid,
							Msg: msg,
						})
					} else {
						msg := handleInterract(interact.Data.Uid, welcomeInteract(interact.Data.Uname), w.svc)
						ms := strings.Split(msg, "\n")
						if len(ms) > 1 {
							for i, s := range ms {
								logic.PushToInterractChan(&logic.InterractData{
									Uid: interact.Data.Uid + int64(i),
									Msg: s,
								})
							}
						} else {
							logic.PushToInterractChan(&logic.InterractData{
								Uid: interact.Data.Uid,
								Msg: msg,
							})
						}
					}
				}
			}
		} else if interact.Data.MsgType == 2 || interact.Data.MsgType == 5 {
			if w.svc.Config.ThanksFocus {
				if len(interact.Data.Uname) == 0 {
					return
				}
				msg := ""
				if w.svc.Config.WelcomeUseAt {
					msg = "感谢关注!" + w.svc.Config.FocusDanmu[random.Intn(len(w.svc.Config.FocusDanmu))]
					logic.PushToBulletSender(msg, &entity.DanmuMsgTextReplyInfo{
						ReplyUid: strconv.FormatInt(interact.Data.Uid, 10),
					})
				} else {
					msg := "感谢 " + shortName(interact.Data.Uname, 8, w.svc.Config.DanmuLen) + " 的关注!"
					logic.PushToBulletSender(msg)
					if w.svc.Config.FocusDanmu != nil && len(w.svc.Config.FocusDanmu) > 0 {
						logic.PushToBulletSender(w.svc.Config.FocusDanmu[random.Intn(len(w.svc.Config.FocusDanmu))])
					}
				}
			}
		} else if interact.Data.MsgType == 3 {
			if w.svc.Config.ThanksShare {
				if len(interact.Data.Uname) == 0 {
					return
				}
				msg := ""
				if w.svc.Config.WelcomeUseAt {
					msg = "感谢分享!" + w.svc.Config.FocusDanmu[random.Intn(len(w.svc.Config.FocusDanmu))]
					logic.PushToBulletSender(msg, &entity.DanmuMsgTextReplyInfo{
						ReplyUid: strconv.FormatInt(interact.Data.Uid, 10),
					})
				} else {
					msg = "感谢 " + shortName(interact.Data.Uname, 8, w.svc.Config.DanmuLen) + " 的分享!"
					logic.PushToBulletSender(msg)
					if w.svc.Config.FocusDanmu != nil && len(w.svc.Config.FocusDanmu) > 0 {
						logic.PushToBulletSender(w.svc.Config.FocusDanmu[random.Intn(len(w.svc.Config.FocusDanmu))])
					}
				}
			}
		} else {
			logx.Info(">>>>>>>>>>>>> 未识别的类型:", s)
		}
	})
}
func inWide(target string, src []string) bool {
	if src != nil {
		for _, s := range src {
			if strings.Contains(target, s) {
				return true
			}
		}
	}
	return false
}
func in(target string, src []string) bool {
	if src != nil {
		sort.Strings(src)
		index := sort.SearchStrings(src, target)
		return index < len(src) && src[index] == target
	}
	return false
}

func handleInterractByTime(uid int64, uname string, svcCtx *svc.ServiceContext) string {

	if _, ook := svcCtx.OtherSideUid[uid]; ook {
		return handleInterract(uid, uname, svcCtx)
	}
	// 凌晨 - Early morning   2:00--5:00
	// 早晨 - Morning   5:00--9:00
	// 上午 - Late morning / Mid-morning  9:00--11:00
	// 中午 - Noon  11:00--14:00
	// 下午 - Afternoon 14:00 -- 20:00
	// 晚上 - Evening / Night 20:00--00:00
	// 午夜 - Midnight 00:00 -- 2:00
	//s := []rune(uname)
	r := "{user}"
	rep := r + "\n"

	if svcCtx.Config.InteractWordByTime &&
		svcCtx.Config.WelcomeDanmuByTime != nil &&
		len(svcCtx.Config.WelcomeDanmuByTime) > 0 {

		key := getRandomDanmuKeyByTime()

		for _, danmuCfg := range svcCtx.Config.WelcomeDanmuByTime {
			if danmuCfg.Key == key {
				if danmuCfg.Enabled && len(danmuCfg.Danmu) > 0 {
					szWelcomOrig := danmuCfg.Danmu[random.Intn(len(danmuCfg.Danmu))]

					if svcCtx.Config.WelcomeUseAt {
						rep = "，"
						r = " {user}"
						szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", rep)
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", rep)
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", rep)
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r, "")
						return szWelcomTmp
					} else {
						welcome := strings.ReplaceAll(szWelcomOrig, r, shortName(uname, 3, svcCtx.Config.DanmuLen))
						rWelcome := []rune(welcome)
						if len(rWelcome) > svcCtx.Config.DanmuLen {
							szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", rep)
							szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", rep)
							szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", rep)
							return strings.ReplaceAll(szWelcomTmp, r, uname)
						} else {
							return welcome
						}
					}
				} else {
					return handleInterract(uid, uname, svcCtx)
				}
			}
		}
		return handleInterract(uid, uname, svcCtx)
	} else {
		return handleInterract(uid, uname, svcCtx)
	}
}
func handleInterract(uid int64, uname string, svcCtx *svc.ServiceContext) string {
	s := []rune(uname)
	r := "{user}"
	rep := r + "\n"
	if _, ook := svcCtx.OtherSideUid[uid]; ook {
		if svcCtx.Config.WelcomeUseAt {
			return "欢迎过来串门~"
		} else {
			szWelcom := "欢迎  过来串门~"
			maxLen := (svcCtx.Config.DanmuLen - len([]rune(szWelcom)))
			if len(s) > maxLen && maxLen > 0 {
				return "欢迎 " + string(s[0:maxLen-1]) + "… 过来串门~"
			} else {
				return "欢迎 " + uname + " 过来串门~"
			}
		}
	} else {
		szWelcomOrig := svcCtx.Config.WelcomeDanmu[random.Intn(len(svcCtx.Config.WelcomeDanmu))]

		if svcCtx.Config.WelcomeUseAt {
			rep = "，"
			r = " {user}"
			szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", rep)
			szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", rep)
			szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", rep)
			szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r, rep)
			return szWelcomTmp
		} else {
			welcome := strings.ReplaceAll(szWelcomOrig, r, shortName(uname, 3, svcCtx.Config.DanmuLen))
			rWelcome := []rune(welcome)
			if len(rWelcome) > svcCtx.Config.DanmuLen {
				szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", rep)
				szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", rep)
				szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", rep)
				szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r, rep)
				return strings.ReplaceAll(szWelcomTmp, r, uname)
			} else {
				return welcome
			}
		}
	}
}
func shortName(uname string, alreadyLen, danmuLen int) string {
	s := []rune(uname)
	maxLen := (danmuLen - alreadyLen)
	if len(s) > maxLen && maxLen > 0 {
		return string(s[0:maxLen-1]) + "…"
	} else {
		return uname
	}
}
func welcomeInteract(name string) string {
	if strings.Contains(name, "欢迎") {
		name = strings.Replace(name, "欢迎", "", 1)
		return name
	} else {
		return name
	}
}
