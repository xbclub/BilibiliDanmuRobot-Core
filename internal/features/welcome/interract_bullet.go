package welcome

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

var interractGiver = &InterractGiver{
	interractFilter: map[int64]time.Time{},
	locked:          new(sync.Mutex),
	interractChan:   make(chan *InterractData, 1000),
}

type InterractGiver struct {
	interractFilter map[int64]time.Time
	locked          *sync.Mutex
	interractChan   chan *InterractData
}
type InterractData struct {
	Uid   int64
	Msg   string
	Reply *entity.DanmuMsgTextReplyInfo
}

func PushToInterractChan(g *InterractData) {
	interractGiver.interractChan <- g
}

// Run 启动欢迎弹幕去重 worker，同一用户 10 秒内只欢迎一次
func Run(ctx context.Context, svcCtx *svc.ServiceContext) {
	var g *InterractData
	var w = 10 * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			if len(interractGiver.interractFilter) > 0 {
				interractGiver.locked.Lock()
				for k, v := range interractGiver.interractFilter {
					if v.Add(w).Unix() < time.Now().Unix() {
						delete(interractGiver.interractFilter, k)
						logx.Debugf("用户 %v 已从重复过滤列表移除", k)
					}
				}
				interractGiver.locked.Unlock()
			}

			t.Reset(w)
		case g = <-interractGiver.interractChan:
			interractGiver.locked.Lock()
			if value, ok := interractGiver.interractFilter[g.Uid]; ok && value.Add(w).Unix() >= time.Now().Unix() {
				logx.Debugf("用户 %v 10秒内重复欢迎已被过滤", g.Uid)
			} else {
				parts := strings.Split(g.Msg, "\n")
				for _, s := range parts {
					if svcCtx.Config.WelcomeUseAt {
						g.Reply = &entity.DanmuMsgTextReplyInfo{
							ReplyUid: strconv.FormatInt(g.Uid, 10),
						}
						logic.PushToBulletSender(s, g.Reply)
					} else {
						logic.PushToBulletSender(s)
					}
					logx.Debug(s)
				}
				interractGiver.interractFilter[g.Uid] = time.Now()
			}
			interractGiver.locked.Unlock()
			logx.Debugf("用户%v 已进入重复过滤列表", g.Uid)
		}
	}
END:
}
