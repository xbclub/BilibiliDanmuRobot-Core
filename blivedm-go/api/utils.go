package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/xbclub/BilibiliDanmuRobot-Core/browserclient"
	"github.com/xbclub/BilibiliDanmuRobot-Core/webprofile"
)

func HttpGet(url string, headers *http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = *headers
	c := browserclient.Client()
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetJson(url string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// 直播域请求补齐浏览器特征头，避免暴露 Go 客户端身份
	for k, v := range webprofile.LiveHeaders() {
		req.Header.Set(k, v)
	}
	c := browserclient.Client()
	resp, err := c.Do(req)
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

func GetJsonWithHeader(url string, headers *http.Header, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// 统一浏览器身份；域名相关的 Origin/Referer 由调用方通过 headers 补充
	headers.Set("User-Agent", webprofile.UserAgent)
	headers.Set("Accept", webprofile.HeaderAccept)
	headers.Set("Accept-Language", webprofile.HeaderAcceptLang)
	req.Header = *headers
	c := browserclient.Client()
	resp, err := c.Do(req)
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
