package handler

import (
	"context"
	"github.com/Akegarasu/blivedm-go/client"
	_ "github.com/Akegarasu/blivedm-go/utils"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/http"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/xbclub/BilibiliDanmuRobot-Core/utiles"
	"github.com/zeromicro/go-zero/core/logx"
)

type wsHandler struct {
	client *client.Client
	svc    *svc.ServiceContext
	// 机器人
	robotBulletCtx    context.Context
	robotBulletCancel context.CancelFunc
	// 弹幕发送
	sendBulletCtx    context.Context
	sendBulletCancel context.CancelFunc
}

func NewWsHandler(svc *svc.ServiceContext) WsHandler {
	ws := new(wsHandler)
	ws.starthttp()
	ws.client = client.NewClient(svc.Config.RoomId)
	ws.client.SetCookie(http.CookieStr)
	ws.svc = svc
	return ws
}

type WsHandler interface {
	StartWsClient() error
	StopWsClient()
	starthttp()
}

func (w *wsHandler) StartWsClient() error {
	w.startLogic()
	w.welcomedanmu()
	return w.client.Start()
}
func (w *wsHandler) StopWsClient() {

}
func (w *wsHandler) startLogic() {
	w.sendBulletCtx, w.sendBulletCancel = context.WithCancel(context.Background())
	go logic.StartSendBullet(w.sendBulletCtx, w.svc)
	logx.Info("弹幕推送已开启...")

	w.robotBulletCtx, w.robotBulletCancel = context.WithCancel(context.Background())
	go logic.StartBulletRobot(w.robotBulletCtx, w.svc)
	logx.Info("弹幕机器人已开启")
}
func (w *wsHandler) starthttp() {
	var err error
	http.InitHttpClient()
	// 判断是否存在历史cookie
	if http.FileExists("token/bili_token.txt") && http.FileExists("token/bili_token.json") {
		err = http.SetHistoryCookie()
		if err != nil {
			if err = w.userlogin(); err != nil {
				logx.Errorf("用户登录失败：%v", err)
				return
			}
		}
		logx.Info("用户登录成功")
	} else {
		if err = w.userlogin(); err != nil {
			logx.Errorf("用户登录失败：%v", err)
			return
		}
		logx.Info("用户登录成功")
	}
}
func (w *wsHandler) userlogin() error {
	var err error
	http.InitHttpClient()
	var loginUrl *entity.LoginUrl
	if loginUrl, err = http.GetLoginUrl(); err != nil {
		logx.Error(err)
		return err
	}

	if err = utiles.GenerateQr(loginUrl.Data.Url); err != nil {
		logx.Error(err)
		return err
	}

	if _, err = http.GetLoginInfo(loginUrl.Data.OauthKey); err != nil {
		logx.Error(err)
		return err
	}

	return err
}
