package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/xbclub/BilibiliDanmuRobot-Core/browserclient"
	"github.com/xbclub/BilibiliDanmuRobot-Core/webprofile"
)

// HeadersForURL 返回 urlStr 对应域名的完整浏览器特征头。
//
// 直播域（live.bilibili.com / api.live.bilibili.com）→ LiveAPIHeaders，
// 其余（主站 api.bilibili.com、passport 等）→ MainAPIHeaders。
// 所有调用点都从这里取默认头，避免同一客户端出现两套行为；
// 需要 cookie 等额外头的调用方在返回值上追加即可。
func HeadersForURL(urlStr string) http.Header {
	headers := http.Header{}
	profile := webprofile.MainAPIHeaders()
	if u, err := url.Parse(urlStr); err == nil {
		host := strings.ToLower(u.Hostname())
		if strings.HasPrefix(host, "live.") || strings.HasPrefix(host, "api.live.") {
			profile = webprofile.LiveAPIHeaders()
		}
	}
	for k, v := range profile {
		headers.Set(k, v)
	}
	return headers
}

// HttpGet 用统一浏览器身份请求 url，返回响应体。
//
// headers 里的项会覆盖默认特征头（例如 cookie）；为 nil 时只用默认头。
func HttpGet(url string, headers *http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = HeadersForURL(url)
	if headers != nil {
		mergeHeader(req.Header, *headers)
	}
	resp, err := browserclient.Client().Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return io.ReadAll(resp.Body)
}

// GetJson 请求 url 并把响应体解析到 result（默认头按域名注入）。
func GetJson(url string, result interface{}) error {
	return GetJsonWithHeader(url, nil, result)
}

// GetJsonWithHeader 请求 url 并把响应体解析到 result。
//
// 默认特征头由 HeadersForURL 按域名统一注入（含 Accept-Encoding、Sec-Fetch-*、
// sec-ch-ua*、Pragma/Cache-Control），调用方传入的 headers（cookie、需要覆盖的
// Origin/Referer）在其上追加。这样 getDanmuInfo / getBuvid3A4 / room_init
// 三条路径的请求头完全一致，不再出现「同一客户端两套行为」。
func GetJsonWithHeader(url string, headers *http.Header, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header = HeadersForURL(url)
	if headers != nil {
		mergeHeader(req.Header, *headers)
	}
	resp, err := browserclient.Client().Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}

// mergeHeader 把 src 的每一项写入 dst（同键覆盖，其余原样追加）。
func mergeHeader(dst http.Header, src http.Header) {
	for k, vs := range src {
		dst.Del(k)
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}
