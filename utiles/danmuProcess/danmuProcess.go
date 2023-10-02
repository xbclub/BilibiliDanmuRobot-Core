package danmuProcess

// 所有弹幕处理类
type DanmuProcessClass struct {
	GptClass *Gpt
}

type DanmuProcess interface {
	// 弹幕处理函数
	DoDanmuProcess()
}
