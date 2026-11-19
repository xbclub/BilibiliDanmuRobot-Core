package biliapi

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestSetHistoryCookieOldTokenCompat 验证老版本写入的 token 文件仍能正确加载：
// 老 txt 为 "k=v;k=v;"（无空格、尾分号）格式，老 json 与 txt 同源。
// 新逻辑：合并 txt 独有字段后以浏览器格式 "k=v; k=v" 重建 CookieStr
func TestSetHistoryCookieOldTokenCompat(t *testing.T) {
	oldTxt := "SESSDATA=abc123;bili_jct=xyz;DedeUserID=42;buvid3=AB==CDinfoc;"
	oldJSON := `{"SESSDATA":"abc123","bili_jct":"xyz"}`

	cookieMu.Lock()
	CookieStr = oldTxt
	CookieList = make(map[string]string)
	if err := json.Unmarshal([]byte(oldJSON), &CookieList); err != nil {
		cookieMu.Unlock()
		t.Fatal(err)
	}
	mergeCookieStrIntoList(CookieStr)
	rebuildCookieStr()
	cookieMu.Unlock()
	t.Cleanup(func() {
		cookieMu.Lock()
		CookieList = make(map[string]string)
		CookieStr = ""
		cookieMu.Unlock()
	})

	// json 与 txt 共有的字段 + txt 独有的字段都必须保留
	want := map[string]string{
		"SESSDATA":   "abc123",
		"bili_jct":   "xyz",
		"DedeUserID": "42",
		"buvid3":     "AB==CDinfoc", // 值内含 '=' 不能被截断
	}
	cookieMu.RLock()
	defer cookieMu.RUnlock()
	if len(CookieList) != len(want) {
		t.Errorf("字段数 %d != %d: %v", len(CookieList), len(want), CookieList)
	}
	for k, v := range want {
		if CookieList[k] != v {
			t.Errorf("字段 %s = %q, 期望 %q", k, CookieList[k], v)
		}
	}

	// CookieStr 为浏览器格式：无尾分号、分号后带空格，可无损解析回相同键值
	if strings.HasSuffix(CookieStr, ";") {
		t.Errorf("CookieStr 不应有尾分号: %q", CookieStr)
	}
	got := make(map[string]string)
	for _, part := range strings.Split(CookieStr, "; ") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			t.Fatalf("CookieStr 分段无法解析: %q", part)
		}
		got[kv[0]] = kv[1]
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("CookieStr 回解析 %s = %q, 期望 %q", k, got[k], v)
		}
	}
}
