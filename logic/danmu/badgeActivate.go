package danmu

import (
	"context"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/model"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func BadgeActiveCheckProcess(msg, uid, username string, svcCtx *svc.ServiceContext, reply ...*entity.DanmuMsgTextReplyInfo) {
	todayDate := svcCtx.DanmuCntModel.GetDateStr(0)
	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		logx.Error(err)
		logic.PushToBulletSender(info, reply...)
	}

	todayDanmuCnt, err := svcCtx.DanmuCntModel.FindOne(context.Background(), id, todayDate)
	switch err {
	case nil:
		err := svcCtx.DanmuCntModel.UpdateCount(context.Background(), id)
		if err != nil {
			logic.PushToBulletSender(info, reply...)
			logx.Error(err)
			return
		}
		todayDanmuCnt.Count = todayDanmuCnt.Count + 1

		if todayDanmuCnt.Count == 10 {
			logic.PushToBulletSender(fmt.Sprintf("好耶！今天发了%v条弹幕了耶！", todayDanmuCnt.Count), reply...)
		}
	case model.ErrNotFound:
		data := model.DanmuCntBase{
			Uid:   id,
			Date:  todayDate,
			Count: 1,
		}
		err := svcCtx.DanmuCntModel.Insert(context.Background(), nil, &data)
		if err != nil {
			logic.PushToBulletSender(info, reply...)
			logx.Error(err)
			return
		}
	default:
		//logic.PushToBulletSender(info, reply...)
		logx.Error(err)
		return
	}

	if msg == "查询弹幕" {
		yesterdayDate := svcCtx.DanmuCntModel.GetDateStr(1)
		beforeyesterdayDate := svcCtx.DanmuCntModel.GetDateStr(2)
		yesterdayDanmuCnt, err1 := svcCtx.DanmuCntModel.FindOne(context.Background(), id, yesterdayDate)
		beforeyesterdayDanmuCnt, err2 := svcCtx.DanmuCntModel.FindOne(context.Background(), id, beforeyesterdayDate)
		todayNum := int64(0)
		yesterdayNum := int64(0)
		beforeYesterdayNum := int64(0)
		if err == nil {
			todayNum = todayDanmuCnt.Count
		}
		if err1 == nil {
			yesterdayNum = yesterdayDanmuCnt.Count
		}
		if err2 == nil {
			beforeYesterdayNum = beforeyesterdayDanmuCnt.Count
		}
		logic.PushToBulletSender(fmt.Sprintf("今/昨/前天各发送了：%v，%v，%v条弹幕", todayNum, yesterdayNum, beforeYesterdayNum), reply...)
	}

}
