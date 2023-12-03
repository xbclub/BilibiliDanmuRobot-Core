package danmu

import (
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"math/rand"
	"strings"
)

func DodrawByLotProcess(msg, username string) {
	if strings.Compare("抽签", msg) == 0 {
		var strMeg string
		iRes := rand.Intn(25)
		var strRes string
		switch iRes {
		case 0:
			strRes = "上上签"
		case 1, 2, 3:
			strRes = "上中签"
		case 4, 5, 6, 7, 8:
			strRes = "上下签"
		case 9, 10, 11, 12, 13, 14, 15:
			strRes = "中上签"
		case 16, 17, 18, 19, 20, 21, 22, 23, 24:
			strRes = "中中签"
		}
		strMeg = fmt.Sprintf("%v, 结果是%v 哟。", username, strRes)
		logic.PushToBulletSender(strMeg)
	}
}
