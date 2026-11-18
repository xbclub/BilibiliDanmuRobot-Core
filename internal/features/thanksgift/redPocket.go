package thanksgift

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/xbclub/BilibiliDanmuRobot-Core/blivedm-go/client"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var _redPocketCnt = 0
var locked *sync.Mutex = new(sync.Mutex)

// ResetRedPocketState 换房时重置红包计数：残留计数会让新房间
// WINNER_LIST 的「计数归零才恢复欢迎」判定永远无法满足
func ResetRedPocketState() {
	locked.Lock()
	_redPocketCnt = 0
	locked.Unlock()
}

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
			// 原地写 Config 字段属于无锁写，统一走变更器写时复制
			svcCtx.MutateConfig(func(c *config.Config) {
				c.LotteryEnable = false
			})
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
