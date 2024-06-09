package main

import (
	"encoding/json"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/handler"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

func main() {
	cls := handler.NewWsHandler()

	if cls != nil {
		cls.StartWsClient()
	} else {
		fmt.Println("cls is nil")
	}
	fmt.Println(cls.GetSvc().Config.RoomId)
	x := cls.GetSvc()
	z := *x.Config
	z.SignInEnable = false
	z.RoomId = 4699397
	z.CronDanmu = false
	marshal, err := json.Marshal(z)
	if err != nil {
		return
	}
	WriteConfig(string(marshal))
	cls.ReloadConfig()
	fmt.Println(cls.GetSvc().Config.RoomId)
	fmt.Println(cls.GetUserinfo())
	time.Sleep(20 * time.Second)
	z.SignInEnable = true
	marshal, err = json.Marshal(z)
	if err != nil {
		return
	}
	WriteConfig(string(marshal))
	cls.ReloadConfig()
	//time.Sleep(15 * time.Second)
	//cls.StopWsClient()
	select {}
}
func WriteConfig(data string) *ConfigResponse {
	var c config.Config
	resp := new(ConfigResponse)
	err := json.Unmarshal([]byte(data), &c)
	if err != nil {
		logx.Error(err)
		resp.Code = false
		resp.Msg = err.Error()
		return resp
	}
	yamlBytes, err := yaml.Marshal(&c)
	if err != nil {
		logx.Error("Failed to marshal YAML: %v", err)
		resp.Code = false
		resp.Msg = err.Error()
		return resp
	}
	if _, err = os.Stat("./etc"); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.Mkdir("./etc", 0755)
		if err != nil {
			logx.Error(err)
			resp.Code = false
			resp.Msg = err.Error()
			return resp
		}
	}
	file, err := os.OpenFile("etc/bilidanmaku-api.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logx.Errorf("打开文件错误：", err)
		resp.Code = false
		resp.Msg = "打开文件错误：" + err.Error()
		return resp
	}
	_, err = file.Write(yamlBytes)
	if err != nil {
		logx.Errorf("文件写入错误：", err)
		resp.Code = false
		resp.Msg = "文件写入错误：" + err.Error()
		return resp
	}
	file.Close()

	//err = Mustload(&c)
	//if err != nil {
	//	logx.Error(err)
	//	resp.Code = false
	//	resp.Msg = err.Error()
	//	return resp
	//}
	resp.Code = true
	return resp
}

type ConfigResponse struct {
	Code bool
	Msg  string
	Form config.Config
}
