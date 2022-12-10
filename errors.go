package chatgpt

import "errors"

// OverMaxSequenceTimes 超过最大对话时间
var OverMaxSequenceTimes = errors.New("maximum conversation times exceeded")

// OverMaxTextLength 超过最大文本长度
var OverMaxTextLength = errors.New("maximum text length exceeded")
