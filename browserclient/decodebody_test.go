package browserclient

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

// wantBody 是所有压缩用例解压后应还原的原始响应体。
const wantBody = `{"code":0,"data":{"room_id":23058,"title":"离线解压测试"},"message":"0"}`

// 压缩数据全部在测试内现场构造，不依赖网络。

func gzipEncode(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		t.Fatalf("gzip 写入失败: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("gzip Close 失败: %v", err)
	}
	return buf.Bytes()
}

func zlibEncode(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		t.Fatalf("zlib 写入失败: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("zlib Close 失败: %v", err)
	}
	return buf.Bytes()
}

func brotliEncode(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := brotli.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		t.Fatalf("brotli 写入失败: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("brotli Close 失败: %v", err)
	}
	return buf.Bytes()
}

func zstdEncode(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w, werr := zstd.NewWriter(&buf)
	if werr != nil {
		t.Fatalf("zstd NewWriter 失败: %v", werr)
	}
	if _, err := w.Write(data); err != nil {
		t.Fatalf("zstd 写入失败: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("zstd Close 失败: %v", err)
	}
	return buf.Bytes()
}

// newEncodedResponse 手工构造带压缩体的 http.Response，
// 不经过真实网络，直接喂给 decodeBody。
func newEncodedResponse(enc string, body []byte) *http.Response {
	return &http.Response{
		StatusCode:    200,
		Header:        http.Header{"Content-Encoding": {enc}, "Content-Length": {strconv.Itoa(len(body))}},
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
	}
}

// TestDecodeBodyDecompresses 覆盖四种压缩编码（含 x-gzip 别名与大小写/空白归一化）：
// 断言解压还原原文、Content-Encoding/Content-Length 头被删除、ContentLength=-1、Uncompressed=true。
func TestDecodeBodyDecompresses(t *testing.T) {
	raw := []byte(wantBody)
	cases := []struct {
		name string
		enc  string // 写入响应头的原始值
		body []byte
	}{
		{name: "gzip", enc: "gzip", body: gzipEncode(t, raw)},
		{name: "x-gzip 别名", enc: "x-gzip", body: gzipEncode(t, raw)},
		{name: "deflate 按 zlib 封装", enc: "deflate", body: zlibEncode(t, raw)},
		{name: "br", enc: "br", body: brotliEncode(t, raw)},
		{name: "zstd", enc: "zstd", body: zstdEncode(t, raw)},
		{name: "编码值带空白与大写仍识别为 gzip", enc: "  GZIP ", body: gzipEncode(t, raw)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeBody(newEncodedResponse(tc.enc, tc.body))
			if err != nil {
				t.Fatalf("decodeBody 返回错误: %v", err)
			}
			data, rerr := io.ReadAll(got.Body)
			if rerr != nil {
				t.Fatalf("读取解压后响应体失败: %v", rerr)
			}
			if string(data) != wantBody {
				t.Errorf("解压结果不还原\n got: %q\nwant: %q", data, wantBody)
			}
			if cerr := got.Body.Close(); cerr != nil {
				t.Errorf("Close 解压后响应体失败: %v", cerr)
			}
			if v := got.Header.Get("Content-Encoding"); v != "" {
				t.Errorf("Content-Encoding 头未删除，仍为 %q", v)
			}
			if v := got.Header.Get("Content-Length"); v != "" {
				t.Errorf("Content-Length 头未删除，仍为 %q", v)
			}
			if got.ContentLength != -1 {
				t.Errorf("ContentLength = %d，解压后长度已不成立，期望 -1", got.ContentLength)
			}
			if !got.Uncompressed {
				t.Error("Uncompressed 应为 true")
			}
		})
	}
}

// TestDecodeBodyUnknownEncoding 未知编码（如 "foo"）不解压：
// 原样返回原始字节，且保留 Content-Encoding 头交回调用方处理。
func TestDecodeBodyUnknownEncoding(t *testing.T) {
	raw := []byte("these are raw bytes of an unknown encoding")
	got, err := decodeBody(newEncodedResponse("foo", raw))
	if err != nil {
		t.Fatalf("decodeBody 返回错误: %v", err)
	}
	data, rerr := io.ReadAll(got.Body)
	if rerr != nil {
		t.Fatalf("读取响应体失败: %v", rerr)
	}
	if !bytes.Equal(data, raw) {
		t.Errorf("未知编码应原样返回\n got: %q\nwant: %q", data, raw)
	}
	if cerr := got.Body.Close(); cerr != nil {
		t.Errorf("Close 响应体失败: %v", cerr)
	}
	if v := got.Header.Get("Content-Encoding"); v != "foo" {
		t.Errorf("未知编码应保留 Content-Encoding 头，got %q", v)
	}
}

// TestDecodeBodyPassthrough 无 Content-Encoding（或 identity）时直接透传：
// 同一个 Response 指针原样返回，体与其他头均不动。
func TestDecodeBodyPassthrough(t *testing.T) {
	cases := []struct {
		name string
		hdr  http.Header
	}{
		{name: "无 Content-Encoding 头", hdr: http.Header{"Content-Type": {"application/json"}}},
		{name: "identity 编码", hdr: http.Header{"Content-Type": {"application/json"}, "Content-Encoding": {"identity"}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := []byte(wantBody)
			resp := &http.Response{
				StatusCode: 200,
				Header:     tc.hdr.Clone(),
				Body:       io.NopCloser(bytes.NewReader(raw)),
			}
			got, err := decodeBody(resp)
			if err != nil {
				t.Fatalf("decodeBody 返回错误: %v", err)
			}
			if got != resp {
				t.Fatal("应原样返回同一个 *http.Response")
			}
			data, rerr := io.ReadAll(got.Body)
			if rerr != nil {
				t.Fatalf("读取响应体失败: %v", rerr)
			}
			if string(data) != wantBody {
				t.Errorf("透传体被改动\n got: %q\nwant: %q", data, wantBody)
			}
			if cerr := got.Body.Close(); cerr != nil {
				t.Errorf("Close 响应体失败: %v", cerr)
			}
			if got.Uncompressed {
				t.Error("未做解压，Uncompressed 不应被置位")
			}
		})
	}
}

// TestDecodeBodyEmptyBody 是空 body 边界的回归测试：
// 响应体为空但带 Content-Encoding 时（修复前 gzip.NewReader 对空流返回 io.EOF，
// 被当作解压错误），应跳过解压直接返回空体，且不报错、头部原样保留。
func TestDecodeBodyEmptyBody(t *testing.T) {
	for _, enc := range []string{"gzip", "deflate", "br", "zstd"} {
		t.Run("空体+"+enc, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: 200,
				Header:     http.Header{"Content-Encoding": {enc}},
				Body:       io.NopCloser(bytes.NewReader(nil)),
			}
			got, err := decodeBody(resp)
			if err != nil {
				t.Fatalf("空 body 不应报错（修复前 %s 会把 io.EOF 当解压失败）: %v", enc, err)
			}
			data, rerr := io.ReadAll(got.Body)
			if rerr != nil {
				t.Fatalf("读取空响应体失败: %v", rerr)
			}
			if len(data) != 0 {
				t.Errorf("空体被改写为 %q", data)
			}
			if cerr := got.Body.Close(); cerr != nil {
				t.Errorf("Close 空响应体失败: %v", cerr)
			}
			if v := got.Header.Get("Content-Encoding"); v != enc {
				t.Errorf("跳过解压应保留 Content-Encoding 头，got %q", v)
			}
		})
	}
}
