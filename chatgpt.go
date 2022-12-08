package chatgpt

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	"time"
)

type ChatGPT struct {
	client   *gogpt.Client
	ctx      context.Context
	userId   string
	maxToken int
	timeOut  time.Duration
}

func New(ApiKey, UserId string, timeOut time.Duration) *ChatGPT {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return &ChatGPT{
		client:   gogpt.NewClient(ApiKey),
		ctx:      ctx,
		userId:   UserId,
		maxToken: 1024,
	}
}
func (c *ChatGPT) Close() {
	c.ctx.Done()
}

func (c *ChatGPT) SetMaxToken(maxToken int) {
	if maxToken > 4096 {
		maxToken = 4096
		return
	}
	c.maxToken = maxToken
}

func (c *ChatGPT) Chat(question string) (answer string, err error) {
	if len(question)+c.maxToken > 4096 {
		question = question[:4096-c.maxToken]
	}
	req := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		MaxTokens:        c.maxToken,
		Prompt:           question,
		Temperature:      0.9,
		TopP:             1,
		N:                1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.5,
		User:             c.userId,
	}
	resp, err := c.client.CreateCompletion(c.ctx, req)
	if err != nil {
		return "", err
	}
	answer = resp.Choices[0].Text
	for len(answer) > 0 {
		if answer[0] == '\n' || answer[0] == ' ' {
			answer = answer[1:]
		} else {
			break
		}
	}
	return resp.Choices[0].Text, err
}
