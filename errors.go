package chatgpt

// OverMaxSequenceTimes 超过最大对话时间
type OverMaxSequenceTimes struct {
}

func (e *OverMaxSequenceTimes) Error() string {
	return "maximum conversation times exceeded"
}

// OverMaxTextLength 超过最大文本长度
type OverMaxTextLength struct {
}

func (e *OverMaxTextLength) Error() string {
	return "maximum text length exceeded"
}
