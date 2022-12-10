package chatgpt

import (
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"strings"
)

var (
	DefaultAiRole    = "AI"
	DefaultHumanRole = "Human"

	DefaultCharacter  = []string{"helpful", "creative", "clever", "friendly", "lovely", "talkative"}
	DefaultBackground = "The following is a conversation with AI assistant. The assistant is %s"
	DefaultPreset     = "\n%s: 你好，让我们开始愉快的谈话！\n%s: 我是 AI assistant ，请问你有什么问题？"
)

type ChatContext struct {
	background  string // 对话背景
	preset      string // 预设对话
	maxSeqTimes int    // 最大对话次数
	aiRole      *role  // AI角色
	humanRole   *role  // 人类角色

	old        []conversation // 旧对话
	restartSeq string         // 重新开始对话的标识
	startSeq   string         // 开始对话的标识

	seqTimes int // 对话次数
}

type conversation struct {
	role   *role
	prompt string
}

type role struct {
	name string
}

func NewContext() *ChatContext {
	return &ChatContext{
		aiRole:      &role{name: DefaultAiRole},
		humanRole:   &role{name: DefaultHumanRole},
		background:  fmt.Sprintf(DefaultBackground, strings.Join(DefaultCharacter, ", ")+"."),
		maxSeqTimes: 10,
		preset:      fmt.Sprintf(DefaultPreset, DefaultHumanRole, DefaultAiRole),
		old:         []conversation{},
		seqTimes:    0,
		restartSeq:  "\n" + DefaultHumanRole + ": ",
		startSeq:    "\n" + DefaultAiRole + ": ",
	}
}

func (c *ChatContext) SetHumanRole(role string) {
	c.humanRole.name = role
	c.restartSeq = "\n" + c.humanRole.name + ": "
}

func (c *ChatContext) SetAiRole(role string) {
	c.aiRole.name = role
	c.startSeq = "\n" + c.aiRole.name + ": "
}

func (c *ChatContext) SetMaxSeqTimes(times int) {
	c.maxSeqTimes = times
}

func (c *ChatContext) GetMaxSeqTimes() int {
	return c.maxSeqTimes
}

func (c *ChatContext) SetBackground(background string) {
	c.background = background
}

func (c *ChatContext) SetPreset(preset string) {
	c.preset = preset
}

func (c *ChatGPT) ChatWithContext(question string) (answer string, err error) {
	question = question + "."
	if len(question) > c.maxQuestionLen {
		return "", OverMaxQuestionLength
	}
	if c.ChatContext.seqTimes >= c.ChatContext.maxSeqTimes {
		return "", OverMaxSequenceTimes
	}
	var promptTable []string
	promptTable = append(promptTable, c.ChatContext.background)
	promptTable = append(promptTable, c.ChatContext.preset)
	for _, v := range c.ChatContext.old {
		if v.role == c.ChatContext.humanRole {
			promptTable = append(promptTable, "\n"+v.role.name+": "+v.prompt)
		} else {
			promptTable = append(promptTable, v.role.name+": "+v.prompt)
		}
	}
	promptTable = append(promptTable, "\n"+c.ChatContext.restartSeq+question)
	prompt := strings.Join(promptTable, "\n")
	prompt += c.ChatContext.startSeq
	if len(prompt) > c.maxText-c.maxAnswerLen {
		return "", OverMaxTextLength
	}
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
		Stop:             []string{c.ChatContext.aiRole.name + ":", c.ChatContext.humanRole.name + ":"},
	}
	resp, err := c.client.CreateCompletion(c.ctx, req)
	if err != nil {
		return "", err
	}
	resp.Choices[0].Text = formatAnswer(resp.Choices[0].Text)
	c.ChatContext.old = append(c.ChatContext.old, conversation{
		role:   c.ChatContext.humanRole,
		prompt: question,
	})
	c.ChatContext.old = append(c.ChatContext.old, conversation{
		role:   c.ChatContext.aiRole,
		prompt: resp.Choices[0].Text,
	})
	c.ChatContext.seqTimes++
	return resp.Choices[0].Text, nil
}
