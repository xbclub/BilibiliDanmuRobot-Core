package logic

// SplitMessage 按弹幕长度限制（rune 计数）把长文本切分为多条弹幕，
// limit 小于等于 0 时不切分，原文返回
func SplitMessage(msg string, limit int) []string {
	if msg == "" {
		return nil
	}
	if limit <= 0 {
		return []string{msg}
	}
	runes := []rune(msg)
	if len(runes) <= limit {
		return []string{msg}
	}
	var out []string
	for len(runes) > limit {
		out = append(out, string(runes[:limit]))
		runes = runes[limit:]
	}
	if len(runes) > 0 {
		out = append(out, string(runes))
	}
	return out
}
