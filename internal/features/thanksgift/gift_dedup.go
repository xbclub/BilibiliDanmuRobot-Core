package thanksgift

import (
	"sync"
	"time"
)

// giftTidWindow V1/V2 双事件 tid 去重窗口：过渡期同一礼物会同时经
// SEND_GIFT（JSON）与 SEND_GIFT_V2（pb）下发，两版 tid 为同一交易号；
// 双发间隔通常在毫秒级，60s 窗口留足余量
const giftTidWindow = 60 * time.Second

var (
	giftTidMu   sync.Mutex
	giftTidSeen = map[string]time.Time{}
)

// giftTidDuplicate 标记并判断该 tid 是否已在窗口内出现过：
// 首次返回 false（已记录），窗口内重复返回 true；tid 为空不去重
// （避免无 tid 的事件被互相误杀）
func giftTidDuplicate(tid string) bool {
	return giftTidDuplicateAt(tid, time.Now())
}

func giftTidDuplicateAt(tid string, now time.Time) bool {
	if tid == "" {
		return false
	}
	giftTidMu.Lock()
	defer giftTidMu.Unlock()
	if t, ok := giftTidSeen[tid]; ok && now.Sub(t) < giftTidWindow {
		return true
	}
	for k, v := range giftTidSeen {
		if now.Sub(v) >= giftTidWindow {
			delete(giftTidSeen, k)
		}
	}
	giftTidSeen[tid] = now
	return false
}
