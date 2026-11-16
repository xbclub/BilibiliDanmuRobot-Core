package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func (w *wsHandler) welcomeInteractWord() {
	w.client.RegisterCustomEventHandler("INTERACT_WORD", func(s string) {
		var interact entity.InteractWordText
		if err := json.Unmarshal([]byte(s), &interact); err != nil {
			logx.Errorf("INTERACT_WORD 解析失败: %v", err)
			return
		}
		w.handleInteractEvent(interact.Data.Uid, interact.Data.Uname, interact.Data.MsgType, interact.Data.FansMedal)
	})
	w.client.RegisterCustomEventHandler("INTERACT_WORD_V2", func(s string) {
		var interact entity.InteractWordV2Text
		if err := json.Unmarshal([]byte(s), &interact); err != nil {
			logx.Errorf("INTERACT_WORD_V2 解析失败: %v", err)
			return
		}
		if interact.Data.Pb == "" {
			return
		}
		ev, err := entity.ParseInteractWordV2(interact.Data.Pb)
		if err != nil {
			logx.Errorf("INTERACT_WORD_V2 pb 解析失败: %v", err)
			return
		}
		w.handleInteractEvent(ev.UID, ev.Uname, ev.MsgType, ev.FansMedal)
	})
}

func (w *wsHandler) handleInteractEvent(uid int64, uname string, msgType int32, medal *entity.FansMedalInfo) {
	switch msgType {
	case 1:
		w.handleEnter(uid, uname, medal)
	case 2, 5:
		w.handleFollow(uid, uname)
	case 3:
		w.handleShare(uid, uname)
	default:
		logx.Infof("未识别的互动类型: uid=%v uname=%v msg_type=%v", uid, uname, msgType)
	}
}

func (w *wsHandler) handleEnter(uid int64, uname string, medal *entity.FansMedalInfo) {
	if !w.svc.Config.InteractSelf && strconv.FormatInt(uid, 10) == w.svc.RobotID {
		return
	}
	if !w.svc.Config.InteractAnchor && uid == w.svc.UserID {
		return
	}

	if v, ok := w.svc.Config.WelcomeString[fmt.Sprint(uid)]; w.svc.Config.WelcomeSwitch && ok {
		logic.PushToInterractChan(&logic.InterractData{Uid: uid, Msg: v})
		return
	}
	if !w.svc.Config.InteractWord {
		return
	}
	if inWide(uname, w.svc.Config.WelcomeBlacklistWide) || in(uname, w.svc.Config.WelcomeBlacklist) {
		return
	}

	name := welcomeInteract(uname)
	if _, ok := w.svc.OtherSideUid[uid]; ok {
		w.pushInteractMsg(uid, handleInterract(uid, name, w.svc))
		return
	}
	if w.svc.Config.WelcomeTier {
		if msg, ok := w.tierWelcome(uid, name, medal); ok {
			w.pushInteractMsg(uid, msg)
			return
		}
	}
	if w.svc.Config.InteractWordByTime {
		w.pushInteractMsg(uid, handleInterractByTime(uid, name, w.svc))
		return
	}
	w.pushInteractMsg(uid, handleInterract(uid, name, w.svc))
}

func (w *wsHandler) pushInteractMsg(uid int64, msg string) {
	parts := strings.Split(msg, "\n")
	if len(parts) > 1 {
		for i, s := range parts {
			logic.PushToInterractChan(&logic.InterractData{Uid: uid + int64(i), Msg: s})
		}
		return
	}
	logic.PushToInterractChan(&logic.InterractData{Uid: uid, Msg: msg})
}

func (w *wsHandler) handleFollow(uid int64, uname string) {
	if !w.svc.Config.ThanksFocus || len(uname) == 0 {
		return
	}
	if w.svc.Config.WelcomeUseAt {
		msg := "感谢关注!"
		if len(w.svc.Config.FocusDanmu) > 0 {
			msg += w.svc.Config.FocusDanmu[rand.Intn(len(w.svc.Config.FocusDanmu))]
		}
		logic.PushToBulletSender(msg, &entity.DanmuMsgTextReplyInfo{ReplyUid: strconv.FormatInt(uid, 10)})
		return
	}
	msg := "感谢 " + shortName(uname, 8, w.svc.Config.DanmuLen) + " 的关注!"
	logic.PushToBulletSender(msg)
	if len(w.svc.Config.FocusDanmu) > 0 {
		logic.PushToBulletSender(w.svc.Config.FocusDanmu[rand.Intn(len(w.svc.Config.FocusDanmu))])
	}
}

func (w *wsHandler) handleShare(uid int64, uname string) {
	if !w.svc.Config.ThanksShare || len(uname) == 0 {
		return
	}
	if w.svc.Config.WelcomeUseAt {
		msg := "感谢分享!"
		if len(w.svc.Config.FocusDanmu) > 0 {
			msg += w.svc.Config.FocusDanmu[rand.Intn(len(w.svc.Config.FocusDanmu))]
		}
		logic.PushToBulletSender(msg, &entity.DanmuMsgTextReplyInfo{ReplyUid: strconv.FormatInt(uid, 10)})
		return
	}
	msg := "感谢 " + shortName(uname, 8, w.svc.Config.DanmuLen) + " 的分享!"
	logic.PushToBulletSender(msg)
	if len(w.svc.Config.FocusDanmu) > 0 {
		logic.PushToBulletSender(w.svc.Config.FocusDanmu[rand.Intn(len(w.svc.Config.FocusDanmu))])
	}
}

func (w *wsHandler) tierWelcome(uid int64, uname string, medal *entity.FansMedalInfo) (string, bool) {
	if len(w.svc.Config.WelcomeFan) == 0 && len(w.svc.Config.WelcomeNew) == 0 && len(w.svc.Config.WelcomeBack) == 0 {
		return "", false
	}
	first := false
	if w.svc.VisitModel != nil {
		var err error
		first, err = w.svc.VisitModel.Touch(context.Background(), uid, time.Now().Unix())
		if err != nil {
			logx.Errorf("记录进场失败 uid=%v: %v", uid, err)
		}
	}

	var pool []string
	switch {
	case w.isFan(medal) && len(w.svc.Config.WelcomeFan) > 0:
		pool = w.svc.Config.WelcomeFan
	case first && len(w.svc.Config.WelcomeNew) > 0:
		pool = w.svc.Config.WelcomeNew
	case !first && len(w.svc.Config.WelcomeBack) > 0:
		pool = w.svc.Config.WelcomeBack
	default:
		return "", false
	}
	return renderWelcome(pool[rand.Intn(len(pool))], uname, medal, w.svc), true
}

func (w *wsHandler) isFan(medal *entity.FansMedalInfo) bool {
	if medal == nil || medal.MedalName == "" {
		return false
	}
	if medal.TargetID == 0 || medal.TargetID != w.svc.UserID {
		return false
	}
	return medal.MedalLevel >= int64(w.svc.Config.WelcomeFanMedalLevel)
}

func renderWelcome(tpl, uname string, medal *entity.FansMedalInfo, svcCtx *svc.ServiceContext) string {
	name := uname
	if svcCtx.Config.WelcomeUseAt {
		tpl = strings.ReplaceAll(tpl, "{at}", "")
	} else {
		name = shortName(uname, 4, svcCtx.Config.DanmuLen)
		tpl = strings.ReplaceAll(tpl, "{at}", "@"+name+" ")
	}
	medalName, medalLevel := "", ""
	if medal != nil {
		medalName = medal.MedalName
		medalLevel = strconv.FormatInt(medal.MedalLevel, 10)
	}
	return strings.NewReplacer(
		"{user}", name,
		"{medal}", medalName,
		"{medal_level}", medalLevel,
		"{anchor}", svcCtx.AnchorName,
	).Replace(tpl)
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
	r := "{user}"
	rep := r + "\n"

	if svcCtx.Config.InteractWordByTime &&
		svcCtx.Config.WelcomeDanmuByTime != nil &&
		len(svcCtx.Config.WelcomeDanmuByTime) > 0 {

		key := getRandomDanmuKeyByTime()

		for _, danmuCfg := range svcCtx.Config.WelcomeDanmuByTime {
			if danmuCfg.Key == key {
				if danmuCfg.Enabled && len(danmuCfg.Danmu) > 0 {
					szWelcomOrig := danmuCfg.Danmu[rand.Intn(len(danmuCfg.Danmu))]

					if svcCtx.Config.WelcomeUseAt {
						rep = "，"
						r = " {user}"
						szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", rep)
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", rep)
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", rep)
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r, "")
						return szWelcomTmp
					}
					welcome := strings.ReplaceAll(szWelcomOrig, r, shortName(uname, 3, svcCtx.Config.DanmuLen))
					rWelcome := []rune(welcome)
					if len(rWelcome) > svcCtx.Config.DanmuLen {
						szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", rep)
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", rep)
						szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", rep)
						return strings.ReplaceAll(szWelcomTmp, r, uname)
					}
					return welcome
				}
				return handleInterract(uid, uname, svcCtx)
			}
		}
		return handleInterract(uid, uname, svcCtx)
	}
	return handleInterract(uid, uname, svcCtx)
}

func handleInterract(uid int64, uname string, svcCtx *svc.ServiceContext) string {
	s := []rune(uname)
	r := "{user}"
	rep := r + "\n"
	if _, ook := svcCtx.OtherSideUid[uid]; ook {
		if svcCtx.Config.WelcomeUseAt {
			return "欢迎过来串门~"
		}
		szWelcom := "欢迎  过来串门~"
		maxLen := (svcCtx.Config.DanmuLen - len([]rune(szWelcom)))
		if len(s) > maxLen && maxLen > 0 {
			return "欢迎 " + string(s[0:maxLen-1]) + "… 过来串门~"
		}
		return "欢迎 " + uname + " 过来串门~"
	}
	szWelcomOrig := svcCtx.Config.WelcomeDanmu[rand.Intn(len(svcCtx.Config.WelcomeDanmu))]

	if svcCtx.Config.WelcomeUseAt {
		rep = "，"
		r = " {user}"
		szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", rep)
		szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", rep)
		szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", rep)
		szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r, rep)
		return szWelcomTmp
	}
	welcome := strings.ReplaceAll(szWelcomOrig, r, shortName(uname, 3, svcCtx.Config.DanmuLen))
	rWelcome := []rune(welcome)
	if len(rWelcome) > svcCtx.Config.DanmuLen {
		szWelcomTmp := strings.ReplaceAll(szWelcomOrig, r+", ", rep)
		szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+",", rep)
		szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r+"，", rep)
		szWelcomTmp = strings.ReplaceAll(szWelcomTmp, r, rep)
		return strings.ReplaceAll(szWelcomTmp, r, uname)
	}
	return welcome
}

func shortName(uname string, alreadyLen, danmuLen int) string {
	s := []rune(uname)
	maxLen := (danmuLen - alreadyLen)
	if len(s) > maxLen && maxLen > 0 {
		return string(s[0:maxLen-1]) + "…"
	}
	return uname
}

func welcomeInteract(name string) string {
	if strings.Contains(name, "欢迎") {
		return strings.Replace(name, "欢迎", "", 1)
	}
	return name
}
