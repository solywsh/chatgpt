package chatgpt

import (
	"context"
	gpt3 "github.com/PullRequestInc/go-gpt3"
	"strings"
)

type ChatGPT struct {
	client gpt3.Client
	ctx    context.Context
}

func New(ApiKey string) *ChatGPT {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return &ChatGPT{
		client: gpt3.NewClient(ApiKey),
		ctx:    ctx,
	}
}
func (c *ChatGPT) Close() {
	c.ctx.Done()
}

func (c *ChatGPT) Chat(question string) (answer string, err error) {
	var ans strings.Builder
	err = c.client.CompletionStreamWithEngine(c.ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(response *gpt3.CompletionResponse) {
		ans.WriteString(response.Choices[0].Text)
	})
	if err != nil {
		return answer, err
	}
	answer = ans.String()
	for {
		if answer[0] == '\n' {
			answer = answer[1:]
		} else {
			break
		}
	}
	return answer, err
}
