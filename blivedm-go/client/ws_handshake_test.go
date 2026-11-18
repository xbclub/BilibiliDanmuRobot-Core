package client

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// newLocalUpgrader 返回接受任意 Origin 的 upgrader（本地测试用）。
func newLocalUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}
}

// TestWSHandshakeHeaders 抓取真实 WS 握手请求头（TASK §8.5.1）：
// 必须含 Sec-Fetch-Dest/Mode、sec-ch-ua*、Accept-Encoding、Accept-Language，
// 且不得含 Sec-Fetch-User、Accept、Referer。
func TestWSHandshakeHeaders(t *testing.T) {
	var got http.Header
	up := newLocalUpgrader()
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		conn, err := up.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		_ = conn.Close()
	}))
	defer srv.Close()

	// 用 httptest 的自签证书走本地 TLS，目的只是抓 HTTP/1.1 upgrade 头；
	// TLS 层本身由 browserclient 在真实连接中提供（指纹相关测试见 browserclient）。
	d := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	url := "wss" + strings.TrimPrefix(srv.URL, "https")
	conn, _, err := d.Dial(url, wsHeaders())
	if err != nil {
		t.Fatalf("dial local wss: %v", err)
	}
	_ = conn.Close()

	required := map[string]string{
		"Sec-Fetch-Dest":     "websocket",
		"Sec-Fetch-Mode":     "websocket",
		"Sec-Fetch-Site":     "same-site",
		"Accept-Encoding":    "gzip, deflate, br, zstd",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Pragma":             "no-cache",
		"Cache-Control":      "no-cache",
		"Origin":             "https://live.bilibili.com",
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": `"macOS"`,
	}
	for k, want := range required {
		if v := got.Get(k); v != want {
			t.Errorf("握手头 %s = %q，期望 %q", k, v, want)
		}
	}
	for _, k := range []string{"Sec-Ch-Ua", "User-Agent", "Upgrade", "Connection", "Sec-Websocket-Key"} {
		if got.Get(k) == "" {
			t.Errorf("握手头缺少 %s", k)
		}
	}
	// 不应出现的头（Chrome 的 WS 握手不发）。
	for _, k := range []string{"Sec-Fetch-User", "Accept", "Referer"} {
		if v := got.Get(k); v != "" {
			t.Errorf("握手头不应包含 %s（实际 %q）", k, v)
		}
	}
	// 不能重复（gorilla 会直接报错，这里确认过滤逻辑生效）。
	if n := len(got.Values("Sec-Websocket-Version")); n != 1 {
		t.Errorf("Sec-WebSocket-Version 出现 %d 次，期望 1", n)
	}
}
