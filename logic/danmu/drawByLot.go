package danmu

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"math/rand"
	"strings"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// 抽签过程函数
func DodrawByLotProcess(msg, username string, svcCtx *svc.ServiceContext, reply ...*entity.DanmuMsgTextReplyInfo) {
	// 判断抽签结果是否为空
	if strings.Compare("抽签", msg) == 0 {
		if svcCtx.Config.DrawLotsList != nil && len(svcCtx.Config.DrawLotsList) > 0 {
			// 随机选择抽签结果
			randomIndex := rand.Intn(len(svcCtx.Config.DrawLotsList))
			logic.PushToBulletSender(svcCtx.Config.DrawLotsList[randomIndex], reply...)
		} else {
			// 如果抽签列表为空，返回提示信息
			response := "别抽签，抽主播!"
			logic.PushToBulletSender(response, reply...)
		}
	}
}