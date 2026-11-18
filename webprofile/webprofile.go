// Package webprofile 统一对外请求的浏览器身份与特征头，
// 避免不同模块携带互相矛盾的 UA 或暴露 Go 默认客户端标志（Go-http-client/1.1）
package webprofile

import (
	"fmt"

	"github.com/xbclub/BilibiliDanmuRobot-Core/browserclient"
)

// buildUserAgent 用 utls 实际使用的 ClientHello profile 拼 UA，版本号不写死。
//
// 版本来源固定为 browserclient.ChromeProfile（= utls.HelloChrome_Auto，
// 随 utls 升级自动跟随新 Chrome）。UA 与 TLS 指纹读的是同一个常量，
// 由 TestUserAgentMatchesTLSProfile 断言两者一致，
// 因此不会出现「UA 说 Chrome/133、ClientHello 却是 140」这种自相矛盾。
//
// 真实 Chrome 的 UA 在副版本上做 UA Reduction（Chrome/<主版本>.0.0.0），
// 主版本与 ClientHello 大版本一致即可，副版本固定为 0.0.0 是真实行为。
func buildUserAgent() string {
	return fmt.Sprintf(
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s.0.0.0 Safari/537.36",
		ChromeVersion(),
	)
}

// UserAgent 全模块共用的浏览器身份（Chrome / macOS），由 buildUserAgent 从 utls 派生。
//
// 这里用 var 而非 func：既有调用点（api/utils.go、client/client.go）
// 把它当字符串值直接用。它**没有**第二个硬编码版本号，版本仍只来自
// browserclient.ChromeProfile。
//
// 包级 var 在 init 阶段求值；utls 包已完成初始化，profile 的 Version 已就绪。
var UserAgent = buildUserAgent()

// ChromeVersion 返回 utls profile 里的主版本号（形如 "133"），
// 供需要单独拼版本号的场景使用（如测试断言 UA 与指纹一致）。
// 这里不写死具体数字：值来自 browserclient.ChromeProfile（HelloChrome_Auto，
// 随 utls 升级自动跟随；utls v1.8.2 下 Auto 恰为 133）。
func ChromeVersion() string {
	return browserclient.ChromeProfile.Version
}

// SecChUA 返回 sec-ch-ua 头的值，主版本号同样从 utls 派生，不写死。
//
// 唯一的固定部分是 GREASE 品牌名（"Not_A Brand" / v=24）：utls 的 ClientHelloID
// 只暴露主版本号，不含品牌列表，无法派生。拿到新版 Chrome 的真实头时再校准这一处。
func SecChUA() string {
	v := ChromeVersion()
	return fmt.Sprintf(`"Chromium";v="%s", "Not_A Brand";v="24", "Google Chrome";v="%s"`, v, v)
}

// 通用特征头取值，与 Chrome XHR 行为一致
const (
	HeaderAccept     = "application/json, text/plain, */*"
	HeaderAcceptLang = "zh-CN,zh;q=0.9"

	// HeaderAcceptEncoding 显式声明，避免 Go 传输层自动只加 gzip（或 WS 一个都不加），
	// 导致「同一个客户端两种 Accept-Encoding」的矛盾。
	HeaderAcceptEncoding = "gzip, deflate, br, zstd"

	HeaderSecChUAMobile   = "?0"
	HeaderSecChUAPlatform = `"macOS"`

	// 注意：这里刻意**不**发 Sec-Fetch-User。
	// 真实 Chrome 只在「用户激活的导航请求」上带它，fetch/XHR 与 WS 握手都不带；
	// 任务书 §4.1-A 的清单里也没有它，多发一个头本身就是破绽。

	MainOrigin  = "https://www.bilibili.com"
	MainReferer = "https://www.bilibili.com/"
	LiveOrigin  = "https://live.bilibili.com"
	LiveReferer = "https://live.bilibili.com/"
)

// clientHints 返回三枚 Client Hints 头（Chrome 对 https 请求必发）。
func clientHints() map[string]string {
	return map[string]string{
		"sec-ch-ua":          SecChUA(),
		"sec-ch-ua-mobile":   HeaderSecChUAMobile,
		"sec-ch-ua-platform": HeaderSecChUAPlatform,
	}
}

// WSHeaders 返回 WebSocket 握手用的特征头（Origin 为直播域）。
//
// 取值对齐真实 Chrome 的 WS 握手：
//   - Sec-Fetch-Dest/Mode 均为 websocket、Site 为 same-site —— 这是浏览器强制头，
//     页面脚本无法伪造，缺失是最明显的破绽
//   - Pragma / Cache-Control: no-cache 是 Chrome 握手固带的
//
// 与 API 组不同，这里**不放** Accept / Referer：Chrome 发 WS 握手时不带这两个头，
// 补上反而与真实浏览器不符。
//
// Connection / Upgrade / Sec-WebSocket-Key / Sec-WebSocket-Version 由
// gorilla/websocket 自己写，不要在这里重复。
func WSHeaders() map[string]string {
	h := map[string]string{
		"User-Agent":      UserAgent,
		"Origin":          LiveOrigin,
		"Accept-Language": HeaderAcceptLang,
		"Accept-Encoding": HeaderAcceptEncoding,
		"Pragma":          "no-cache",
		"Cache-Control":   "no-cache",
		"Sec-Fetch-Dest":  "websocket",
		"Sec-Fetch-Mode":  "websocket",
		"Sec-Fetch-Site":  "same-site",
	}
	for k, v := range clientHints() {
		h[k] = v
	}
	return h
}

// LiveAPIHeaders 返回访问直播域（api.live.bilibili.com 等）时的完整特征头。
//
// 对齐真实 Chrome 从 live.bilibili.com 页面发起的 fetch/XHR：
// Sec-Fetch-Dest=empty、Mode=cors、Site=same-site。
func LiveAPIHeaders() map[string]string {
	h := map[string]string{
		"User-Agent":      UserAgent,
		"Accept":          HeaderAccept,
		"Accept-Language": HeaderAcceptLang,
		"Accept-Encoding": HeaderAcceptEncoding,
		"Origin":          LiveOrigin,
		"Referer":         LiveReferer,
		"Pragma":          "no-cache",
		"Cache-Control":   "no-cache",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-site",
	}
	for k, v := range clientHints() {
		h[k] = v
	}
	return h
}

// MainAPIHeaders 返回访问主站域（api.bilibili.com、passport.bilibili.com）时的完整特征头。
// 语义同 LiveAPIHeaders，只把 Origin/Referer 换成主站。
func MainAPIHeaders() map[string]string {
	h := map[string]string{
		"User-Agent":      UserAgent,
		"Accept":          HeaderAccept,
		"Accept-Language": HeaderAcceptLang,
		"Accept-Encoding": HeaderAcceptEncoding,
		"Origin":          MainOrigin,
		"Referer":         MainReferer,
		"Pragma":          "no-cache",
		"Cache-Control":   "no-cache",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-site",
	}
	for k, v := range clientHints() {
		h[k] = v
	}
	return h
}

// MainHeaders 兼容旧调用点，转发到 MainAPIHeaders，避免两套拼法漂移。
func MainHeaders() map[string]string {
	return MainAPIHeaders()
}

// LiveHeaders 兼容旧调用点，转发到 LiveAPIHeaders，避免两套拼法漂移。
func LiveHeaders() map[string]string {
	return LiveAPIHeaders()
}
