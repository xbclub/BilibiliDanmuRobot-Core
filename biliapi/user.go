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
		kv := strings.SplitN(pair[0], "=", 2)
		if len(kv) < 2 {
			continue
		}
		if CookieValue(kv[0]) == "" {
			AddCookie(kv[0], kv[1])
		}
	}
	//使用追加模式打开文件（持读锁快照，避免与 Logout 等写方竞争）
	cookieMu.RLock()
	snapshotStr := CookieStr
	snapshotList, _ := json.Marshal(CookieList)
	cookieMu.RUnlock()
	// 落盘失败仅记日志（登录本身已成功），不再对 nil file 调 WriteString
	if file, err = os.OpenFile("token/bili_token.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		logx.Errorf("打开文件错误：%v", err)
	} else {
		if _, werr := file.WriteString(snapshotStr); werr != nil {
			logx.Errorf("写入 bili_token.txt 失败：%v", werr)
		}
		file.Close()
	}
	//使用追加模式打开文件
	if file, err = os.OpenFile("token/bili_token.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		logx.Errorf("打开文件错误：%v", err)
		return data, nil
	}
	if _, werr := file.WriteString(string(snapshotList)); werr != nil {
		logx.Errorf("写入 bili_token.json 失败：%v", werr)
	}
	file.Close()
	return data, nil
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
	// 文件 IO 在锁外完成，锁内只做内存解析与重建
	cookieJSON, err := os.ReadFile("token/bili_token.json")
	if err != nil {
		logx.Errorf("打开文件错误：%v", err)
		return err
	}
	cookieMu.Lock()
	CookieStr = string(cookie)
	if err = json.Unmarshal(cookieJSON, &CookieList); err != nil {
		cookieMu.Unlock()
		return err
	}
	// txt 与 json 理论上同源；合并 txt 独有的字段后以浏览器格式（"k=v; k=v"）重建 CookieStr
	mergeCookieStrIntoList(CookieStr)
	rebuildCookieStr()
	cookieMu.Unlock()
	// 补齐 buvid3/buvid4 设备指纹，降低后续请求被风控(412)的概率
	ensureBuvid()
	return err
}
