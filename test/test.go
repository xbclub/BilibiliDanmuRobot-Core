package main

import (
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/handler"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"os"
)

func main() {
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
	ctx := svc.NewServiceContext(c)
	cls := handler.NewWsHandler(ctx)
	cls.StartWsClient()
	//time.Sleep(15 * time.Second)
	//cls.StopWsClient()
	select {}
}
