package welcome

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func registerInteractWord(c *client.Client, svcCtx *svc.ServiceContext) {
	c.RegisterCustomEventHandler("INTERACT_WORD", func(s string) {
		var interact entity.InteractWordText
		if err := json.Unmarshal([]byte(s), &interact); err != nil {
			logx.Errorf("INTERACT_WORD 解析失败: %v", err)
			return
		}
		handleInteractEvent(svcCtx, interact.Data.Uid, interact.Data.Uname, interact.Data.MsgType, interact.Data.FansMedal)
	})
	c.RegisterCustomEventHandler("INTERACT_WORD_V2", func(s string) {
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
		handleInteractEvent(svcCtx, ev.UID, ev.Uname, ev.MsgType, ev.FansMedal)
	})
}

func handleInteractEvent(svcCtx *svc.ServiceContext, uid int64, uname string, msgType int32, medal *entity.FansMedalInfo) {
	// 1 进场 2 关注 3 分享 5(互关)
	if msgType == 1 {
		logx.Infof("进场事件 uid=%v uname=%v", uid, uname)
	}
	switch msgType {
	case 1:
		handleEnter(svcCtx, uid, uname, medal)
	case 2, 5:
		handleFollow(svcCtx, uid, uname)
	case 3:
		handleShare(svcCtx, uid, uname)
	default:
		logx.Infof("未识别的互动类型: uid=%v uname=%v msg_type=%v", uid, uname, msgType)
	}
}

func handleEnter(svcCtx *svc.ServiceContext, uid int64, uname string, medal *entity.FansMedalInfo) {
	if !svcCtx.Config.InteractSelf && strconv.FormatInt(uid, 10) == svcCtx.RobotID {
		return
	}
	if !svcCtx.Config.InteractAnchor && uid == svcCtx.UserID {
		return
	}

	if v, ok := svcCtx.Config.WelcomeString[fmt.Sprint(uid)]; svcCtx.Config.WelcomeSwitch && ok {
		PushToInterractChan(&InterractData{Uid: uid, Msg: v})
		return
	}
	if !svcCtx.Config.InteractWord {
		logx.Infof("进场跳过: InteractWord=false uid=%v uname=%v", uid, uname)
		return
	}
	if inWide(uname, svcCtx.Config.WelcomeBlacklistWide) || in(uname, svcCtx.Config.WelcomeBlacklist) {
		logx.Infof("进场跳过: 黑名单 uid=%v uname=%v", uid, uname)
		return
	}

	name := welcomeInteract(uname)
	if svcCtx.OtherSideUid.Has(uid) {
		pushInteractMsg(uid, handleInterract(uid, name, svcCtx))
		return
	}
	if svcCtx.Config.WelcomeTier {
		if msg, ok := tierWelcome(svcCtx, uid, name, medal); ok {
			pushInteractMsg(uid, msg)
			return
		}
	}
	if svcCtx.Config.InteractWordByTime {
		pushInteractMsg(uid, handleInterractByTime(uid, name, svcCtx))
		return
	}
	pushInteractMsg(uid, handleInterract(uid, name, svcCtx))
}

func pushInteractMsg(uid int64, msg string) {
	parts := strings.Split(msg, "\n")
	if len(parts) > 1 {
		for i, s := range parts {
			PushToInterractChan(&InterractData{Uid: uid + int64(i), Msg: s})
		}
		return
	}
	PushToInterractChan(&InterractData{Uid: uid, Msg: msg})
}

func handleFollow(svcCtx *svc.ServiceContext, uid int64, uname string) {
	if !svcCtx.Config.ThanksFocus || len(uname) == 0 {
		return
	}
	if svcCtx.Config.WelcomeUseAt {
		msg := "感谢关注!"
		if len(svcCtx.Config.FocusDanmu) > 0 {
			msg += svcCtx.Config.FocusDanmu[rand.Intn(len(svcCtx.Config.FocusDanmu))]
		}
		logic.PushToBulletSender(msg, &entity.DanmuMsgTextReplyInfo{ReplyUid: strconv.FormatInt(uid, 10)})
		return
	}
	msg := "感谢 " + shortName(uname, 8, svcCtx.Config.DanmuLen) + " 的关注!"
	logic.PushToBulletSender(msg)
	if len(svcCtx.Config.FocusDanmu) > 0 {
		logic.PushToBulletSender(svcCtx.Config.FocusDanmu[rand.Intn(len(svcCtx.Config.FocusDanmu))])
	}
}

func handleShare(svcCtx *svc.ServiceContext, uid int64, uname string) {
	if !svcCtx.Config.ThanksShare || len(uname) == 0 {
		return
	}
	if svcCtx.Config.WelcomeUseAt {
		msg := "感谢分享!"
		if len(svcCtx.Config.FocusDanmu) > 0 {
			msg += svcCtx.Config.FocusDanmu[rand.Intn(len(svcCtx.Config.FocusDanmu))]
		}
		logic.PushToBulletSender(msg, &entity.DanmuMsgTextReplyInfo{ReplyUid: strconv.FormatInt(uid, 10)})
		return
	}
	msg := "感谢 " + shortName(uname, 8, svcCtx.Config.DanmuLen) + " 的分享!"
	logic.PushToBulletSender(msg)
	if len(svcCtx.Config.FocusDanmu) > 0 {
		logic.PushToBulletSender(svcCtx.Config.FocusDanmu[rand.Intn(len(svcCtx.Config.FocusDanmu))])
	}
}

func tierWelcome(svcCtx *svc.ServiceContext, uid int64, uname string, medal *entity.FansMedalInfo) (string, bool) {
	if len(svcCtx.Config.WelcomeFan) == 0 && len(svcCtx.Config.WelcomeNew) == 0 && len(svcCtx.Config.WelcomeBack) == 0 {
		logx.Infof("三档欢迎: uid=%v uname=%v 三组模板均为空，回退旧逻辑", uid, uname)
		return "", false
	}
	first := false
	if svcCtx.VisitModel != nil {
		var err error
		first, err = svcCtx.VisitModel.Touch(context.Background(), uid, time.Now().Unix())
		if err != nil {
			logx.Errorf("记录进场失败 uid=%v: %v", uid, err)
		}
	}
	isFan := isFan(svcCtx, medal)

	var pool []string
	tier := ""
	switch {
	case isFan && len(svcCtx.Config.WelcomeFan) > 0:
		pool = svcCtx.Config.WelcomeFan
		tier = "fan"
	case first && len(svcCtx.Config.WelcomeNew) > 0:
		pool = svcCtx.Config.WelcomeNew
		tier = "new"
	case !first && len(svcCtx.Config.WelcomeBack) > 0:
		pool = svcCtx.Config.WelcomeBack
		tier = "back"
	default:
		logx.Infof("三档欢迎: uid=%v uname=%v first=%v isFan=%v 无匹配模板，回退旧逻辑", uid, uname, first, isFan)
		return "", false
	}
	logx.Infof("三档欢迎: uid=%v uname=%v first=%v isFan=%v tier=%s", uid, uname, first, isFan, tier)
	return renderWelcome(pool[rand.Intn(len(pool))], uname, medal, svcCtx), true
}

func isFan(svcCtx *svc.ServiceContext, medal *entity.FansMedalInfo) bool {
	if medal == nil || medal.MedalName == "" {
		return false
	}
	if medal.TargetID == 0 || medal.TargetID != svcCtx.UserID {
		return false
	}
	return medal.MedalLevel >= int64(svcCtx.Config.WelcomeFanMedalLevel)
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
	if svcCtx.OtherSideUid.Has(uid) {
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
	if svcCtx.OtherSideUid.Has(uid) {
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
