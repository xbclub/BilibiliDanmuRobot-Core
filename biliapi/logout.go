package biliapi

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/xbclub/BilibiliDanmuRobot-Core/webprofile"
	"github.com/zeromicro/go-zero/core/logx"
)

// logoutResp B 站注销接口响应（只关心 code）
type logoutResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Logout 退出登录（独立方法，不作为配置项）：
//  1. best-effort 调 B 站注销接口，让 cookie 在服务端失效（失败不阻塞本地清理，
//     5s 超时防止黑洞网络挂住 GUI 调用）
//  2. 清空进程内会话状态（CookieStr / CookieList，写锁内完成）
//  3. 删除登录凭证文件 token/bili_token.txt 与 token/bili_token.json
//
// 本地清理失败（文件删除出错）才返回 error；未登录状态下调用是安全的（幂等）。
// 注意：已建立的弹幕 WebSocket 连接持有旧 cookie 的快照，GUI 注销后应
// StopWsClient 停掉 handler 再重新走扫码登录流程。
func Logout() error {
	// 服务端注销：POST passport.bilibili.com/login/exit/v2，需要 bili_jct
	csrf := CookieValue("bili_jct")
	cookieHeader := CookieHeader()
	serverCalled := false
	if csrf != "" && cookieHeader != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var resp logoutResp
		_, err := cli.R().
			SetContext(ctx).
			SetHeaders(webprofile.MainHeaders()).
			SetHeader("Cookie", cookieHeader).
			SetFormData(map[string]string{"bili_jct": csrf}).
			SetResult(&resp).
			Post("https://passport.bilibili.com/login/exit/v2")
		serverCalled = true
		if err != nil || resp.Code != 0 {
			// 服务端注销失败不影响本地清理（cookie 删除后自然失效）
			logx.Errorf("服务端注销失败（本地凭证仍会清除）: err=%v code=%d msg=%s",
				err, resp.Code, resp.Message)
		} else {
			logx.Info("服务端注销成功")
		}
	}

	// 注销请求期间若发生了新的扫码登录（bili_jct 变化），放弃清空与删文件
	// ——新会话优先；比对与清空在同一写锁块内原子完成，彻底收口 TOCTOU
	skipClear := false
	cookieMu.Lock()
	if serverCalled && CookieList["bili_jct"] != csrf {
		skipClear = true
	} else {
		CookieList = make(map[string]string)
		CookieStr = ""
	}
	cookieMu.Unlock()
	if skipClear {
		logx.Info("检测到注销期间发生新的登录，跳过本地清理")
		return nil
	}

	// 删除本地凭证文件（不存在视为成功，幂等）
	var delErrs []string
	for _, path := range []string{"token/bili_token.txt", "token/bili_token.json"} {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			delErrs = append(delErrs, path+": "+err.Error())
		}
	}
	if len(delErrs) > 0 {
		return errors.New("删除凭证文件失败: " + strings.Join(delErrs, "; "))
	}
	logx.Info("退出登录完成：会话状态已清空，本地凭证已删除")
	return nil
}

// LoggedIn 返回当前进程是否持有登录态（依据 bili_jct/SESSDATA 是否齐备）
func LoggedIn() bool {
	return CookieValue("bili_jct") != "" && CookieValue("SESSDATA") != ""
}

// LogoutLocal 仅执行 Logout 的本地清理部分：清空进程内会话状态（写锁内）
// 并删除凭证文件 token/bili_token.{txt,json}，不发起任何网络请求。
// 适用于「登录态已被服务端判定失效」的场景——会话已死，调注销接口注定失败，
// 直接本地收尾即可。幂等，未登录状态下调用安全；清理失败只记日志不返回错误。
func LogoutLocal() {
	cookieMu.Lock()
	CookieList = make(map[string]string)
	CookieStr = ""
	cookieMu.Unlock()
	for _, path := range []string{"token/bili_token.txt", "token/bili_token.json"} {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			logx.Errorf("删除凭证文件失败 %s: %v", path, err)
		}
	}
	logx.Info("本地登录态已清空（未调用服务端注销接口）")
}
