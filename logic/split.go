package logic

// SplitDanmu 按弹幕长度限制（rune 计数）把长文本切分为多条弹幕，
// danmuLen 小于 1 时按 1 处理，避免除零
func SplitDanmu(msg string, danmuLen int) []string {
	if danmuLen < 1 {
		danmuLen = 1
	}
	var msgdata []string
	msgrun := []rune(msg)
	msgLen := len(msgrun)
	msgcount := msgLen / danmuLen
	if msgLen%danmuLen != 0 {
		msgcount += 1
	}
	for m := 1; m <= msgcount; m++ {
		if msgLen < m*danmuLen {
			msgdata = append(msgdata, string(msgrun[(m-1)*danmuLen:msgLen]))
			continue
		}
		msgdata = append(msgdata, string(msgrun[(m-1)*danmuLen:danmuLen*m]))
	}
	return msgdata
}
