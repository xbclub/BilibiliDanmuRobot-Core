package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/zeromicro/go-zero/core/logx"
)

func GetDanmuToken(roomid int, spiInfo *entity.SPIInfo) (danmuAuthDatas *entity.DanmuAuthData, err error) {
	var url = fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%v&type=0", roomid)
	var resp *resty.Response
	cookies := CookieStr + fmt.Sprintf("buvid3=%s;", spiInfo.Data.B3) + fmt.Sprintf("buvid4=%s;", spiInfo.Data.B4)

	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		SetHeader("Cookie", cookies).
		Get(url); err != nil {
		logx.Error("弹幕流秘钥获取失败：", err)
		return nil, err
	}
	// 先解析响应状态
	logx.Alert(string(resp.Body()))
	danmuAuthDatas = &entity.DanmuAuthData{}
	if err = json.Unmarshal(resp.Body(), danmuAuthDatas); err != nil {
		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}
	if danmuAuthDatas.Code != 0 {
		logx.Error(danmuAuthDatas.Message)
		return nil, errors.New(danmuAuthDatas.Message)
	}
	return danmuAuthDatas, nil
}
