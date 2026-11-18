// Package browserclient 提供统一使用 utls 伪装 Chrome TLS 指纹（JA3/JA4）的 HTTP 传输层，
// 供 resty、net/http 直连与 WebSocket 握手共用，避免暴露 Go 标准库的 TLS 指纹。
//
// 已知不完美（Phase A 范围）：HTTP/2 层仍为 x/net/http2 的默认 SETTINGS 帧
// 与 header 顺序，与 Chrome 存在差异；TLS ClientHello 层与 Chrome 一致。
package browserclient

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

const (
	dialTimeout    = 15 * time.Second
	requestTimeout = 30 * time.Second
	h2IdleTimeout  = 30 * time.Second
)

var errNegotiatedH2 = errors.New("对端协商了 h2，应走 h2 transport")

// ChromeProfile 是 utls 的 Chrome ClientHello profile，随 utls 升级自动跟随。
//
// 用 HelloChrome_Auto 而不是写死版本（如 HelloChrome_133）：utls 升级到新版
// Chrome 指纹时，UA 与 sec-ch-ua 会自动跟着变，永远不会出现
// 「UA 说 Chrome/133、ClientHello 已是 140」这种自相矛盾。
// 版本号只有这一个来源 —— webprofile.UserAgent / SecChUA 都读 ChromeProfile.Version，
// 全仓库没有第二个硬编码的 Chrome 主版本号。
// 由 webprofile.TestUserAgentMatchesTLSProfile 守护 UA 与指纹版本一致。
var ChromeProfile = utls.HelloChrome_Auto

// newChromeUConn 用固定的 Chrome profile 构造 UConn，并把预设里的 ALPN 扩展
// 替换成 nextProtos。
//
// 为什么不能只写 utls.Config.NextProtos：HelloChrome_* 预设走
// applyPresetByID → UTLSIdToSpec（utls v1.8.2 u_parrots.go:2736-2761），
// 扩展列表按预设原样复制，该路径不经过 ALPNExtension.writeToUConn，
// 因此 config.NextProtos 根本不会写进 ClientHello 的 ALPN 扩展 ——
// 传什么都恒为预设的 [h2, http/1.1]。
//
// 正确做法：取出预设 spec，替换其中的 ALPNExtension，再用
// HelloCustom + ApplyPreset 应用。其余扩展的顺序/取值仍与 Chrome 133 一致。
func newChromeUConn(raw net.Conn, host string, nextProtos []string) (*utls.UConn, error) {
	spec, err := utls.UTLSIdToSpec(ChromeProfile)
	if err != nil {
		_ = raw.Close()
		return nil, err
	}
	alpn := &utls.ALPNExtension{AlpnProtocols: append([]string(nil), nextProtos...)}
	replaced := false
	for i, ext := range spec.Extensions {
		if _, ok := ext.(*utls.ALPNExtension); ok {
			spec.Extensions[i] = alpn
			replaced = true
		}
	}
	if !replaced {
		spec.Extensions = append(spec.Extensions, alpn)
	}

	uconn := utls.UClient(raw, &utls.Config{
		ServerName: host,
		// 同时设置 config.NextProtos：握手校验（checkALPN）读的是该字段，
		// 而线报文里的 ALPN 扩展来自上面替换后的 ALPNExtension。
		NextProtos: append([]string(nil), nextProtos...),
	}, utls.HelloCustom)
	// ApplyPreset 会复用 spec 中的扩展对象，这里每次调用都新建 spec/alpn，
	// 不会跨连接共享可变状态。
	if err := uconn.ApplyPreset(&spec); err != nil {
		_ = uconn.Close()
		return nil, err
	}
	return uconn, nil
}

// dialTLS 用 utls 的 Chrome ClientHello 建立 TLS 连接，返回连接与 ALPN 协商结果。
// nextProtos 决定 ALPN 列表（会被计入 JA4 指纹），HTTP 请求用 ["h2","http/1.1"]，
// WebSocket 握手只能跑 HTTP/1.1，用 ["http/1.1"]
func dialTLS(ctx context.Context, network, addr string, nextProtos []string) (net.Conn, string, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	d := net.Dialer{Timeout: dialTimeout, KeepAlive: 30 * time.Second}
	raw, err := d.DialContext(ctx, network, addr)
	if err != nil {
		return nil, "", err
	}
	uconn, err := newChromeUConn(raw, host, nextProtos)
	if err != nil {
		return nil, "", err
	}
	if err := uconn.HandshakeContext(ctx); err != nil {
		_ = uconn.Close()
		return nil, "", err
	}
	return uconn, uconn.ConnectionState().NegotiatedProtocol, nil
}

// DialTLSContext 供 gorilla/websocket 的 NetDialTLSContext 使用。
// 返回的连接已完成 TLS 握手（Chrome 133 指纹）。
//
// ALPN 只给 ["http/1.1"]：B 站弹幕端点 wss://<host>/sub 只做 HTTP/1.1 upgrade，
// gorilla 也只在这条连接上发 h1 请求。若协商成 h2，gorilla 会把服务端的
// h2 SETTINGS 帧当成 HTTP/1.1 响应解析，报 malformed HTTP response 并永远连不上。
// 真实 Chrome 连该端点同样只协商到 h1.1，所以这不是降级。
//
// 注意覆盖方式：必须替换预设里的 ALPN 扩展（见 newChromeUConn）。
// 只在 utls.Config.NextProtos 里写 ["http/1.1"] 是无效的，ClientHello 仍是
// [h2, http/1.1]，这正是之前 WS 一直连不上的原因。
func DialTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, proto, err := dialTLS(ctx, network, addr, []string{"http/1.1"})
	if err != nil {
		return nil, err
	}
	if proto != "" && proto != "http/1.1" {
		_ = conn.Close()
		return nil, fmt.Errorf("ALPN 意外协商为 %q，WebSocket 握手只能走 http/1.1", proto)
	}
	return conn, nil
}

type transport struct {
	h2 *http2.Transport
	h1 *http.Transport

	mu   sync.RWMutex
	noH2 map[string]bool // 已确认不支持 h2 的 host:port，避免重复试探
}

// NewTransport 返回伪装 Chrome TLS 指纹的 RoundTripper：https 请求优先 h2，
// 对端不支持 h2 时自动回退 h1；http 明文请求直接走 h1
func NewTransport() http.RoundTripper {
	t := &transport{noH2: make(map[string]bool)}
	t.h2 = &http2.Transport{
		ReadIdleTimeout: h2IdleTimeout,
		DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			conn, proto, err := dialTLS(ctx, network, addr, []string{"h2", "http/1.1"})
			if err != nil {
				return nil, err
			}
			if proto != "h2" {
				_ = conn.Close()
				t.markNoH2(addr)
				return nil, errNegotiatedH2
			}
			return conn, nil
		},
	}
	t.h1 = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     false,
		DialContext:           (&net.Dialer{Timeout: dialTimeout, KeepAlive: 30 * time.Second}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   dialTimeout,
		ExpectContinueTimeout: 1 * time.Second,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, proto, err := dialTLS(ctx, network, addr, []string{"h2", "http/1.1"})
			if err != nil {
				return nil, err
			}
			if proto == "h2" {
				// h1 transport 不能处理 h2 连接；正常不会走到这里（h2 优先已在上层分流）
				_ = conn.Close()
				return nil, errNegotiatedH2
			}
			return conn, nil
		},
	}
	return t
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.roundTrip(req)
	if err != nil {
		return nil, err
	}
	return decodeBody(resp)
}

func (t *transport) roundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme != "https" {
		return t.h1.RoundTrip(req)
	}
	addr := addrOf(req.URL)
	if !t.isNoH2(addr) {
		resp, err := t.h2.RoundTrip(req)
		if err == nil {
			return resp, nil
		}
		if t.isNoH2(addr) {
			// 首次连接发现对端不支持 h2，已缓存，改走 h1
			return t.h1.RoundTrip(req)
		}
		return nil, err
	}
	return t.h1.RoundTrip(req)
}

// decodeBody 按 Content-Encoding 解压响应体。
//
// 为什么需要它：webprofile 显式发送 Accept-Encoding: gzip, deflate, br, zstd
// （与真实 Chrome 一致，也是任务书 §4.1-C 的要求）。一旦显式设置该头，
// net/http 与 x/net/http2 的透明 gzip 解压就会关闭
// （见 x/net/internal/httpcommon/request.go 的 IsRequestGzip：仅在该头缺失时才返回 true），
// 于是 B 站返回的 br/zstd 压缩体必须自己解，否则 json.Decode 直接报错。
// 实测 api.live.bilibili.com 与 api.bilibili.com 均回 content-encoding: br。
//
// 只处理单层编码（真实服务不会叠加多层）；未压缩或未知编码原样返回。
func decodeBody(resp *http.Response) (*http.Response, error) {
	enc := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding")))
	if enc == "" || enc == "identity" {
		return resp, nil
	}
	var r io.Reader
	switch enc {
	case "gzip", "x-gzip":
		zr, err := gzip.NewReader(resp.Body)
		if err != nil {
			_ = resp.Body.Close()
			return nil, err
		}
		r = zr
	case "deflate":
		// HTTP 的 deflate 理论上是 zlib 封装（2 字节头，CMF/FLG 且 (CMF<<8|FLG)%31==0），
		// 但部分服务端直发裸 deflate。用 Peek 判定而不试读，避免把已消费的字节丢掉。
		br := bufio.NewReader(resp.Body)
		hdr, perr := br.Peek(2)
		var zr io.ReadCloser
		var zerr error
		if perr == nil && len(hdr) == 2 && hdr[0]&0x0f == 8 && (uint16(hdr[0])<<8|uint16(hdr[1]))%31 == 0 {
			zr, zerr = zlib.NewReader(br)
		}
		if zr != nil && zerr == nil {
			r = zr
		} else {
			r = flate.NewReader(br)
		}
	case "br":
		r = brotli.NewReader(resp.Body)
	case "zstd":
		zr, err := zstd.NewReader(resp.Body)
		if err != nil {
			_ = resp.Body.Close()
			return nil, err
		}
		r = zr.IOReadCloser()
	default:
		// 未知编码：交回调用方，读到的仍是原始字节
		return resp, nil
	}
	resp.Body = &decodedBody{Reader: r, closer: resp.Body}
	// 解压后长度已不成立，且内容已被改写，必须去掉这两个头，
	// 否则调用方再按 Content-Encoding 解一次会失败。
	resp.Header.Del("Content-Encoding")
	resp.Header.Del("Content-Length")
	resp.ContentLength = -1
	resp.Uncompressed = true
	return resp, nil
}

// decodedBody 把解压 Reader 与底层的原始 Body 绑在一起，Close 时一并关闭。
type decodedBody struct {
	io.Reader
	closer io.Closer
}

func (d *decodedBody) Close() error {
	var err error
	if c, ok := d.Reader.(io.Closer); ok {
		err = c.Close()
	}
	if e := d.closer.Close(); err == nil {
		err = e
	}
	return err
}

func (t *transport) markNoH2(addr string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.noH2[addr] = true
}

func (t *transport) isNoH2(addr string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.noH2[addr]
}

func addrOf(u *url.URL) string {
	host := u.Host
	if u.Port() == "" {
		if u.Scheme == "http" {
			host = net.JoinHostPort(host, "80")
		} else {
			host = net.JoinHostPort(host, "443")
		}
	}
	return host
}

var (
	sharedOnce sync.Once
	sharedRT   http.RoundTripper
)

// Transport 进程级共享的 Chrome 指纹 RoundTripper，resty 通过 SetTransport 注入，
// 与 Client() 共用连接池
func Transport() http.RoundTripper {
	sharedOnce.Do(func() { sharedRT = NewTransport() })
	return sharedRT
}

// Client 返回使用共享 RoundTripper 的 HTTP 客户端（超时 30s）
func Client() *http.Client {
	return &http.Client{Transport: Transport(), Timeout: requestTimeout}
}
