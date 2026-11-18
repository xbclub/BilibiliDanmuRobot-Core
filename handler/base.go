package handler

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/robfig/cron/v3"
	"github.com/xbclub/BilibiliDanmuRobot-Core/biliapi"
	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	_ "github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/utils"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/internal/features/danmuproc"
	"github.com/xbclub/BilibiliDanmuRobot-Core/internal/features/notice"
	"github.com/xbclub/BilibiliDanmuRobot-Core/internal/features/pk"
	"github.com/xbclub/BilibiliDanmuRobot-Core/internal/features/robot"
	"github.com/xbclub/BilibiliDanmuRobot-Core/internal/features/thanksgift"
	"github.com/xbclub/BilibiliDanmuRobot-Core/internal/features/welcome"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/qrcode"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"
)

type wsHandler struct {
	client *client.Client
	svc    *svc.ServiceContext
	// 各功能 worker 的统一生命周期
	logicCtx    context.Context
	logicCancel context.CancelFunc
	//定时弹幕
	corndanmu           *cron.Cron
	mapCronDanmuSendIdx map[int]int
	userId              int
	initStart           bool
}

func NewWsHandler() WsHandler {
	ctx, err := mustloadConfig()
	if err != nil {
		return nil
	}
	ws := new(wsHandler)
	ws.initStart = false
	err = ws.starthttp()
	if err != nil {
		logx.Error(err)
		return nil
	}
	ws.client = client.NewClient(ctx.Config.RoomId)
	ws.client.SetCookie(biliapi.CookieStr)
	ws.svc = ctx
	//初始化定时弹幕
	ws.corndanmu = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)))
	ws.mapCronDanmuSendIdx = make(map[int]int)

	// 设置uid作为基本配置
	strUserId, ok := biliapi.CookieList["DedeUserID"]
	if !ok {
		logx.Infof("uid加载失败，请重新登录")
		return nil
	}
	ws.userId, err = strconv.Atoi(strUserId)
	ctx.RobotID = strUserId
	roominfo, err := biliapi.RoomInit(ctx.Config.RoomId)
	if err != nil {
		logx.Error(err)
	} else {
		ctx.UserID = roominfo.Data.Uid
	}
	if userinfo, err := biliapi.Userinfo(ctx.Config.RoomId); err == nil && userinfo != nil {
		ctx.AnchorName = userinfo.Data.Info.Uname
	}
	return ws
}
func (ws *wsHandler) ReloadConfig() error {
	c, err := loadConfig()
	if err != nil {
		return err
	}
	oldconfig := *ws.svc.Config

	if c.RoomId == oldconfig.RoomId {
		// 仅配置热更：只替换 Config 指针，复用原有 DB 连接与 models，
		// 避免每次 reload 重建 sqlite 连接与 AutoMigrate（PR #33 评审 P2）
		ws.svc.Config = &c
		ws.reloadCronIfNeeded(c, oldconfig)
		return nil
	}

	logx.Infof("房间号更改，更换房间号 ：%v", c.RoomId)
	// 先探测新房信息，失败则中止 reload，保持旧连接与旧上下文完全不变，
	// 避免沿用旧房主播 UID 导致老粉档静默失效（PR #33 评审 P2）
	roominfo, err := biliapi.RoomInit(c.RoomId)
	if err != nil {
		logx.Errorf("新房间信息获取失败，保持原房间 %v: %v", oldconfig.RoomId, err)
		return err
	}
	anchorName := ws.svc.AnchorName
	if userinfo, err := biliapi.Userinfo(c.RoomId); err == nil && userinfo != nil {
		anchorName = userinfo.Data.Info.Uname
	}

	// 先把新连接完全立起来，任何一步失败都直接返回，旧连接与旧上下文零改动
	newCli := client.NewClient(c.RoomId)
	newCli.SetCookie(biliapi.CookieStr)
	if err := newCli.Start(); err != nil {
		logx.Errorf("新房间连接失败，保持原房间 %v: %v", oldconfig.RoomId, err)
		return err
	}

	ctx := svc.NewServiceContext(c)
	ctx.RobotID = ws.svc.RobotID
	ctx.UserID = roominfo.Data.Uid
	ctx.AnchorName = anchorName
	ctx.Autointerract = ws.svc.Autointerract

	// 新连接就绪后再切换：先停旧连接收口旧房事件，再原地更新原有
	// ServiceContext（常驻 goroutine 在启动时捕获了 ws.svc 指针，整体替换
	// 会使它们读写旧实例，PR #33 评审 P0），最后注册新 handler——
	// 保证 handler 注册时读到的已是新房上下文，且换房期间无旧房事件写入新房表
	oldCli := ws.client
	oldCli.Stop()
	ws.svc.Config = ctx.Config
	ws.svc.RobotID = ctx.RobotID
	ws.svc.UserID = ctx.UserID
	ws.svc.AnchorName = ctx.AnchorName
	ws.svc.Autointerract = ctx.Autointerract
	ws.svc.SignInModel = ctx.SignInModel
	ws.svc.DanmuCntModel = ctx.DanmuCntModel
	ws.svc.BlindBoxStatModel = ctx.BlindBoxStatModel
	ws.svc.VisitModel = ctx.VisitModel
	// 换房后清空旧房的 PK 串门 uid 集合；不替换 OtherSideUid 实例，
	// 保证 pk goroutine 与新注册的 handler 始终读写同一集合
	ws.svc.OtherSideUid.Clear()

	ws.client = newCli
	ws.registerHandler()
	ws.reloadCronIfNeeded(c, oldconfig)
	return nil
}

func (ws *wsHandler) reloadCronIfNeeded(newCfg, oldCfg config.Config) {
	if newCfg.CronDanmu != oldCfg.CronDanmu || !areSlicesEqual(newCfg.CronDanmuList, oldCfg.CronDanmuList) {
		logx.Info("识别到定时弹幕配置发生变化，重新加载")
		for _, i := range ws.corndanmu.Entries() {
			ws.corndanmu.Remove(i.ID)
		}
		ws.corndanmuStart()
	}
}

type WsHandler interface {
	InitStartWsClient()
	StopWsClient()
	SayGoodbye()
	StopChanel()
	StartWsClient() error
	starthttp() error
	ReloadConfig() error
	GetSvc() svc.ServiceContext
	GetUserinfo() *entity.UserinfoLite
}

func (w *wsHandler) InitStartWsClient() {
	w.startLogic()
}
func (w *wsHandler) StartWsClient() error {
	if w.svc.Config.EntryMsg != "off" {
		err := biliapi.Send(w.svc.Config.EntryMsg, w.svc)
		if err != nil {
			logx.Error(err)
		}
	}
	w.corndanmu.Start()
	w.client = client.NewClient(w.svc.Config.RoomId)
	w.client.SetCookie(biliapi.CookieStr)
	w.registerHandler()
	return w.client.Start()
}
func (w *wsHandler) GetUserinfo() *entity.UserinfoLite {
	return biliapi.GetUserInfo()
}
func (w *wsHandler) GetSvc() svc.ServiceContext {
	return *w.svc
}
func (w *wsHandler) StopWsClient() {
	w.corndanmu.Stop()
	w.client.Stop()
	//w.svc.Db.Db.Close()
}
func (w *wsHandler) StopChanel() {
	if w.logicCancel != nil {
		w.logicCancel()
	}
	for _, i := range w.corndanmu.Entries() {
		w.corndanmu.Remove(i.ID)
	}
}
func (w *wsHandler) SayGoodbye() {
	if len(w.svc.Config.GoodbyeInfo) > 0 {
		for _, msgs := range logic.SplitMessage(w.svc.Config.GoodbyeInfo, w.svc.Config.DanmuLen) {
			err := biliapi.Send(msgs, w.svc)
			if err != nil {
				logx.Errorf("下播弹幕发送失败：%s msg: %s", err, msgs)
			}
			time.Sleep(1 * time.Second) // 防止弹幕发送过快
		}
	}
}
func (w *wsHandler) startLogic() {
	w.logicCtx, w.logicCancel = context.WithCancel(context.Background())
	// 中央弹幕发送队列
	go logic.StartSendBullet(w.logicCtx, w.svc)
	logx.Info("弹幕推送已开启...")
	// 弹幕解析分发
	go danmuproc.Run(w.logicCtx, w.svc)
	// 机器人
	go robot.Run(w.logicCtx, w.svc)
	logx.Info("弹幕机器人已开启")
	// 欢迎模块
	go welcome.Run(w.logicCtx, w.svc)
	logx.Info("欢迎模块已开启")
	// 礼物感谢
	go thanksgift.Run(w.logicCtx, w.svc)
	logx.Info("礼物感谢已开启")
	// pk提醒
	go pk.Run(w.logicCtx, w.svc)

	// 下播提醒
	// w.sayGoodbyeByWs()

	//定时弹幕
	w.corndanmuStart()
}
func (w *wsHandler) registerHandler() {
	// 欢迎（进场特效、进场消息、天选互锁）
	welcome.Register(w.client, w.svc)
	// 弹幕解析分发
	danmuproc.Register(w.client, w.svc)
	logx.Info("弹幕处理已开启")
	// pk提醒
	pk.Register(w.client, w.svc)
	logx.Info("pk提醒已开启")
	// 禁言用户提醒
	notice.Register(w.client, w.svc)
	// 礼物感谢 + 红包
	thanksgift.Register(w.client, w.svc)
}
func (w *wsHandler) starthttp() error {
	var err error
	biliapi.InitHttpClient()
	// 判断是否存在历史cookie
	if biliapi.FileExists("token/bili_token.txt") && biliapi.FileExists("token/bili_token.json") {
		err = biliapi.SetHistoryCookie()
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
	biliapi.InitHttpClient()
	var loginUrl *entity.LoginUrl
	if loginUrl, err = biliapi.GetLoginUrl(); err != nil {
		logx.Error(err)
		return err
	}

	if err = qrcode.Generate(loginUrl.Data.Url); err != nil {
		logx.Error(err)
		return err
	}

	if _, err = biliapi.GetLoginInfo(loginUrl.Data.OauthKey); err != nil {
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
// loadConfig 读取配置文件并完成日志与目录初始化，不建立数据库连接
func loadConfig() (config.Config, error) {
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
	if err != nil || !info.IsDir() {
		err = os.MkdirAll(c.DBPath, 0777)
		if err != nil {
			logx.Errorf("文件夹创建失败：%s", c.DBPath)
			return c, err
		}
	}
	return c, nil
}

func mustloadConfig() (*svc.ServiceContext, error) {
	c, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return svc.NewServiceContext(c), nil
}

// 比较两个 Person 切片是否相同
func areSlicesEqual(a, b []config.CronDanmuList) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}

	return true
}
