package danmu

import (
	"context"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-module/carbon/v2"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/model"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

var info string = "签到服务异常"

func DosignInProcess(msg, uid, username string, svcCtx *svc.ServiceContext, reply ...*entity.DanmuMsgTextReplyInfo) {
	if msg != "签到" && msg != "打卡" {
		return
	}
	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		logx.Error(err)
		logic.PushToBulletSender(info, reply...)
	}
	// 获取当前时间
	now := carbon.Now(carbon.Local)
	signInfo, err := svcCtx.SignInModel.FindOne(context.Background(), id)
	switch err {
	case nil:
		lastdate := carbon.CreateFromTimestamp(signInfo.LastDay, carbon.Local)
		if lastdate.Year() != now.Year() || lastdate.Month() != now.Month() || lastdate.Day() != now.Day() {
			err := svcCtx.SignInModel.UpdateCount(context.Background(), id)
			if err != nil {
				logic.PushToBulletSender(info, reply...)
				logx.Error(err)
				return
			}
			logic.PushToBulletSender(fmt.Sprintf("已签到%v天", signInfo.Count+1), reply...)
		} else {
			logic.PushToBulletSender(fmt.Sprintf("今天已经签到过了,已签到%v天", signInfo.Count), reply...)
		}
	case model.ErrNotFound:
		data := model.SingInBase{
			Uid:     id,
			LastDay: now.Timestamp(),
			Count:   1,
		}
		err := svcCtx.SignInModel.Insert(context.Background(), nil, &data)
		if err != nil {
			logic.PushToBulletSender(info, reply...)
			logx.Error(err)
			return
		}
		logic.PushToBulletSender("已签到1天", reply...)
	default:
		logic.PushToBulletSender(info, reply...)
		logx.Error(err)
		return
	}

}
