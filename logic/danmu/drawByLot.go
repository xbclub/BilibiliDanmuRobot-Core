package danmu

import (
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"math/rand"
	"strings"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// 抽签过程函数
func DodrawByLotProcess(msg, username string, svcCtx *svc.ServiceContext) {
	// 判断抽签结果是否为空
	if svcCtx.Config.DrawLotsList != nil && len(svcCtx.Config.DrawLotsList) > 0 {
		if strings.Compare("抽签", msg) == 0 {
			// 随机选择抽签结果
			randomIndex := rand.Intn(len(svcCtx.Config.DrawLotsList))
			fanfanstr := svcCtx.Config.DrawLotsList[randomIndex]
			response := fmt.Sprintf("%v", fanfanstr)
			// 将抽签结果发送给弹幕发送器
			strMeg := fmt.Sprintf("抽签成功，%v的抽签结果为:",username)
			logic.PushToBulletSender(strMeg)
			logic.PushToBulletSender(response)
		}
	} else {
		// 如果抽签列表为空，返回提示信息
		response := "别抽签，抽我"
		logic.PushToBulletSender(response)
	}
}