package main

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/handler"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/conf"
)

func main() {
	var c config.Config
	conf.MustLoad("etc/bilidanmaku-api.yaml", &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)
	cls := handler.NewWsHandler(ctx)
	cls.StartWsClient()
	select {}
}
