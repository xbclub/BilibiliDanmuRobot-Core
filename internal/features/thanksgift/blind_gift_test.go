package thanksgift

import (
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
)

func TestBlindBoxGiftThanks(t *testing.T) {
	send := &entity.SendGiftText{}
	send.Data.Uname = "巴多_Bado"
	send.Data.GiftName = "星光铃铛"
	send.Data.Price = 5200
	send.Data.Num = 1
	send.Data.BlindGift.OriginalGiftName = "幸运盲盒"
	send.Data.BlindGift.OriginalGiftPrice = 5000

	if got, want := blindBoxGiftThanks(send), "感谢巴多_Bado赠送的幸运盲盒，爆出星光铃铛，血赚了200电池"; got != want {
		t.Fatalf("blindBoxGiftThanks() = %q, want %q", got, want)
	}

	send.Data.Price = 4800
	if got, want := blindBoxGiftThanks(send), "感谢巴多_Bado赠送的幸运盲盒，爆出星光铃铛，血亏了200电池"; got != want {
		t.Fatalf("blindBoxGiftThanks() = %q, want %q", got, want)
	}
}
