package danmu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/http"

	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

// DoLotteryProcess 执行抽奖
func DoLotteryProcess(msg, uid, username string, svcCtx *svc.ServiceContext) {
	if strings.Compare(msg, "抽奖") != 0 {
		return
	}

	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		logx.Error(err)
		logic.PushToBulletSender("抽奖服务异常")
	}

	_ = id

	// 获取抽奖地址
	url := svcCtx.Config.LotteryUrl
	if url == "" {
		logic.PushToBulletSender("未配置抽奖地址")
		return
	}

	// 请求抽奖
	req := &entity.LotteryRequest{
		Msg:      msg,
		Uid:      id,
		Username: username,
		RoomID:   int64(svcCtx.Config.RoomId),
		Version:  "1.0",
	}

	// 请求抽奖地址
	resp, err := http.GetLucky(url, req)
	if err != nil {
		logx.Error(err)
		logic.PushToBulletSender("抽奖服务异常")
		return
	}

	logx.Infof("抽奖结果: %s", resp.Msg)
	// 返回抽奖结果
	logic.PushToBulletSender(fmt.Sprintf("@%s, %s", username, resp.Msg))
}
