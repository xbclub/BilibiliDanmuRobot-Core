package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/Akegarasu/blivedm-go/client"
	_ "github.com/Akegarasu/blivedm-go/utils"
	_ "github.com/glebarez/go-sqlite"
	"github.com/robfig/cron/v3"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/http"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic/danmu"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/xbclub/BilibiliDanmuRobot-Core/utiles"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"os"
	"strconv"
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
	// 特效欢迎
	ineterractCtx    context.Context
	ineterractCancel context.CancelFunc
	//礼物感谢
	thanksGiftCtx   context.Context
	thankGiftCancel context.CancelFunc
	//pk提醒
	pkCtx    context.Context
	pkCancel context.CancelFunc
	//弹幕处理
	danmuLogicCtx    context.Context
	danmuLogicCancel context.CancelFunc
	//定时弹幕
	corndanmu           *cron.Cron
	mapCronDanmuSendIdx map[int]int
	userId              int
}

func NewWsHandler() WsHandler {
	ctx, err := mustloadConfig()
	if err != nil {
		return nil
	}
	ws := new(wsHandler)
	err = ws.starthttp()
	if err != nil {
		logx.Error(err)
		return nil
	}
	ws.client = client.NewClient(ctx.Config.RoomId)
	ws.client.SetCookie(http.CookieStr)
	ws.svc = ctx
	//初始化定时弹幕
	ws.corndanmu = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)))
	ws.mapCronDanmuSendIdx = make(map[int]int)

	// 设置uid作为基本配置
	strUserId, ok := http.CookieList["DedeUserID"]
	if !ok {
		logx.Infof("uid加载失败，请重新登录")
		return nil
	}
	ws.userId, err = strconv.Atoi(strUserId)
	ctx.RobotID = strUserId
	roominfo, err := http.RoomInit(ctx.Config.RoomId)
	if err != nil {
		logx.Error(err)
		//return nil
	}
	ctx.UserID = roominfo.Data.Uid
	return ws
}
func (w *wsHandler) ReloadConfig() error {
	ctx, err := mustloadConfig()
	w.svc = ctx
	if err != nil {
		return err
	}
	if ctx.Config.RoomId != w.svc.Config.RoomId {
		ws := new(wsHandler)
		err = ws.starthttp()
		if err != nil {
			logx.Error(err)
			return err
		}
		ws.client = client.NewClient(ctx.Config.RoomId)
		ws.client.SetCookie(http.CookieStr)
		ws.svc = ctx
		//初始化定时弹幕
		ws.corndanmu = cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
		)))
		ws.mapCronDanmuSendIdx = make(map[int]int)

		// 设置uid作为基本配置
		strUserId, ok := http.CookieList["DedeUserID"]
		if !ok {
			logx.Infof("uid加载失败，请重新登录")
			return errors.New("uid加载失败，请重新登录")
		}
		ws.userId, err = strconv.Atoi(strUserId)
		ctx.RobotID = strUserId
		roominfo, err := http.RoomInit(ctx.Config.RoomId)
		if err != nil {
			logx.Error(err)
			//return err
		}
		ctx.UserID = roominfo.Data.Uid
	}
	return nil
}

type WsHandler interface {
	StartWsClient() error
	StopWsClient()
	SayGoodbye()
	starthttp() error
	ReloadConfig() error
}

func (w *wsHandler) StartWsClient() error {
	w.startLogic()
	if w.svc.Config.EntryMsg != "off" {
		err := http.Send(w.svc.Config.EntryMsg, w.svc)
		if err != nil {
			logx.Error(err)
		}
	}
	return w.client.Start()
}
func (w *wsHandler) StopWsClient() {
	if w.sendBulletCancel != nil {
		w.sendBulletCancel()
	}
	if w.robotBulletCancel != nil {
		w.robotBulletCancel()
	}
	if w.thankGiftCancel != nil {
		w.thankGiftCancel()
	}
	if w.ineterractCancel != nil {
		w.ineterractCancel() // 关闭弹幕姬goroutine
	}
	if w.pkCancel != nil {
		w.pkCancel()
	}
	if w.danmuLogicCancel != nil {
		w.danmuLogicCancel()
	}
	for _, i := range w.corndanmu.Entries() {
		w.corndanmu.Remove(i.ID)
	}
	w.corndanmu.Stop()
	w.client.Stop()
	//w.svc.Db.Db.Close()
}
func (w *wsHandler) SayGoodbye() {
	if len(w.svc.Config.GoodbyeInfo) > 0 {
		err := http.Send(w.svc.Config.GoodbyeInfo, w.svc)
		if err != nil {
			logx.Error(err)
		}
	}
}
func (w *wsHandler) startLogic() {
	w.sendBulletCtx, w.sendBulletCancel = context.WithCancel(context.Background())
	go logic.StartSendBullet(w.sendBulletCtx, w.svc)
	logx.Info("弹幕推送已开启...")
	// 机器人
	w.robotBulletCtx, w.robotBulletCancel = context.WithCancel(context.Background())
	go logic.StartBulletRobot(w.robotBulletCtx, w.svc)
	// 弹幕逻辑
	w.danmuLogicCtx, w.danmuLogicCancel = context.WithCancel(context.Background())
	go danmu.StartDanmuLogic(w.danmuLogicCtx, w.svc)
	w.receiveDanmu()
	logx.Info("弹幕机器人已开启")
	// 特效欢迎
	w.ineterractCtx, w.ineterractCancel = context.WithCancel(context.Background())
	go logic.Interact(w.ineterractCtx)
	w.welcomeEntryEffect()
	w.welcomeInteractWord()
	logx.Info("欢迎模块已开启")
	// 天选自动关闭欢迎
	w.anchorLot()
	// 礼物感谢
	w.thanksGiftCtx, w.thankGiftCancel = context.WithCancel(context.Background())
	go logic.ThanksGift(w.thanksGiftCtx, w.svc)
	w.thankGifts()
	logx.Info("礼物感谢已开启")
	// pk提醒
	w.pkCtx, w.pkCancel = context.WithCancel(context.Background())
	go logic.PK(w.pkCtx, w.svc)
	w.pkBattleStart()
	w.pkBattleEnd()
	// 禁言用户提醒
	w.blockUser()
	// 下播提醒
	// w.sayGoodbyeByWs()
	logx.Info("pk提醒已开启")

	logx.Info("弹幕处理已开启")
	//定时弹幕
	w.corndanmuStart()

}
func (w *wsHandler) starthttp() error {
	var err error
	http.InitHttpClient()
	// 判断是否存在历史cookie
	if http.FileExists("token/bili_token.txt") && http.FileExists("token/bili_token.json") {
		err = http.SetHistoryCookie()
		if err != nil {
			logx.Error("用户登录失败")
			return err
		}
		logx.Info("用户登录成功")
	} else {
		//if err = w.userlogin(); err != nil {
		//	logx.Errorf("用户登录失败：%v", err)
		//	return
		//}
		//logx.Info("用户登录成功")
		logx.Error("用户登录失败")
		return errors.New("用户登录失败")
	}
	return nil
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
func (w *wsHandler) corndanmuStart() {
	if w.svc.Config.CronDanmu == false {
		return
	}
	for n, danmux := range w.svc.Config.CronDanmuList {
		if danmux.Danmu != nil {
			i := n
			danmus := danmux
			_, err := w.corndanmu.AddFunc(danmus.Cron, func() {
				if len(danmus.Danmu) > 0 {
					if danmus.Random {
						logic.PushToBulletSender(danmus.Danmu[rand.Intn(len(danmus.Danmu))])
					} else {
						_, ok := w.mapCronDanmuSendIdx[i]
						if !ok {
							w.mapCronDanmuSendIdx[i] = 0
						}
						w.mapCronDanmuSendIdx[i] = w.mapCronDanmuSendIdx[i] + 1
						logic.PushToBulletSender(danmus.Danmu[w.mapCronDanmuSendIdx[i]%len(danmus.Danmu)])
					}
				}
			})
			if err != nil {
				logx.Errorf("第%d条定时弹幕配置出现错误: %v", i+1, err)
			}
		}
	}
	w.corndanmu.Start()
}
func mustloadConfig() (*svc.ServiceContext, error) {
	dir := "./token"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(fmt.Sprintf("无法创建token文件夹 请手动创建:%s", err))
		}
	}

	var c config.Config
	conf.MustLoad("etc/bilidanmaku-api.yaml", &c, conf.UseEnv())
	logx.MustSetup(c.Log)
	logx.DisableStat()
	//配置数据库文件夹
	info, err := os.Stat(c.DBPath)
	if os.IsNotExist(err) || !info.IsDir() {
		err = os.MkdirAll(c.DBPath, 0777)
		if err != nil {
			logx.Errorf("文件夹创建失败：%s", c.DBPath)
			return nil, err
		}
	}
	ctx := svc.NewServiceContext(c)
	return ctx, err
}
