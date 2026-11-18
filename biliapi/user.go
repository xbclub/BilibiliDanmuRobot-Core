package biliapi

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/webprofile"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
	"time"
)

func GetLoginUrl() (*entity.LoginUrl, error) {
	var err error
	var resp *resty.Response
	var url = "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"

	r := &entity.LoginUrl{}
	if resp, err = cli.R().
		SetHeaders(webprofile.MainHeaders()).
		Get(url); err != nil {
		logx.Error("请求getLoginUrl失败：", err)
		return nil, err
	}
	if err = json.Unmarshal(resp.Body(), r); err != nil {
		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}

	logx.Info("oauthKey:", r.Data.OauthKey)

	return r, err
}

// 验证登录的同时，将cookie赋值
func GetLoginInfo(oauthKey string) (*entity.LoginInfoData, error) {
	var err error
	var url = "https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key=" + oauthKey
	var resp *resty.Response
	var data *entity.LoginInfoData
	var file *os.File
	pre := &entity.LoginInfoPre{}
	logx.Info("等待扫码登录...")
	for {

		if resp, err = cli.R().
			SetHeaders(webprofile.MainHeaders()).
			Get(url); err != nil {
			logx.Error("请求getLoginInfo失败：", err)
			return nil, err
		}
		if err = json.Unmarshal(resp.Body(), pre); err != nil {
			logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
			return nil, err
		}

		if pre.Code == 0 {
			data = &entity.LoginInfoData{}
			if err = json.Unmarshal(resp.Body(), data); err != nil {
				logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
				return nil, err
			}
			if data.Data.Code == 0 {
				logx.Info("登录成功！")
				break
			} else if data.Data.Code == 86038 {
				logx.Error(data.Data.Message)
				return nil, errors.New(data.Data.Message)
			}
		} else {
			return nil, errors.New(pre.Message)
		}

		time.Sleep(5 * time.Second)
	}
	for _, v := range resp.Header().Values("Set-Cookie") {
		pair := strings.Split(v, ";")
		kv := strings.Split(pair[0], "=")
		if len(kv) < 2 {
			continue
		}
		if _, ok := CookieList[kv[0]]; !ok {
			AddCookie(kv[0], kv[1])
		}
	}
	//使用追加模式打开文件
	file, err = os.OpenFile("token/bili_token.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logx.Errorf("打开文件错误：%v", err)
	}
	file.WriteString(CookieStr)
	file.Close()
	//使用追加模式打开文件
	file, err = os.OpenFile("token/bili_token.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logx.Errorf("打开文件错误：%v", err)
	}
	tokenstr, _ := json.Marshal(CookieList)
	file.WriteString(string(tokenstr))
	file.Close()
	return data, err
}
func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func SetHistoryCookie() error {
	var cookie []byte
	var err error
	cookie, err = os.ReadFile("token/bili_token.txt")
	if err != nil {
		logx.Errorf("打开文件错误：%v", err)
		return err
	}
	CookieStr = string(cookie)
	cookie, err = os.ReadFile("token/bili_token.json")
	if err != nil {
		logx.Errorf("打开文件错误：%v", err)
		return err
	}
	err = json.Unmarshal(cookie, &CookieList)
	if err != nil {

		return err
	}
	// txt 与 json 理论上同源；合并 txt 独有的字段后以浏览器格式（"k=v; k=v"）重建 CookieStr
	mergeCookieStrIntoList(CookieStr)
	rebuildCookieStr()
	// 补齐 buvid3/buvid4 设备指纹，降低后续请求被风控(412)的概率
	ensureBuvid()
	return err
}
