package main

import (
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/handler"
)

func main() {
	cls := handler.NewWsHandler()

	if cls != nil {
		cls.StartWsClient()
	} else {
		fmt.Println("cls is nil")
	}
	//time.Sleep(15 * time.Second)
	//cls.StopWsClient()
	select {}
}
