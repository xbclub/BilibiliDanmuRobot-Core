package config

import (
	"github.com/zeromicro/go-zero/core/logx"
)

type Config struct {
	//rest.RestConf
	Log logx.LogConf

	// 核心设置
	RoomId      int    `json:",default=4699397"`
	WsServerUrl string `json:",default=wss://broadcastlv.chat.bilibili.com:2245/sub"`

	// 常规设置
	DanmuLen     int    `json:",default=20"`    // 弹幕限制长度
	EntryMsg     string `json:",default=off"`   // 进房间自动发送的文本
	PKNotice     bool   `json:",default=true"`  // PK信息开关
	ShowBlockMsg bool   `json:",default=false"` // 禁言提醒开关
	GoodbyeInfo  string `json:",optional"`      // 下播自动发送的话

	// 关键字回复
	KeywordReply     bool              `json:",default=false"` //关键词回复开关
	KeywordReplyList map[string]string `json:",optional"`      // 关键词回复列表

	// AI聊天相关
	TalkRobotCmd  string   `json:",default=test"`                                // 机器人聊天关键字
	FuzzyMatchCmd bool     `json:",default=false"`                               // 模糊匹配关键字
	RobotName     string   `json:",default=花花"`                                 // 机器人名称
	RobotMode     string   `json:",default=QingYunKe,options=QingYunKe|ChatGPT"` // 机器人服务
	ChatGPT       struct { // GPT的配置
		APIUrl   string `json:",default=https://api.openai.com/v1"`
		APIToken string `json:",optional"`
		Prompt   string `json:",default=你是一个非常幽默的机器人助理，可以使用emoji表情符号，可以使用颜文字"`
		Limit    bool   `json:",default=true"`
		Model    string `json:",default=gpt-3.5-turbo"`
	}

	// 欢迎配置
	InteractWord       bool       `json:",default=false"`         // 欢迎弹幕开关
	WelcomeDanmu       []string   `json:",default='欢迎 {user} ~'"` // 欢迎语列表
	InteractWordByTime bool       `json:",default=false"`         // 按时间段欢迎
	WelcomeDanmuByTime []struct { // 分时段欢迎配置列表
		Enabled bool     `json:",optional"`      // 是否启用
		Key     string   `json:",optional"`      // 时间
		Random  bool     `json:",default=false"` // 是否随机
		Danmu   []string `json:",optional"`      // 内容列表
	} `json:",optional"`
	EntryEffect             bool              `json:",default=false"` // 特效欢迎开关
	WelcomeHighWealthy      bool              `json:",default=false"` // 欢迎高财富等级用户
	WelcomeHighWealthyLevel int               `json:",default=20"`    // 从多少级财富等级开始欢迎
	ThanksFocus             bool              `json:",default=false"` // 关注感谢开关
	ThanksShare             bool              `json:",default=false"` // 分享感谢开关
	InteractSelf            bool              `json:",default=true"`  // 欢迎自己
	FocusDanmu              []string          `json:",optional"`      // 关注的感谢列表
	WelcomeSwitch           bool              `json:",default=false"` // 指定欢迎开关
	WelcomeString           map[string]string `json:",optional"`      // 指定欢迎配置列表
	WelcomeBlacklistWide    []string          `json:",optional"`      // 不欢迎黑名单模糊匹配
	WelcomeBlacklist        []string          `json:",optional"`      // 不欢迎黑名单精确匹配

	// 答谢设置
	ThanksGift            bool `json:",default=false"` // 感谢送礼
	ThanksGiftTimeout     int  `json:",default=3"`     // 礼物统计时间
	ThanksBlindBoxTimeout int  `json:",default=6"`     // 盲盒统计时间
	ThanksMinCost         int  `json:",default=0"`     // 最小感谢礼物价值

	// 定时弹幕配置
	CronDanmu bool `json:",default=false"` // 定时弹幕开关
	// CronSupportSec bool `json:",default=false"`
	CronDanmuList []struct { // 定时弹幕配置列表
		Cron   string   `json:",optional"`      // 定时表达式
		Random bool     `json:",default=false"` // 是否随机发送
		Danmu  []string `json:",optional"`      // 内容列表
	} `json:",optional"`

	// 抽签设置
	DrawByLot    bool     `json:",default=true"` // 抽签开关
	DrawLotsList []string `json:",optional"`     // 抽签话术列表

	// 签到设置
	SignInEnable bool   `json:",default=true"` // 签到
	DBPath       string `json:",default=./db"`
	DBName       string `json:",default=sqliteDataBase.db"`

	// 杂项设置 GUI无界面配置
	CustomizeBullet bool `json:",default=false"` // 手动弹幕发送(命令行)	GUI不要有选项

	// 抽奖设置
	LotteryEnable bool   `json:",default=true"` // 抽奖开关
	LotteryUrl    string `json:",optional"`     // 抽奖地址
}
