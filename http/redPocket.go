package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

// 自己UID 房间号 主播UID 抽奖ID
func PostRedPocket(uid, rid, ruid, lotId string) (int, []byte, error) {
	url := "https://api.live.bilibili.com/xlive/lottery-interface/v1/popularityRedPocket/RedPocketDraw?csrf=" + CookieList["bili_jct"]
	ua := "Mozilla/5.0 BiliDroid/6.79.0 (bbcallen@gmail.com) os/android model/Redmi K30 Pro mobi_app/android build/6790300 channel/360 innerVer/6790310 osVer/11 network/2"

	jsonData := []byte(fmt.Sprintf(`{
		"uid": %s,
		"room_id": %s,
		"ruid": %s,
		"lot_id": %s,
		"spm_id": "live.live-room-detail.red-envelope.extract",
		"jump_from": "29000",
		"session_id": "",
		"statistics": "{\"appId\":1,\"platform\":1,\"version\":\"8.0.0\",\"abtest\":\"\"}",
		"ts": 1718861739
	}`, uid, rid, ruid, lotId))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", CookieStr)
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
