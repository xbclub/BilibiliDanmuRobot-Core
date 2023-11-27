package danmuProcess

import (
	"context"
	"fmt"
	"github.com/Akegarasu/blivedm-go/message"
	_ "github.com/glebarez/go-sqlite"
	"github.com/golang-module/carbon/v2"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/model"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SignIn struct {
	danmuContent string
	fromUser     message.User
	svcCtx       *svc.ServiceContext
}

func (signIn *SignIn) DoDanmuProcess() {
	if signIn.danmuContent != "签到" {
		return
	}
	// 获取当前时间
	now := carbon.Now(carbon.Local)
	signInfo, err := signIn.svcCtx.SininModel.FindOne(context.Background(), int64(signIn.fromUser.Uid))
	switch err {
	case nil:
		lastdate := carbon.CreateFromTimestamp(signInfo.LastDay, carbon.Local)
		if lastdate.Year() != now.Year() || lastdate.Month() != now.Month() || lastdate.Day() != now.Day() {
			err := signIn.svcCtx.SininModel.UpdateCount(context.Background(), int64(signIn.fromUser.Uid))
			if err != nil {
				logic.PushToBulletSender("签到服务异常")
				logx.Error(err)
				return
			}
			logic.PushToBulletSender(fmt.Sprintf("%s,已签到%v天", signIn.fromUser.Uname, signInfo.Count+1))
		} else {
			logic.PushToBulletSender(fmt.Sprintf("%s,今天已经签到过了,已签到%v天", signIn.fromUser.Uname, signInfo.Count))
		}
	case model.ErrNotFound:
		data := model.SingInBase{
			Uid:     int64(signIn.fromUser.Uid),
			LastDay: now.Timestamp(),
			Count:   1,
		}
		err := signIn.svcCtx.SininModel.Insert(context.Background(), nil, &data)
		if err != nil {
			logic.PushToBulletSender("签到服务异常")
			logx.Error(err)
			return
		}
		logic.PushToBulletSender(fmt.Sprintf("%s,已签到1天", signIn.fromUser.Uname))
	default:
		logic.PushToBulletSender("签到服务异常")
		logx.Error(err)
		return
	}

}

func (signIn *SignIn) Create() DanmuProcess {
	return new(SignIn)
}

func (signIn *SignIn) SetConfig(svcCtx *svc.ServiceContext) {
	signIn.svcCtx = svcCtx
}

func (signIn *SignIn) SetDanmu(content string, user message.User) {
	signIn.danmuContent = content
	signIn.fromUser = user
}
