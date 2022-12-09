package chatgpt

import "strings"

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

func (c *ChatGPT) ChatWithContext(question string) {
	//promptTable := strings.Builder{}
	//promptTable.WriteString(c.ChatContext.Role)
	//promptTable.WriteString("\n")
	//promptTable.WriteString(c.ChatContext.Head)
	//promptTable.WriteString("\n")
	// 性能去他妈
	var promptTable []string
	promptTable = append(promptTable, c.ChatContext.Role)
	promptTable = append(promptTable, c.ChatContext.Head)
	for _, v := range c.ChatContext.Old {
		promptTable = append(promptTable, v.role+" "+v.prompt)
	}
	promptTable = append(promptTable, c.ChatContext.RestartSequence+question+".")

}

//func (c *ChatGPT) cutPrompt(prompt []string) []string {
//	//extra := len(prompt) + c.ChatContext.MaxSequence + ""
//}
