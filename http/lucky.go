package http

import (
	"encoding/json"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
)

func GetLucky(url string, req *entity.LotteryRequest) (*entity.LotteryResponse, error) {
	var response entity.LotteryResponse
	if resp, err := cli.R().
		SetHeader("user-agent", userAgent).
		SetBody(req).
		Post(url); err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal(resp.Body(), &response); err != nil {
			return nil, err
		}
		return &response, nil
	}
}
