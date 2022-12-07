package chatgpt

import (
	"context"
	"fmt"
	gpt3 "github.com/solywsh/go-gpt3"
)

type ChatGPT struct {
	client gpt3.Client
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
		client: gpt3.NewClient(ApiKey),
		ctx:    ctx,
		userId: UserId,
	}
}
func (c *ChatGPT) Close() {
	c.ctx.Done()
}

func (c *ChatGPT) Chat(question string) (answer string, err error) {
	resp, err := c.client.CompletionWithEngine(c.ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens: gpt3.IntPtr(30),
		User:      c.userId,
	})
	if err != nil {
		return "", err
	}
	fmt.Printf("resp: %+v", resp)
	answer = resp.Choices[0].Text
	for {
		if answer[0] == '\n' {
			answer = answer[1:]
		} else {
			break
		}
	}
	return answer, err
}
