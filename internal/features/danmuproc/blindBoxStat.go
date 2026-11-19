package danmuproc

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"

	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-module/carbon/v2"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/model"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var errInfo string = "盲盒统计服务异常"

var (
	monthBlindBoxReg = regexp.MustCompile(`^([0-9]+)月盲盒$`)
	dayBlindBoxReg   = regexp.MustCompile(`^(今日|当日|今天|昨日|昨天)盲盒$`)
)

// SaveBlindBoxStat 记录盲盒礼物开销，礼物事件到达时调用
func SaveBlindBoxStat(g *entity.SendGiftText, svcCtx *svc.ServiceContext) {
	logx.Info(g.Data.BlindGift.OriginalGiftName)
	if g.Data.BlindGift.OriginalGiftName == "" {
		return
	}
	now := carbon.Now(carbon.Local)
	err := svcCtx.BlindBoxStatModel.Insert(context.Background(), nil, &model.BlindBoxStatBase{
		Uid:               int64(g.Data.UID),
		BlindBoxName:      g.Data.BlindGift.OriginalGiftName,
		Price:             int32(g.Data.Price),
		OriginalGiftPrice: int32(g.Data.BlindGift.OriginalGiftPrice),
		Cnt:               int32(g.Data.Num),
		Year:              int16(now.Year()),
		Month:             int16(now.Month()),
		Day:               int16(now.Day()),
	})
	if err != nil {
		logx.Alert("保存盲盒数据出错!!! " + err.Error())
	} else {
		logx.Info("盲盒数据保存成功!!! ")
	}
}

// DoBlindBoxStat 处理「今日/昨日盲盒」「X月盲盒」查询
func DoBlindBoxStat(msg, uid, username string, svcCtx *svc.ServiceContext, reply ...*entity.DanmuMsgTextReplyInfo) {
	if !svcCtx.Config.BlindBoxStat {
		return
	}

	if match := dayBlindBoxReg.FindStringSubmatch(msg); len(match) == 2 {
		date := carbon.Now(carbon.Local)
		label := "今日"
		if match[1] == "昨日" || match[1] == "昨天" {
			date = date.SubDays(1)
			label = "昨日"
		}
		id, err := strconv.ParseInt(uid, 10, 64)
		if err != nil {
			logx.Error(err)
			logic.PushToBulletSender(errInfo, reply...)
			return
		}
		var ret *model.Result
		if svcCtx.UserID == id {
			ret, err = svcCtx.BlindBoxStatModel.GetTotal(context.Background(), int16(date.Year()), int16(date.Month()), int16(date.Day()))
		} else {
			ret, err = svcCtx.BlindBoxStatModel.GetTotalOnePersion(context.Background(), id, int16(date.Year()), int16(date.Month()), int16(date.Day()))
		}
		pushBlindBoxResult(label, ret, err, blindBoxDayPrefix(uid, username, svcCtx), reply...)
		return
	}

	match := monthBlindBoxReg.FindStringSubmatch(msg)
	if len(match) != 2 {
		return
	}
	month, err := strconv.Atoi(match[1])
	if err != nil || month < 1 || month > 12 {
		logic.PushToBulletSender(fmt.Sprintf("月份「%s」不正确!", match[1]), reply...)
		return
	}
	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		logx.Error(err)
		logic.PushToBulletSender(errInfo, reply...)
		return
	}
	now := carbon.Now(carbon.Local)
	var ret *model.Result
	if svcCtx.UserID == id {
		ret, err = svcCtx.BlindBoxStatModel.GetTotal(context.Background(), int16(now.Year()), int16(month), 0)
	} else {
		ret, err = svcCtx.BlindBoxStatModel.GetTotalOnePersion(context.Background(), id, int16(now.Year()), int16(month), 0)
	}
	pushBlindBoxResult(fmt.Sprintf("%d年%d月", now.Year(), month), ret, err, "", reply...)
}

func blindBoxDayPrefix(uid, username string, svcCtx *svc.ServiceContext) string {
	if alias := svcCtx.Config.UserNicknames[uid]; alias != "" {
		return alias + " "
	}
	return "@" + username + " "
}

func pushBlindBoxResult(label string, ret *model.Result, err error, prefix string, reply ...*entity.DanmuMsgTextReplyInfo) {
	if errors.Is(err, model.ErrNotFound) {
		ret, err = &model.Result{}, nil
	}
	if err != nil {
		logic.PushToBulletSender(errInfo, reply...)
		logx.Alert("盲盒统计出错了!" + err.Error())
		return
	}
	if ret == nil {
		ret = &model.Result{}
	}
	cost := float64(ret.Cost) / 1000.0
	if ret.R >= 0 {
		logic.PushToBulletSender(fmt.Sprintf("%s%s净赚%.2f元，开%d盒，花费%.2f元", prefix, label, float64(ret.R)/1000.0, ret.C, cost), reply...)
	} else {
		logic.PushToBulletSender(fmt.Sprintf("%s%s净亏%.2f元，开%d盒，花费%.2f元", prefix, label, math.Abs(float64(ret.R)/1000.0), ret.C, cost), reply...)
	}
}
