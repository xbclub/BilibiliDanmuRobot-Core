package client

import (
	"strings"
	"testing"
)

func TestBuvidFromCookie(t *testing.T) {
	cases := []struct {
		name   string
		cookie string
		want   string
	}{
		{"优先 buvid3", "SESSDATA=x; buvid3=ABC; _uuid=DEF", "ABC"},
		{"无 buvid3 回退 _uuid", "SESSDATA=x; _uuid=DEF", "DEF"},
		{"buvid3 在串尾无分号", "SESSDATA=x; buvid3=ABC", "ABC"},
		{"都没有返回空", "SESSDATA=x", ""},
	}
	for _, c := range cases {
		if got := buvidFromCookie(c.cookie); got != c.want {
			t.Fatalf("%s: buvidFromCookie(%q) = %q, want %q", c.name, c.cookie, got, c.want)
		}
	}
}

func TestAppendBuvidCookies(t *testing.T) {
	// 两个都缺：补齐 buvid3 + buvid4
	got := appendBuvidCookies("SESSDATA=x", "B3", "B4")
	if got != "SESSDATA=x; buvid3=B3; buvid4=B4" {
		t.Fatalf("both missing: got %q", got)
	}
	// 已有 buvid3：只补 buvid4，不产生重复键
	got = appendBuvidCookies("SESSDATA=x; buvid3=OLD", "B3", "B4")
	if got != "SESSDATA=x; buvid3=OLD; buvid4=B4" {
		t.Fatalf("buvid3 exists: got %q", got)
	}
	// 两个都有：原样返回
	got = appendBuvidCookies("buvid3=A; buvid4=B", "B3", "B4")
	if got != "buvid3=A; buvid4=B" {
		t.Fatalf("both exist: got %q", got)
	}
	// 空 cookie 直接拼，不带前导分号
	got = appendBuvidCookies("", "B3", "B4")
	if got != "buvid3=B3; buvid4=B4" {
		t.Fatalf("empty cookie: got %q", got)
	}
	// b4 为空只补 buvid3
	got = appendBuvidCookies("SESSDATA=x;", "B3", "")
	if got != "SESSDATA=x; buvid3=B3" {
		t.Fatalf("b4 empty: got %q", got)
	}
	if strings.Count(got, "buvid3") != 1 {
		t.Fatalf("duplicate buvid3 keys: %q", got)
	}
}
