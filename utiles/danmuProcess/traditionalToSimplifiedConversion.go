package danmuProcess

import (
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/liuzl/gocc"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type TraditionalToSimplifiedConversion struct {
	danmuContent *string
	fromUser     *message.User
	svcCtx       *svc.ServiceContext
	t2s          *gocc.OpenCC
}

func (traditionalToSimplifiedConversion *TraditionalToSimplifiedConversion) DoDanmuProcess() {
	if traditionalToSimplifiedConversion.t2s != nil {
		out, err := traditionalToSimplifiedConversion.t2s.Convert(*traditionalToSimplifiedConversion.danmuContent)
		if err != nil {
			logx.Infof(err.Error())
		}

		if err == nil && out != *traditionalToSimplifiedConversion.danmuContent {
			logic.PushToBulletSender(out)
		}
	}
}

func (traditionalToSimplifiedConversion *TraditionalToSimplifiedConversion) Create() DanmuProcess {
	return new(TraditionalToSimplifiedConversion)
}

func (traditionalToSimplifiedConversion *TraditionalToSimplifiedConversion) SetConfig(svcCtx *svc.ServiceContext) {
	traditionalToSimplifiedConversion.svcCtx = svcCtx
	t2s, err := gocc.New("t2s")
	if err != nil {
		logx.Infof(err.Error())
		return
	}
	traditionalToSimplifiedConversion.t2s = t2s
}

func (traditionalToSimplifiedConversion *TraditionalToSimplifiedConversion) SetDanmu(content *string, user *message.User) {
	traditionalToSimplifiedConversion.danmuContent = content
	traditionalToSimplifiedConversion.fromUser = user
}
