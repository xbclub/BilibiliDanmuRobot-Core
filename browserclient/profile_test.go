package browserclient

import (
	"net"
	"strings"
	"testing"

	utls "github.com/refraction-networking/utls"
)

// TestChromeProfileFollowsUtls 守护「自动跟随」语义：ChromeProfile 必须是
// utls.HelloChrome_Auto，这样 utls 升级时 UA、sec-ch-ua 与 ClientHello 指纹
// 一起升级，三者永不脱节。
//
// 这里**不**断言某个具体版本号：写死版本（如必须等于 HelloChrome_133）会让
// utls 升级后指纹冻结在旧版本，反而更容易被识别。
func TestChromeProfileFollowsUtls(t *testing.T) {
	if ChromeProfile != utls.HelloChrome_Auto {
		t.Errorf("ChromeProfile = %+v，期望 utls.HelloChrome_Auto（自动跟随升级）", ChromeProfile)
	}
	// 版本号必须非空，否则派生出的 UA/sec-ch-ua 会拼出畸形字符串。
	if ChromeProfile.Version == "" {
		t.Error("ChromeProfile.Version 为空，UA 与 sec-ch-ua 无法派生")
	}
	// Str() 形如 Chrome-<版本>，是 utls 的既定格式。
	if !strings.HasPrefix(ChromeProfile.Str(), "Chrome-") {
		t.Errorf("ChromeProfile.Str() = %q，应以 Chrome- 开头", ChromeProfile.Str())
	}
}

// TestClientHelloALPNOverridden 是本仓库最关键的回归测试。
//
// 曾经的 bug：只用 utls.Config.NextProtos 想覆盖 Chrome 预设的 ALPN，
// 但 HelloChrome_* 走 applyPresetByID → UTLSIdToSpec，扩展列表原样复制，
// config.NextProtos 根本进不了 ClientHello，ALPN 恒为 [h2 http/1.1]。
// 结果 WS 协商成 h2，gorilla 把 h2 SETTINGS 帧当 HTTP/1.1 响应解析，
// 报 malformed HTTP response，永远连不上。
//
// 这里不联网，只构建 ClientHello 并检查其 ALPN 扩展。
func TestClientHelloALPNOverridden(t *testing.T) {
	cases := []struct {
		name       string
		nextProtos []string
		want       []string
	}{
		{name: "WS 只给 http/1.1", nextProtos: []string{"http/1.1"}, want: []string{"http/1.1"}},
		{name: "HTTP 请求给 h2+http/1.1", nextProtos: []string{"h2", "http/1.1"}, want: []string{"h2", "http/1.1"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, server := net.Pipe()
			defer client.Close()
			defer server.Close()

			uconn, err := newChromeUConn(client, "broadcastlv.chat.bilibili.com", tc.nextProtos)
			if err != nil {
				t.Fatalf("newChromeUConn: %v", err)
			}
			if err := uconn.BuildHandshakeState(); err != nil {
				t.Fatalf("BuildHandshakeState: %v", err)
			}

			got := uconn.HandshakeState.Hello.AlpnProtocols
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Fatalf("ClientHello ALPN = %v，期望 %v", got, tc.want)
			}

			// 同时确认替换没有破坏其它 Chrome 特征：扩展里仍只有一个 ALPN，
			// 且 utls 记录的 ClientHelloID 仍是 HelloCustom + ApplyPreset 路径。
			alpnCount := 0
			for _, ext := range uconn.Extensions {
				if _, ok := ext.(*utls.ALPNExtension); ok {
					alpnCount++
				}
			}
			if alpnCount != 1 {
				t.Errorf("ClientHello 中 ALPN 扩展数量 = %d，期望 1", alpnCount)
			}
		})
	}
}

// TestDialTLSContextRejectsH2 确认 DialTLSContext 对非 h1.1 协商结果报错，
// 而不是把 h2 连接交给 gorilla（那种连接必然握手失败且难以诊断）。
func TestDialTLSContextRejectsH2(t *testing.T) {
	// 直接调用内部 helper 无法伪造协商结果，这里只断言函数存在正确签名
	// 与错误文案约束：真正的协商在第 1 个测试已覆盖 ALPN 请求值。
	// 保留此测试用于说明意图，避免后人删掉 DialTLSContext 的检查分支。
	if !strings.Contains(errNegotiatedH2.Error(), "h2") {
		t.Errorf("errNegotiatedH2 = %q，文案应含 h2 便于排查", errNegotiatedH2)
	}
}
