package chatgpt

import (
	gogpt "github.com/sashabaranov/go-gpt3"
	"strings"
)

var (
	DefaultHead      = "\nHuman: 你好，让我们开始愉快的谈话！\nAI: 我是 AI assistant ，请问你有什么问题？"
	DefaultCharacter = []string{"helpful", "creative", "clever", "friendly", "lovely", "talkative"}
	DefaultRole      = "The following is a conversation with Ai assistant. The assistant is" + strings.Join(DefaultCharacter, ",") + "."
)

type ChatContext struct {
	Head            string
	Character       []string
	Role            string
	Old             []conversation
	New             []conversation
	MaxSequence     int
	RestartSequence string
	StartSequence   string
}

type conversation struct {
	role   string
	prompt string
}

func NewContext() *ChatContext {
	return &ChatContext{
		Role:            DefaultRole,
		Head:            DefaultHead,
		Character:       DefaultCharacter,
		Old:             []conversation{},
		New:             []conversation{},
		MaxSequence:     10,
		RestartSequence: "\nHuman: ",
		StartSequence:   "\nAI: ",
	}
}

func (c *ChatGPT) ChatWithContext(question string) (answer string, err error) {
	if c.ChatContext.MaxSequence >= len(c.ChatContext.Old) {
		return "", &OverMaxSequenceTimes{}
	}
	var promptTable []string
	promptTable = append(promptTable, c.ChatContext.Role)
	promptTable = append(promptTable, c.ChatContext.Head)
	for _, v := range c.ChatContext.Old {
		promptTable = append(promptTable, v.role+" "+v.prompt)
	}
	textLen := len(strings.Join(promptTable, ",")) + len(question) + len(c.ChatContext.StartSequence) + len(c.ChatContext.RestartSequence)
	if textLen > c.maxText-c.maxAnswerLen {
		return "", &OverMaxTextLength{}
	}
	promptTable = append(promptTable, c.ChatContext.RestartSequence+question)
	prompt := strings.Join(promptTable, ",")
	req := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		MaxTokens:        c.maxAnswerLen,
		Prompt:           prompt,
		Temperature:      0.9,
		TopP:             1,
		N:                1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.5,
		User:             c.userId,
		Stop:             []string{c.ChatContext.StartSequence, c.ChatContext.RestartSequence},
	}
	resp, err := c.client.CreateCompletion(c.ctx, req)
	if err != nil {
		return "", err
	}
	// Human
	c.ChatContext.Old = append(c.ChatContext.Old, conversation{
		role:   c.ChatContext.RestartSequence,
		prompt: question,
	})
	// AI
	c.ChatContext.Old = append(c.ChatContext.Old, conversation{
		role:   c.ChatContext.StartSequence,
		prompt: resp.Choices[0].Text,
	})
	return resp.Choices[0].Text, nil
}
