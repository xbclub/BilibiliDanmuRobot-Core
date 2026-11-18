package thanksgift

import "testing"

func TestRenderThanks(t *testing.T) {
	got := renderThanks("感谢{user}的{gift}，{num}个共{cost}元", "小明", "粉丝团灯牌", 2, 25000)
	want := "感谢小明的粉丝团灯牌，2个共25.00元"
	if got != want {
		t.Fatalf("renderThanks() = %q, want %q", got, want)
	}
}
