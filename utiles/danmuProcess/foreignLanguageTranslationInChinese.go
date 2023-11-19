package danmuProcess

import (
	"fmt"
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/pemistahl/lingua-go"
	"github.com/ulinoyaped/BaiduTranslate"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"strings"
)

var languages = []lingua.Language{
	lingua.English,
	lingua.French,
	lingua.German,
	lingua.Spanish,
	lingua.Chinese,
	lingua.Japanese,
	lingua.Russian,
	lingua.Korean,
}

type ForeignLanguageTranslationInChinese struct {
	danmuContent *string
	fromUser     *message.User
	svcCtx       *svc.ServiceContext
	baiduInfo    BaiduTranslate.BaiduInfo
	detector     lingua.LanguageDetector
}

func (foreignLanguageTranslationInChinese *ForeignLanguageTranslationInChinese) Create() DanmuProcess {
	return new(ForeignLanguageTranslationInChinese)
}

func (foreignLanguageTranslationInChinese *ForeignLanguageTranslationInChinese) DoDanmuProcess() {
	language, exists := foreignLanguageTranslationInChinese.detector.DetectLanguageOf(*foreignLanguageTranslationInChinese.danmuContent)

	if exists && language != lingua.Chinese && len(foreignLanguageTranslationInChinese.baiduInfo.AppID) >= 0 {
		res := foreignLanguageTranslationInChinese.baiduInfo.NormalTr(*foreignLanguageTranslationInChinese.danmuContent, "auto", "zh")
		if res.Dst != "" && strings.Compare(res.Dst, *foreignLanguageTranslationInChinese.danmuContent) != 0 {
			strMeg := fmt.Sprintf("%v：%v", foreignLanguageTranslationInChinese.fromUser.Uname, res.Dst)
			logic.PushToBulletSender(strMeg)
		}
	}
}

func (foreignLanguageTranslationInChinese *ForeignLanguageTranslationInChinese) SetConfig(svcCtx *svc.ServiceContext) {
	foreignLanguageTranslationInChinese.svcCtx = svcCtx
	foreignLanguageTranslationInChinese.detector = lingua.NewLanguageDetectorBuilder().FromLanguages(languages...).Build()
	foreignLanguageTranslationInChinese.baiduInfo.AppID = svcCtx.Config.ForeignLanguageTranslationInChinese.AppID
	foreignLanguageTranslationInChinese.baiduInfo.SecretKey = svcCtx.Config.ForeignLanguageTranslationInChinese.SecretKey
}

func (foreignLanguageTranslationInChinese *ForeignLanguageTranslationInChinese) SetDanmu(content *string, user *message.User) {
	foreignLanguageTranslationInChinese.danmuContent = content
	foreignLanguageTranslationInChinese.fromUser = user
}
