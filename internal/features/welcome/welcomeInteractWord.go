package welcome

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
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
		logx.Debugf("进场事件 uid=%v uname=%v", uid, uname)
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

// enterSeenWindow 进场事件短窗去重窗口，与欢迎发送层 10 秒去重保持一致
const enterSeenWindow = 10 * time.Second

var (
	enterSeenMu sync.Mutex
	enterSeen   = map[int64]time.Time{}
)

// dedupEnter V1/V2 双事件去重：B 站过渡期同一次进场可能同时下发
// INTERACT_WORD 与 INTERACT_WORD_V2，短窗内同一 uid 只处理第一条，
// 避免 visit 表重复计数、真·新人被并发事件误判为回头客（PR #33 评审 P1）
func dedupEnter(uid int64) bool {
	return dedupEnterAt(uid, time.Now())
}

func dedupEnterAt(uid int64, now time.Time) bool {
	enterSeenMu.Lock()
	defer enterSeenMu.Unlock()
	if t, ok := enterSeen[uid]; ok && now.Sub(t) < enterSeenWindow {
		return false
	}
	for k, v := range enterSeen {
		if now.Sub(v) >= enterSeenWindow {
			delete(enterSeen, k)
		}
	}
	enterSeen[uid] = now
	return true
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
		// 天选等场景临时关闭欢迎；去重戳放在此判断之后，
		// 避免暂停期间的进场把恢复后 10 秒内的真实进场也吞掉
		logx.Debugf("进场跳过: InteractWord=false uid=%v uname=%v", uid, uname)
		return
	}
	if !dedupEnter(uid) {
		logx.Debugf("进场事件短窗去重: uid=%v uname=%v", uid, uname)
		return
	}
	if inWide(uname, svcCtx.Config.WelcomeBlacklistWide) || in(uname, svcCtx.Config.WelcomeBlacklist) {
		logx.Debugf("进场跳过: 黑名单 uid=%v uname=%v", uid, uname)
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
	// 多行模板整条入队，由发送层统一按 \n 拆分：每行用真实 uid 构造 @ 回复，
	// 且同 uid 的多行与去重窗口一致——旧的 uid+i 拆行入队会 @ 到不存在的
	// UID 并绕过发送层按 uid 的 10 秒去重
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
	if svcCtx.VisitModel == nil {
		logx.Infof("三档欢迎: uid=%v uname=%v VisitModel 未初始化，回退旧逻辑", uid, uname)
		return "", false
	}
	var err error
	first, err = svcCtx.VisitModel.Touch(context.Background(), uid, time.Now().Unix())
	if err != nil {
		// 记录失败时无法判断首访，回退旧逻辑，避免把真·新人误判成回头客
		logx.Errorf("记录进场失败 uid=%v: %v，回退旧逻辑", uid, err)
		return "", false
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

// clipRunes 超长文本按 rune 截断加省略号，n 小于 1 时原样返回
func clipRunes(s string, n int) string {
	r := []rune(s)
	if n > 0 && len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}

func renderWelcome(tpl, uname string, medal *entity.FansMedalInfo, svcCtx *svc.ServiceContext) string {
	name := shortName(uname, 4, svcCtx.Config.DanmuLen)
	if svcCtx.Config.WelcomeUseAt {
		tpl = strings.ReplaceAll(tpl, "{at}", "")
	} else {
		tpl = strings.ReplaceAll(tpl, "{at}", "@"+name+" ")
	}
	medalName, medalLevel := "", ""
	if medal != nil {
		medalName = medal.MedalName
		medalLevel = strconv.FormatInt(medal.MedalLevel, 10)
	}
	// 主播名/粉丝牌名计入长度预算，避免模板超长被 SplitMessage 硬切在词中间；
	// 极小 DanmuLen 时保底 4 rune，预算机制不至于静默失效
	budget := svcCtx.Config.DanmuLen / 2
	if budget < 4 {
		budget = 4
	}
	return strings.NewReplacer(
		"{user}", name,
		"{medal}", clipRunes(medalName, budget),
		"{medal_level}", medalLevel,
		"{anchor}", clipRunes(svcCtx.AnchorName, budget),
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

// in 精确匹配黑名单；不改动共享配置切片，避免 V1/V2 并发派发下的排序竞争
func in(target string, src []string) bool {
	for _, s := range src {
		if s == target {
			return true
		}
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
