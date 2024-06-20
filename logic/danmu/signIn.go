package danmu

import (
	"context"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-module/carbon/v2"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/model"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

func DosignInProcess(msg, uid, username string, svcCtx *svc.ServiceContext) {
	if msg != "签到" || msg != "打卡" {
		return
	}
	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		logx.Error(err)
		logic.PushToBulletSender("签到服务异常")
	}
	// 获取当前时间
	now := carbon.Now(carbon.Local)
	signInfo, err := svcCtx.SininModel.FindOne(context.Background(), id)
	switch err {
	case nil:
		lastdate := carbon.CreateFromTimestamp(signInfo.LastDay, carbon.Local)
		if lastdate.Year() != now.Year() || lastdate.Month() != now.Month() || lastdate.Day() != now.Day() {
			err := svcCtx.SininModel.UpdateCount(context.Background(), id)
			if err != nil {
				logic.PushToBulletSender("签到服务异常")
				logx.Error(err)
				return
			}
			logic.PushToBulletSender(fmt.Sprintf("@%s,已签到%v天", username, signInfo.Count+1))
		} else {
			logic.PushToBulletSender(fmt.Sprintf("@%s,今天已经签到过了,已签到%v天", username, signInfo.Count))
		}
	case model.ErrNotFound:
		data := model.SingInBase{
			Uid:     id,
			LastDay: now.Timestamp(),
			Count:   1,
		}
		err := svcCtx.SininModel.Insert(context.Background(), nil, &data)
		if err != nil {
			logic.PushToBulletSender("签到服务异常")
			logx.Error(err)
			return
		}
		logic.PushToBulletSender(fmt.Sprintf("@%s,已签到1天", username))
	default:
		logic.PushToBulletSender("签到服务异常")
		logx.Error(err)
		return
	}

}
