package danmuProcess

import (
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
)

// 所有弹幕处理类
type DanmuProcessClass struct {
	GptClass                                 Gpt
	DrawByLotClass                           DrawByLot
	ForeignLanguageTranslationInChineseClass ForeignLanguageTranslationInChinese
	TraditionalToSimplifiedConversionClass   TraditionalToSimplifiedConversion
}

type DanmuProcess interface {
	Create() DanmuProcess
	// 弹幕处理函数
	DoDanmuProcess()
	SetConfig(svcCtx *svc.ServiceContext)
	SetDanmu(content *string, user *message.User)
}
