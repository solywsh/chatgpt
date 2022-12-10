# chatgpt

> chartgpt client for golang

## Usege

Download the package first:

```shell
go get github.com/solywsh/chatgpt
```

Conversation without context:

```go
package main

import (
	"fmt"
	chatgpt "github.com/solywsh/chatgpt"
)

func main() {
    // The timeout is used to control the situation that the session is in a long and multi session situation.
    // If it is set to 0, there will be no timeout. Note that a single request still has a timeout setting of 30s.
    chat := New("openai_key", "user_id(not required)", 30*time.Second) 
	defer chat.Close()
    // 
	//select {
	//case <-chat.GetDoneChan():
	//	fmt.Println("time out/finish")
	//}
	question := "你认为2022年世界杯的冠军是谁？"
	fmt.Printf("Q: %s\n", question)
	answer, err := chat.Chat(question)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("A: %s\n", answer)

	//Q: 你认为2022年世界杯的冠军是谁？
	//A: 这个问题很难回答，因为2022年世界杯还没有开始，所以没有人知道冠军是谁。
}
```

Conversation with context:

```golang
package main

import (
	"fmt"
	chatgpt "github.com/solywsh/chatgpt"
)

func main() {
    chat := New("openai_key", "user_id(not required)", 10*time.Second)
	defer chat.Close()
	//select {
	//case <-chat.GetDoneChan():
	//	fmt.Println("time out")
	//}
	question := "现在你是一只猫，接下来你只能用\"喵喵喵\"回答."
	fmt.Printf("Q: %s\n", question)
	answer, err := chat.ChatWithContext(question)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("A: %s\n", answer)
	question = "你是一只猫吗？"
	fmt.Printf("Q: %s\n", question)
	answer, err = chat.ChatWithContext(question)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("A: %s\n", answer)
    
    // Q: 现在你是一只猫，接下来你只能用"喵喵喵"回答.
	// A: 喵喵喵！
    // Q: 你是一只猫吗？
    // A: 喵喵~!
}
```

