// Package browserclient 提供统一使用 utls 伪装 Chrome TLS 指纹（JA3/JA4）的 HTTP 传输层，
// 供 resty、net/http 直连与 WebSocket 握手共用，避免暴露 Go 标准库的 TLS 指纹。
//
// 已知不完美（Phase A 范围）：HTTP/2 层仍为 x/net/http2 的默认 SETTINGS 帧
// 与 header 顺序，与 Chrome 存在差异；TLS ClientHello 层与 Chrome 一致。
package browserclient

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

const (
	dialTimeout    = 15 * time.Second
	requestTimeout = 30 * time.Second
	h2IdleTimeout  = 30 * time.Second
)

var errNegotiatedH2 = errors.New("对端协商了 h2，应走 h2 transport")

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
	uconn := utls.UClient(raw, &utls.Config{
		ServerName: host,
		NextProtos: nextProtos,
	}, utls.HelloChrome_Auto)
	if err := uconn.HandshakeContext(ctx); err != nil {
		_ = uconn.Close()
		return nil, "", err
	}
	return uconn, uconn.ConnectionState().NegotiatedProtocol, nil
}

// DialTLSContext 供 gorilla/websocket 的 NetDialTLSContext 使用。
// 返回的连接已完成 TLS 握手（Chrome 指纹），ALPN 仅 http/1.1，
// 因为 gorilla 只会在此连接上发送 HTTP/1.1 upgrade 请求
func DialTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	conn, _, err := dialTLS(ctx, network, addr, []string{"http/1.1"})
	return conn, err
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
