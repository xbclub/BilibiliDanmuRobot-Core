package logic

import (
	"reflect"
	"strings"
	"testing"
)

func TestSplitDanmu(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		danmuLen int
		want     []string
	}{
		{"空文本", "", 20, nil},
		{"短文本不切分", "你好", 20, []string{"你好"}},
		{"恰好等长", "一二三四五六七八九十", 10, []string{"一二三四五六七八九十"}},
		{"超长中文切分", strings.Repeat("弹", 45), 20, []string{strings.Repeat("弹", 20), strings.Repeat("弹", 20), strings.Repeat("弹", 5)}},
		{"非整数倍余一条", "一二三四五", 2, []string{"一二", "三四", "五"}},
		{"emoji 按 rune 切分", "😀😅🙂🙃😉😊😐🤐🤨😐", 3, []string{"😀😅🙂", "🙃😉😊", "😐🤐🤨", "😐"}},
		{"danmuLen 为 0 防除零", "一二三", 0, []string{"一", "二", "三"}},
		{"danmuLen 为负防除零", "一二三", -5, []string{"一", "二", "三"}},
		{"danmuLen 为 1 逐字切分", "abc", 1, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitDanmu(tt.msg, tt.danmuLen); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitDanmu(%q, %d) = %q, want %q", tt.msg, tt.danmuLen, got, tt.want)
			}
		})
	}
}

func TestSplitDanmuNoChunkExceedsLimit(t *testing.T) {
	msg := "混合内容：Hello世界！This is a long message 用于验证每段都不超限。"
	for _, limit := range []int{1, 2, 5, 20} {
		for _, chunk := range SplitDanmu(msg, limit) {
			if got := len([]rune(chunk)); got > limit {
				t.Errorf("limit=%d 时切出长度 %d 的分段: %q", limit, got, chunk)
			}
		}
	}
}
