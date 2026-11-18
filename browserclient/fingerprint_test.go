package browserclient

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

// TestFingerprintChrome 通过 tls.peet.ws 回显客户端 TLS 指纹，
// 验证 utls 传输层确实以 Chrome 指纹出现。
// 需要网络，默认跳过；运行方式：
//
//	BILI_FP_CHECK=1 go test ./browserclient/ -run TestFingerprintChrome -v
func TestFingerprintChrome(t *testing.T) {
	if os.Getenv("BILI_FP_CHECK") == "" {
		t.Skip("设置 BILI_FP_CHECK=1 后运行（需要外网访问 tls.peet.ws）")
	}

	resp, err := Client().Get("https://tls.peet.ws/api/all")
	if err != nil {
		t.Fatalf("请求失败：%v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("读取响应失败：%v", err)
	}

	var report struct {
		HTTPVersion string `json:"http_version"`
		TLS         struct {
			TLSVersionNegotiated string `json:"tls_version_negotiated"`
			JA3                  string `json:"ja3"`
			JA3Hash              string `json:"ja3_hash"`
			JA4                  string `json:"ja4"`
			Peetprint            string `json:"peetprint"`
		} `json:"tls"`
		HTTP2 struct {
			AkamaiFingerprint     string `json:"akamai_fingerprint"`
			AkamaiFingerprintHash string `json:"akamai_fingerprint_hash"`
		} `json:"http2"`
	}
	if err := json.Unmarshal(body, &report); err != nil {
		t.Fatalf("解析失败：%v\nbody: %s", err, body)
	}

	t.Logf("TLS 版本: %s (772=TLS1.3)", report.TLS.TLSVersionNegotiated)
	t.Logf("HTTP 版本: %s", report.HTTPVersion)
	t.Logf("JA3:  %s (%s)", report.TLS.JA3, report.TLS.JA3Hash)
	t.Logf("JA4:  %s", report.TLS.JA4)
	t.Logf("peetprint: %s", report.TLS.Peetprint)
	t.Logf("akamai h2 指纹: %s (%s)", report.HTTP2.AkamaiFingerprint, report.HTTP2.AkamaiFingerprintHash)

	// 硬断言：TLS1.3 + h2 + JA4 的 Chrome 特征（t13d 前缀、首段以 h2 结尾）
	if report.TLS.TLSVersionNegotiated != "772" {
		t.Errorf("TLS 版本为 %s，期望 772 (TLS1.3)", report.TLS.TLSVersionNegotiated)
	}
	if report.HTTPVersion != "h2" {
		t.Errorf("HTTP 版本为 %s，期望 h2", report.HTTPVersion)
	}
	head := strings.Split(report.TLS.JA4, "_")[0]
	if !strings.HasPrefix(head, "t13d") || !strings.HasSuffix(head, "h2") {
		t.Errorf("JA4 = %s，首段 %s 不符合 Chrome 的 t13d...h2 特征", report.TLS.JA4, head)
	}
}
