package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/avast/retry-go/v4"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

//func GetDanmuInfo(svcCtx *svc.ServiceContext) (*entity.ResponseBulletInfo, error) {
//	var err error
//	var resp *resty.Response
//	var url = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=" + strconv.Itoa(svcCtx.Config.RoomId) + "&type=0"
//	r := &entity.ResponseBulletInfo{}
//	if resp, err = cli.R().
//		SetHeader("user-agent", userAgent).
//		Get(url); err != nil {
//		logx.Error("请求getDanmuInfo失败：", err)
//		return nil, err
//	}
//	if err = json.Unmarshal(resp.Body(), r); err != nil {
//		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
//		return nil, err
//	}
//
//	return r, nil
//}

func Send(msg string, svcCtx *svc.ServiceContext) error {
	var err error
	var url = "https://api.live.bilibili.com/msg/send"
	var respdata *entity.DanmuResp = new(entity.DanmuResp)
	m := make(map[string]string)
	m["bubble"] = "5"
	m["msg"] = msg
	m["color"] = "4546550"
	//m["mode"] = "4"
	m["fontsize"] = "25"
	m["rnd"] = strconv.FormatInt(time.Now().Unix(), 10)
	m["roomid"] = strconv.Itoa(svcCtx.Config.RoomId)
	m["csrf"] = CookieList["bili_jct"]
	m["csrf_token"] = CookieList["bili_jct"]
	err = retry.Do(func() error {
		_, data, err := postWithFormData(http.MethodPost, url, userAgent, CookieStr, &m)
		if err != nil {
			logx.Errorf("请求send失败：%v", err)
			return err
		}
		err = json.Unmarshal(data, respdata)
		if err != nil {
			logx.Errorf("send弹幕响应解析失败:%v", err)
			return nil
		}
		if respdata.Code != 0 {
			logx.Infof("请求send失败:%s", respdata.Msg)
			return errors.New(respdata.Msg)
		}
		return nil
	}, retry.Attempts(3), retry.Delay(1*time.Second))
	if err != nil {
		logx.Error(err)
	}
	return nil
}
func postWithFormData(method, url, ua, cookie string, postData *map[string]string) (int, []byte, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range *postData {
		w.WriteField(k, v)
	}
	w.Close()
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Cookie", cookie)
	req.Header.Set("user-agent", ua)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logx.Error(err)
		return 0, nil, err
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp.StatusCode, data, nil
}
