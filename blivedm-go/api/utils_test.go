package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xbclub/BilibiliDanmuRobot-Core/webprofile"
)

// TestHeadersForURLDomainSplit 验证按域名选头：直播域与主站域各自成组，
// 且都带齐 Sec-Fetch-* / sec-ch-ua* / Accept-Encoding。
func TestHeadersForURLDomainSplit(t *testing.T) {
	live := HeadersForURL("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=6")
	main := HeadersForURL("https://api.bilibili.com/x/web-interface/nav")
	roomInit := HeadersForURL("https://api.live.bilibili.com/room/v1/Room/room_init?id=6")

	// getDanmuInfo 与 room_init 属于同一域，头必须完全一致 ——
	// 历史上这里曾出现「同一客户端两套行为」，是最明显的破绽。
	if diff := headerDiff(live, roomInit); diff != "" {
		t.Errorf("getDanmuInfo 与 room_init 的请求头不一致: %s", diff)
	}

	for name, h := range map[string]http.Header{"live": live, "main": main} {
		for _, k := range []string{
			"User-Agent", "Accept", "Accept-Language", "Accept-Encoding",
			"Sec-Fetch-Dest", "Sec-Fetch-Mode", "Sec-Fetch-Site",
			"sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform",
			"Pragma", "Cache-Control", "Origin", "Referer",
		} {
			if h.Get(k) == "" {
				t.Errorf("%s 域请求头缺少 %s", name, k)
			}
		}
		if got := h.Get("Accept-Encoding"); got != webprofile.HeaderAcceptEncoding {
			t.Errorf("%s 的 Accept-Encoding = %q，期望 %q", name, got, webprofile.HeaderAcceptEncoding)
		}
	}

	// Origin/Referer 必须按域区分。
	if got := live.Get("Origin"); got != webprofile.LiveOrigin {
		t.Errorf("live Origin = %q，期望 %q", got, webprofile.LiveOrigin)
	}
	if got := main.Get("Origin"); got != webprofile.MainOrigin {
		t.Errorf("main Origin = %q，期望 %q", got, webprofile.MainOrigin)
	}
}

// TestGetJsonWithHeaderMergesCallerHeaders 用 httptest 抓真实出站请求头，
// 确认调用方补的 cookie/Origin 不会把默认浏览器头覆盖成缺项。
func TestGetJsonWithHeaderMergesCallerHeaders(t *testing.T) {
	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0}`))
	}))
	defer srv.Close()

	extra := &http.Header{}
	extra.Set("cookie", "SESSDATA=abc; bili_jct=def")
	var out struct {
		Code int `json:"code"`
	}
	if err := GetJsonWithHeader(srv.URL, extra, &out); err != nil {
		t.Fatalf("GetJsonWithHeader: %v", err)
	}

	if got.Get("Cookie") != "SESSDATA=abc; bili_jct=def" {
		t.Errorf("Cookie = %q", got.Get("Cookie"))
	}
	// 默认特征头仍在（未被调用方覆盖丢）。
	for _, k := range []string{"User-Agent", "Accept-Encoding", "Sec-Fetch-Dest", "sec-ch-ua"} {
		if got.Get(k) == "" {
			t.Errorf("出站请求缺少默认头 %s，实际头: %v", k, got)
		}
	}
	// 与 room_init（GetJson）走同一套默认头：这里用同 URL 再取一次默认头比对。
	want := HeadersForURL(srv.URL)
	if diff := headerDiff(got, want); diff != "" {
		t.Errorf("GetJsonWithHeader 出站头与默认头不一致: %s", diff)
	}
}

// headerDiff 返回 want 中有、got 中缺失或值不同的项（小写键名），空串表示一致。
func headerDiff(got, want http.Header) string {
	var diffs []string
	for k, vs := range want {
		gv := got.Get(k)
		if gv == "" || gv != strings.Join(vs, ", ") {
			diffs = append(diffs, k+"="+gv)
		}
	}
	return strings.Join(diffs, "; ")
}
