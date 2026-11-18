package thanksgift

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var _redPocketCnt = 0
var locked *sync.Mutex = new(sync.Mutex)

// registerRedPocket 红包感谢，红包期间临时关闭欢迎弹幕
func registerRedPocket(c *client.Client, svcCtx *svc.ServiceContext) {
	c.RegisterCustomEventHandler("POPULARITY_RED_POCKET_NEW", func(s string) {
		send := &entity.RedPocketNew{}
		_ = json.Unmarshal([]byte(s), send)
		locked.Lock()
		_redPocketCnt++
		locked.Unlock()
		if svcCtx.Config.ThanksGift {
			if svcCtx.Config.ThanksGiftUseAt {
				logic.PushToBulletSender(fmt.Sprintf("感谢 %d 电池的 %s", send.Data.Price, send.Data.GiftName), &entity.DanmuMsgTextReplyInfo{
					ReplyUid: strconv.Itoa(send.Data.Uid),
				})
			} else {
				logic.PushToBulletSender(fmt.Sprintf("感谢 %s %d电池的 %s", send.Data.Uname, send.Data.Price, send.Data.GiftName))
			}
		}
		if svcCtx.Config.InteractWord || svcCtx.Config.EntryEffect || svcCtx.Config.WelcomeHighWealthy {
			svcCtx.PauseWelcome()
			svcCtx.Config.LotteryEnable = false
			logic.PushToBulletSender("识别到红包，欢迎弹幕已临时关闭")
		}
	})

	c.RegisterCustomEventHandler("POPULARITY_RED_POCKET_WINNER_LIST", func(s string) {
		locked.Lock()
		_redPocketCnt--
		if _redPocketCnt < 0 {
			_redPocketCnt = 0
		}
		locked.Unlock()
		data := &entity.RedPocketWinnerList{}
		_ = json.Unmarshal([]byte(s), data)

		logx.Info("中奖名单:")
		for _, winner := range data.Data.WinnerInfo {
			w := winner.([]interface{})
			logx.Info(" >>> ", fmt.Sprintf("%.0f", w[0].(float64)), w[1].(string))
		}

		if _redPocketCnt <= 0 {
			svcCtx.ResumeWelcome()
			logic.PushToBulletSender("红包结束，欢迎弹幕已恢复默认")
			_redPocketCnt = 0
		}
	})
}
