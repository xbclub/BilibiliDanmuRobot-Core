package biliapi

import (
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/xbclub/BilibiliDanmuRobot-Core/browserclient"
	"github.com/xbclub/BilibiliDanmuRobot-Core/webprofile"
)

// UserAgent 由 webprofile 从 utls ClientHello profile 动态派生，非常量
var (
	UserAgent = webprofile.UserAgent
	userAgent = UserAgent
)

// cookieMu 保护 CookieStr / CookieList：发送 goroutine、GUI 查询、
// Logout/AddCookie 写方并发访问（包外请走 CookieHeader/CookieValue 访问器，
// 包内直接访问全局变量的位置也必须持锁）
var cookieMu sync.RWMutex

// 调用GetLoginInfo后，对全局变量cookie赋值
var CookieStr string
var CookieList = make(map[string]string)

// 全局客户端对象
var cli *resty.Client

func InitHttpClient() {
	// 全部请求走 Chrome TLS 指纹传输层
	cli = resty.New().SetTransport(browserclient.Transport())
}

// CookieHeader 线程安全地读取当前 cookie 头字符串
func CookieHeader() string {
	cookieMu.RLock()
	defer cookieMu.RUnlock()
	return CookieStr
}

// CookieValue 线程安全地读取指定 cookie 字段（不存在返回空串）
func CookieValue(key string) string {
	cookieMu.RLock()
	defer cookieMu.RUnlock()
	return CookieList[key]
}

// AddCookie 写入/覆盖一个 cookie 字段，并以浏览器格式（"k=v; k=v"，无尾分号）重建 CookieStr
func AddCookie(key, value string) {
	if value == "" {
		return
	}
	cookieMu.Lock()
	defer cookieMu.Unlock()
	CookieList[key] = value
	rebuildCookieStr()
}

func rebuildCookieStr() {
	parts := make([]string, 0, len(CookieList))
	for k, v := range CookieList {
		parts = append(parts, k+"="+v)
	}
	CookieStr = strings.Join(parts, "; ")
}

// mergeCookieStrIntoList 把 "k=v; k=v" 形式的 cookie 字符串并入 CookieList（不覆盖已有键）。
// 用于加载本地 token 时保留 txt 独有而 json 缺失的字段（如手工只更新了 txt 的场景）
func mergeCookieStrIntoList(raw string) {
	for _, part := range strings.Split(raw, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		if _, ok := CookieList[kv[0]]; !ok {
			CookieList[kv[0]] = kv[1]
		}
	}
}
