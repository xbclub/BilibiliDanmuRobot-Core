package biliapi

import (
	"os"
	"testing"
)

// TestLogout 离网验证退出登录：本地凭证清理 + 进程内状态清空 + 幂等
// （不种 bili_jct，跳过服务端注销分支，避免测试触网）
func TestLogout(t *testing.T) {
	// 模拟登录态
	CookieList["SESSDATA"] = "fake-sessdata"
	CookieStr = "SESSDATA=fake-sessdata"
	if LoggedIn() {
		t.Fatal("缺少 bili_jct 时 LoggedIn 应为 false")
	}

	// 造本地凭证文件（包目录下 token/）
	if err := os.MkdirAll("token", 0755); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"token/bili_token.txt", "token/bili_token.json"} {
		if err := os.WriteFile(p, []byte("x"), 0666); err != nil {
			t.Fatal(err)
		}
	}

	if err := Logout(); err != nil {
		t.Fatalf("Logout 应成功: %v", err)
	}

	if CookieStr != "" || len(CookieList) != 0 {
		t.Fatalf("会话状态未清空: CookieStr=%q CookieList=%v", CookieStr, CookieList)
	}
	for _, p := range []string{"token/bili_token.txt", "token/bili_token.json"} {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Fatalf("凭证文件未删除: %s", p)
		}
	}

	// 幂等：未登录状态再次调用安全
	if err := Logout(); err != nil {
		t.Fatalf("重复 Logout 应幂等成功: %v", err)
	}
}

func TestLoggedIn(t *testing.T) {
	if LoggedIn() {
		t.Fatal("空状态不应视为已登录")
	}
	CookieList["bili_jct"] = "jct"
	CookieList["SESSDATA"] = "sess"
	if !LoggedIn() {
		t.Fatal("bili_jct+SESSDATA 齐备应视为已登录")
	}
	// 还原，避免污染其他测试
	CookieList = make(map[string]string)
	CookieStr = ""
}
