package biliapi

import (
	"errors"
	"os"
	"strings"

	"github.com/xbclub/BilibiliDanmuRobot-Core/webprofile"
	"github.com/zeromicro/go-zero/core/logx"
)

// logoutResp B 站注销接口响应（只关心 code）
type logoutResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Logout 退出登录（独立方法，不作为配置项）：
//  1. best-effort 调 B 站注销接口，让 cookie 在服务端失效（失败不阻塞本地清理）
//  2. 清空进程内会话状态（CookieStr / CookieList）
//  3. 删除登录凭证文件 token/bili_token.txt 与 token/bili_token.json
//
// 本地清理失败（文件删除出错）才返回 error；未登录状态下调用是安全的（幂等）。
// 注意：已建立的弹幕 WebSocket 连接持有旧 cookie 的快照，GUI 注销后应
// StopWsClient 停掉 handler 再重新走扫码登录流程。
func Logout() error {
	// 服务端注销：POST passport.bilibili.com/login/exit/v2，需要 bili_jct
	if csrf, ok := CookieList["bili_jct"]; ok && csrf != "" && CookieStr != "" {
		var resp logoutResp
		_, err := cli.R().
			SetHeaders(webprofile.MainHeaders()).
			SetHeader("Cookie", CookieStr).
			SetFormData(map[string]string{"bili_jct": csrf}).
			SetResult(&resp).
			Post("https://passport.bilibili.com/login/exit/v2")
		if err != nil || resp.Code != 0 {
			// 服务端注销失败不影响本地清理（cookie 删除后自然失效）
			logx.Errorf("服务端注销失败（本地凭证仍会清除）: err=%v code=%d msg=%s",
				err, resp.Code, resp.Message)
		} else {
			logx.Info("服务端注销成功")
		}
	}

	// 清空进程内会话状态
	CookieList = make(map[string]string)
	CookieStr = ""

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
	return CookieList["bili_jct"] != "" && CookieList["SESSDATA"] != ""
}
