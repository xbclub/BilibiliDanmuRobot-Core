package biliapi

import (
	"sync/atomic"

	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/zeromicro/go-zero/core/logx"
)

// dryRun 弹幕发送出口的 dry-run 开关（进程内原子布尔）。
// 刻意不做成配置项：仅供 GUI/测试在进程内显式调用开启，
// 配置文件里不会出现这个开关，避免线上误开。
var dryRun atomic.Bool

// EnableDryRun 开启发送出口 dry-run：此后所有弹幕不真正请求 B 站，
// 只把「将要发送的内容」打到日志（前缀 [dry-run]），Send 一律返回成功，
// 发送队列按正常 3.5s 节奏继续排空，可完整模拟运行时序。
// GUI 调试用法：NewWsHandler 之前调用一次即可。
func EnableDryRun() {
	dryRun.Store(true)
	logx.Info("[dry-run] 弹幕发送出口已切换为 dry-run 模式（仅记日志，不对外发送）")
}

// DisableDryRun 关闭 dry-run，恢复真实发送
func DisableDryRun() {
	dryRun.Store(false)
	logx.Info("[dry-run] 弹幕发送出口已恢复真实发送")
}

// DryRunEnabled 返回当前是否处于 dry-run 模式
func DryRunEnabled() bool {
	return dryRun.Load()
}

// logDryRunSend 统一的 dry-run 日志格式：消息体 + @回复目标
func logDryRunSend(msg string, reply ...*entity.DanmuMsgTextReplyInfo) {
	if len(reply) > 0 {
		logx.Infof("[dry-run] 弹幕未发送: %q reply=%v", msg, reply[0])
		return
	}
	logx.Infof("[dry-run] 弹幕未发送: %q", msg)
}
