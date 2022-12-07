package chatgpt

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type ChatGPT struct {
	client *gogpt.Client
	ctx    context.Context
	userId string
}

func New(ApiKey, UserId string) *ChatGPT {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return &ChatGPT{
		client: gogpt.NewClient(ApiKey),
		ctx:    ctx,
		userId: UserId,
	}
}
func (c *ChatGPT) Close() {
	c.ctx.Done()
}

func (c *ChatGPT) Chat(question string) (answer string, err error) {
	req := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		MaxTokens:        4096,
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
		if answer[0] == '\n' {
			answer = answer[1:]
		} else {
			break
		}
	}
	return resp.Choices[0].Text, err
}
