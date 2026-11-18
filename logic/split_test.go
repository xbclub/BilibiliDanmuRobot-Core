package logic

import (
	"strings"
	"testing"
)

func TestSplitMessageKeepsFortyRunes(t *testing.T) {
	msg := strings.Repeat("中", 40)
	got := SplitMessage(msg, 40)
	if len(got) != 1 || got[0] != msg {
		t.Fatalf("40 runes should stay in one message, got %#v", got)
	}

	msg = strings.Repeat("中", 41)
	got = SplitMessage(msg, 40)
	if len(got) != 2 || len([]rune(got[0])) != 40 || len([]rune(got[1])) != 1 {
		t.Fatalf("41 runes should split into 40+1, got %#v", got)
	}
}

func TestSplitMessageCountsWideAndNarrowRunesEqually(t *testing.T) {
	msg := strings.Repeat("a中Ａ１", 10)
	got := SplitMessage(msg, 40)
	if len(got) != 1 || got[0] != msg {
		t.Fatalf("40 mixed-width runes should stay in one message, got %#v", got)
	}
}

func TestSplitMessageEdgeCases(t *testing.T) {
	if got := SplitMessage("", 40); got != nil {
		t.Fatalf("empty message should return nil, got %#v", got)
	}
	if got := SplitMessage("abc", 0); len(got) != 1 || got[0] != "abc" {
		t.Fatalf("non-positive limit should not split, got %#v", got)
	}
	// 整除边界：80/40 恰好两条，无空尾巴
	msg := strings.Repeat("中", 80)
	got := SplitMessage(msg, 40)
	if len(got) != 2 || got[0] != strings.Repeat("中", 40) || got[1] != strings.Repeat("中", 40) {
		t.Fatalf("80 runes should split into exactly two 40-rune parts, got %#v", got)
	}
	// limit=1 逐字符拆分
	got = SplitMessage("abc", 1)
	if len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Fatalf("limit=1 should split per rune, got %#v", got)
	}
}
