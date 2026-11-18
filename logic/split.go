package logic

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
