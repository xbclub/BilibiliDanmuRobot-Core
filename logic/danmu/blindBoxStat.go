package danmu

import (
	"context"
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

func DoBlindBoxStat(msg, uid, username string, svcCtx *svc.ServiceContext, reply ...*entity.DanmuMsgTextReplyInfo) {
	if !svcCtx.Config.BlindBoxStat {
		return
	}

	reg := `(?P<month>^[0-9]+)月盲盒$`
	re := regexp.MustCompile(reg)

	match := re.FindStringSubmatch(msg)

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
	}

	// 获取当前时间
	now := carbon.Now(carbon.Local)

	var ret *model.Result
	// 主播查询本月所有数据
	if svcCtx.UserID == id {
		ret, err = svcCtx.BlindBoxStatModel.GetTotal(context.Background(), int16(now.Year()), int16(month), 0)
	} else {
		// 用户查询自己的数据
		ret, err = svcCtx.BlindBoxStatModel.GetTotalOnePersion(context.Background(), id, int16(now.Year()), int16(month), 0)
	}
	if err == nil {
		r := float64(ret.R) / float64(1000.0)
		if ret.R > 0 {
			logic.PushToBulletSender(
				fmt.Sprintf(
					"%s月共开%d个, 赚了＋%.2f元",
					match[1],
					ret.C,
					r,
				),
				reply...,
			)
		} else if ret.R == 0 {
			logic.PushToBulletSender(
				fmt.Sprintf(
					"%s月共开%d个, 没亏没赚!",
					match[1],
					ret.C,
				),
				reply...,
			)
		} else {
			logic.PushToBulletSender(
				fmt.Sprintf(
					"%s月共开%d个, 亏了－%.2f元",
					match[1],
					ret.C,
					math.Abs(r),
				),
				reply...,
			)
		}
	} else {
		logic.PushToBulletSender(errInfo, reply...)
		logx.Alert("盲盒统计出错了!" + err.Error())
	}
}
