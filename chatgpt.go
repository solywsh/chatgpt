package chatgpt

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	"time"
)

type ChatGPT struct {
	client         *gogpt.Client
	ctx            context.Context
	userId         string
	maxQuestionLen int
	maxText        int
	maxAnswerLen   int
	timeOut        time.Duration // 超时时间, 0表示不超时
	timeOutChan    chan struct {
	}
	cancel         func()
	stopFlag       []string // 如果不设置机器会联想并返回联想结束语句
	maxStopFlagLen int

	ChatContext *ChatContext
}

func New(ApiKey, UserId string, timeOut time.Duration) *ChatGPT {
	var ctx context.Context
	var cancel func()
	if timeOut == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), timeOut)
	}
	timeOutChan := make(chan struct{}, 1)
	go func() {
		<-ctx.Done()
		timeOutChan <- struct{}{} // 发送超时信号，或是提示结束，用于聊天机器人场景，配合GetTimeOutChan() 使用
	}()
	return &ChatGPT{
		client:         gogpt.NewClient(ApiKey),
		ctx:            ctx,
		userId:         UserId,
		maxQuestionLen: 1024, // 最大问题长度
		maxAnswerLen:   1024, // 最大答案长度
		maxText:        4096, // 最大文本 = 问题 + 回答
		timeOut:        timeOut,
		timeOutChan:    timeOutChan,
		cancel: func() {
			cancel()
		},
		stopFlag:       []string{"."},
		maxStopFlagLen: 1,
		ChatContext:    NewContext(),
	}
}
func (c *ChatGPT) Close() {
	c.cancel()
}

func (c *ChatGPT) GetTimeOutChan() chan struct{} {
	return c.timeOutChan
}

func (c *ChatGPT) SetMaxQuestionLen(maxQuestionLen int) {
	if maxQuestionLen > c.maxText-c.maxAnswerLen {
		maxQuestionLen = c.maxText - c.maxAnswerLen
	}
	c.maxQuestionLen = maxQuestionLen
}

func (c *ChatGPT) Chat(question string) (answer string, err error) {
	if len(question)+c.maxAnswerLen+c.maxStopFlagLen > c.maxText {
		question = question[:c.maxText-c.maxAnswerLen]
	}
	req := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		MaxTokens:        c.maxAnswerLen,
		Prompt:           question + ".", // 加"."提示AI结束
		Temperature:      0.9,
		TopP:             1,
		N:                1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.5,
		User:             c.userId,
		Stop:             c.stopFlag,
	}
	resp, err := c.client.CreateCompletion(c.ctx, req)
	if err != nil {
		return "", err
	}
	answer = resp.Choices[0].Text
	for len(answer) > 0 {
		if answer[:1] == "\n" || answer[0] == ' ' {
			answer = answer[1:]
		} else {
			break
		}
	}
	return resp.Choices[0].Text, err
}
