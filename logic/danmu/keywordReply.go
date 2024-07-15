package danmu

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"strings"
)

func KeywordReply(danmu string, svcCtx *svc.ServiceContext, reply ...*entity.DanmuMsgTextReplyInfo) {
	if svcCtx.Config.KeywordReplyList != nil &&
		len(svcCtx.Config.KeywordReplyList) > 0 {
		for k, v := range svcCtx.Config.KeywordReplyList {
			if strings.Contains(danmu, k) {
				logic.PushToBulletSender(v, reply...)
				break
			}
		}
	}
}
