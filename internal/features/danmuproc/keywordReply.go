package danmuproc

import (
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"strings"
)

// KeywordReply 关键词回复，命中第一个关键词即回复
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
