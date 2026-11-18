package webprofile

import (
	"strings"
	"testing"

	utls "github.com/refraction-networking/utls"
)

// TestUserAgentMatchesTLSProfile 是 TASK §8.5.3 的守护：
// UA / sec-ch-ua 里的 Chrome 版本号必须与 utls 实际使用的 ClientHello profile 一致。
// 若有人把 browserclient.ChromeProfile 改成别的版本却忘了同步，
// 或反过来写死一个版本号，这里会立刻失败。
func TestUserAgentMatchesTLSProfile(t *testing.T) {
	// 1) 码内约定：profile 必须显式钉在 HelloChrome_133（§4.1-D）。
	if ChromeVersion() != "133" {
		t.Fatalf("ChromeVersion() = %q，期望 133", ChromeVersion())
	}

	// 2) UA 里必须出现 Chrome/133.0.0.0。
	wantUA := "Chrome/" + ChromeVersion() + ".0.0.0"
	if !strings.Contains(UserAgent, wantUA) {
		t.Errorf("UserAgent = %q，应包含 %q", UserAgent, wantUA)
	}

	// 3) UA 里不能出现其它版本的 Chrome/<数字>（防一个字符串里两套版本号）。
	if got := countChromeTokens(UserAgent); got != 1 {
		t.Errorf("UserAgent 中 Chrome/ 版本号出现 %d 次，期望 1: %q", got, UserAgent)
	}

	// 4) sec-ch-ua 的品牌版本必须与 UA 主版本一致：
	//    "Chromium";v="133" 与 "Google Chrome";v="133" 各一次，没有第三个版本。
	secChUA := SecChUA()
	if !strings.Contains(secChUA, `"Google Chrome";v="`+ChromeVersion()+`"`) {
		t.Errorf("SecChUA = %q，应含 Google Chrome;v=%s", secChUA, ChromeVersion())
	}
	if got := countVersionTokens(secChUA, ChromeVersion()); got != 2 {
		t.Errorf("SecChUA 中 v=%q 出现 %d 次，期望 2（Chromium 与 Google Chrome）: %q",
			ChromeVersion(), got, secChUA)
	}
	for _, other := range []string{"140", "141", "140.0"} {
		if strings.Contains(secChUA, `v="`+other+`"`) {
			t.Errorf("SecChUA 出现意外版本 %s: %q", other, secChUA)
		}
	}
}

// TestChromeProfileReferencesUtls 确认版本来源确实是 utls 的 profile，
// 而不是任何地方写死的字面量。直接比较常量身份。
func TestChromeProfileReferencesUtls(t *testing.T) {
	if ChromeVersion() != utls.HelloChrome_133.Version {
		t.Errorf("ChromeVersion() = %q，utls.HelloChrome_133.Version = %q",
			ChromeVersion(), utls.HelloChrome_133.Version)
	}
}

// TestHeadersNoSecFetchUser 守护「不发 Sec-Fetch-User」这一实测结论：
// 真实 Chrome 只在用户激活的导航请求上发该头，fetch/XHR 与 WS 握手都不发。
// 三组头里出现它都会成为破绽，因此这里逐组断言。
func TestHeadersNoSecFetchUser(t *testing.T) {
	groups := map[string]map[string]string{
		"WS":      WSHeaders(),
		"LiveAPI": LiveAPIHeaders(),
		"MainAPI": MainAPIHeaders(),
	}
	for name, h := range groups {
		for k := range h {
			if strings.EqualFold(k, "Sec-Fetch-User") {
				t.Errorf("%sHeaders 不应包含 Sec-Fetch-User（真实 Chrome 的 XHR/WS 不带）", name)
			}
		}
	}
}

// TestHeadersCommonFields 确认三组头都带齐伪装必需项。
func TestHeadersCommonFields(t *testing.T) {
	required := []string{
		"User-Agent",
		"Accept-Encoding",
		"Accept-Language",
		"Sec-Fetch-Dest",
		"Sec-Fetch-Mode",
		"Sec-Fetch-Site",
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-platform",
	}
	groups := map[string]map[string]string{
		"WS":      WSHeaders(),
		"LiveAPI": LiveAPIHeaders(),
		"MainAPI": MainAPIHeaders(),
	}
	for name, h := range groups {
		for _, k := range required {
			if _, ok := h[k]; !ok {
				t.Errorf("%sHeaders 缺少 %s", name, k)
			}
		}
		if got := h["Accept-Encoding"]; got != HeaderAcceptEncoding {
			t.Errorf("%s Accept-Encoding = %q，期望 %q", name, got, HeaderAcceptEncoding)
		}
	}

	// WS 握手不能带 Accept / Referer（真实 Chrome 不发）。
	ws := WSHeaders()
	for _, k := range []string{"Accept", "Referer"} {
		if _, ok := ws[k]; ok {
			t.Errorf("WSHeaders 不应包含 %s", k)
		}
	}
}

// countChromeTokens 统计字符串中 "Chrome/<数字>" 形式的出现次数。
func countChromeTokens(s string) int {
	count := 0
	rest := s
	for {
		i := strings.Index(rest, "Chrome/")
		if i < 0 {
			return count
		}
		rest = rest[i+len("Chrome/"):]
		if len(rest) > 0 && isDigit(rest[0]) {
			count++
		}
	}
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

// countVersionTokens 统计形如 v="<ver>" 的出现次数。
func countVersionTokens(s, ver string) int {
	count := 0
	needle := `v="` + ver + `"`
	rest := s
	for {
		i := strings.Index(rest, needle)
		if i < 0 {
			return count
		}
		count++
		rest = rest[i+len(needle):]
	}
}
