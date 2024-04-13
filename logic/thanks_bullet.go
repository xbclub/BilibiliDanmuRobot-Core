package logic

import (
	"context"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"sync"
	"time"
)

// 检测到礼物，push [uname]->[giftName]->[cost]，number+1
// 每3s统计一次礼物，并进行感谢，礼物价值高于x元加一句大气

var thanksGiver *GiftThanksGiver

type GiftThanksGiver struct {
	giftNotBlindBoxTable map[string]map[string]map[string]int
	giftBlindBoxTable    map[string]map[string]map[string]int
	locked               *sync.Mutex
	tableMu              sync.RWMutex
	giftChan             chan *entity.SendGiftText
}

func PushToGiftChan(g *entity.SendGiftText) {
	thanksGiver.giftChan <- g
}

func ThanksGift(ctx context.Context, svcCtx *svc.ServiceContext) {

	thanksGiver = &GiftThanksGiver{
		giftNotBlindBoxTable: make(map[string]map[string]map[string]int),
		giftBlindBoxTable:    make(map[string]map[string]map[string]int),
		locked:               new(sync.Mutex),
		tableMu:              sync.RWMutex{},
		giftChan:             make(chan *entity.SendGiftText, 1000),
	}

	var g *entity.SendGiftText
	var w = time.Duration(svcCtx.Config.ThanksGiftTimeout) * time.Second
	var t = time.NewTimer(w)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			goto END
		case <-t.C:
			thanksGiver.locked.Lock()
			summarizeGift(svcCtx.Config.DanmuLen, svcCtx.Config.ThanksMinCost)
			thanksGiver.locked.Unlock()
			t.Reset(w)
		case g = <-thanksGiver.giftChan:
			thanksGiver.locked.Lock()
			if _, ok := thanksGiver.giftNotBlindBoxTable[g.Data.Uname]; !ok {
				thanksGiver.giftNotBlindBoxTable[g.Data.Uname] = make(map[string]map[string]int)
			}
			if _, ok := thanksGiver.giftNotBlindBoxTable[g.Data.Uname][g.Data.GiftName]; !ok {
				thanksGiver.giftNotBlindBoxTable[g.Data.Uname][g.Data.GiftName] = make(map[string]int)
			}
			thanksGiver.giftNotBlindBoxTable[g.Data.Uname][g.Data.GiftName]["cost"] += g.Data.Price
			thanksGiver.giftNotBlindBoxTable[g.Data.Uname][g.Data.GiftName]["count"] += g.Data.Num
			//if g.Data.BlindGift.OriginalGiftName == "" {
			//	//fmt.Printf("非盲盒: ")
			//} else {
			//	//fmt.Printf("盲盒: ")
			//	if _, ok := thanksGiver.giftBlindBoxTable[g.Data.Uname]; !ok {
			//		thanksGiver.giftBlindBoxTable[g.Data.Uname] = make(map[string]map[string]int)
			//	}
			//	if _, ok := thanksGiver.giftBlindBoxTable[g.Data.Uname][g.Data.BlindGift.OriginalGiftName]; !ok {
			//		thanksGiver.giftBlindBoxTable[g.Data.Uname][g.Data.BlindGift.OriginalGiftName] = make(map[string]int)
			//	}
			//	thanksGiver.giftBlindBoxTable[g.Data.Uname][g.Data.BlindGift.OriginalGiftName]["count"] += g.Data.Num
			//	thanksGiver.giftBlindBoxTable[g.Data.Uname][g.Data.BlindGift.OriginalGiftName]["profit_and_loss"] += g.Data.Price - g.Data.BlindGift.OriginalGiftPrice
			//}
			thanksGiver.locked.Unlock()
		}
	}
END:
}

func summarizeGift(danmuLen int, minCost int) {
	for name, m := range thanksGiver.giftNotBlindBoxTable {
		sumCost := 0
		giftstring := []string{}
		msg := ""
		for gift, cost := range m {
			giftstring = append(giftstring, fmt.Sprintf("%d个%s", cost["count"], gift))
			// 计算打赏金额
			sumCost += cost["cost"]

			// 感谢完后立刻清空map
			delete(m, gift)
		}

		msgShort := ""

		msg = "感谢" + name + "的"
		for k, v := range giftstring {
			if k == 0 {
				msg += v
				msgShort = v
			} else {
				msg += "，" + v
				msgShort += "，" + v
			}
		}

		ms := []rune(msg)
		if sumCost < minCost {
			// discard
		} else if len(ms) > danmuLen {
			PushToBulletSender("感谢 " + name + " 的")
			PushToBulletSender(msgShort)
		} else {
			PushToBulletSender(msg)
		}

		//fmt.Println("礼物-----", name, giftstring)
		// 总打赏高于x元，加一句大气
		if sumCost >= 50000 { // 50元
			PushToBulletSender(name + "老板大气大气")
		}
		delete(thanksGiver.giftNotBlindBoxTable, name)
	}

	//盲盒礼物
	for name, m := range thanksGiver.giftBlindBoxTable {
		giftstring := []string{}
		msg := ""
		for blindBoxName, blindBoxMap := range m {
			giftstring = append(giftstring, fmt.Sprintf("%d个%s盈亏%+.2f元", blindBoxMap["count"], blindBoxName, float64(blindBoxMap["profit_and_loss"])/1000))
			// 计算打赏金额
			// 感谢完后立刻清空map
			delete(m, blindBoxName)
		}

		msgShort := ""

		msg = "感谢" + name + "的"
		for k, v := range giftstring {
			if k == 0 {
				msg += v
				msgShort = v
			} else {
				msg += "，" + v
				msgShort += "，" + v
			}
		}

		ms := []rune(msg)

		if len(ms) > danmuLen {
			PushToBulletSender("感谢 " + name + " 的")
			PushToBulletSender(msgShort)
		} else {
			PushToBulletSender(msg)
		}
		delete(thanksGiver.giftBlindBoxTable, name)
	}
}
