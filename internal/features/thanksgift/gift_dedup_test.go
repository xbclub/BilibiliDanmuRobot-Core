package thanksgift

import (
	"testing"
	"time"
)

func TestGiftTidDedup(t *testing.T) {
	t.Cleanup(resetGiftTidSeen)
	t0 := time.Now()
	if giftTidDuplicateAt("tid-1", t0) {
		t.Fatal("首次出现的 tid 不应判重")
	}
	if !giftTidDuplicateAt("tid-1", t0.Add(time.Second)) {
		t.Fatal("窗口内重复 tid 应判重")
	}
	if giftTidDuplicateAt("tid-2", t0.Add(time.Second)) {
		t.Fatal("不同 tid 不应判重")
	}
	// 空 tid 不去重（返回 false 且不记录，两次调用都放行）
	if giftTidDuplicateAt("", t0) || giftTidDuplicateAt("", t0) {
		t.Fatal("空 tid 不应判重")
	}
	// 窗口过后同 tid 可再次通过，且过期条目被清理
	if giftTidDuplicateAt("tid-1", t0.Add(61*time.Second)) {
		t.Fatal("窗口过期后 tid 应放行")
	}
	giftTidMu.Lock()
	if _, ok := giftTidSeen["tid-2"]; ok {
		giftTidMu.Unlock()
		t.Fatal("过期条目应被清理")
	}
	giftTidMu.Unlock()
}
