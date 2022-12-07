# chatgpt

> chartgpt client for golang

## Usege

```go
package main

import (
	"fmt"
	chatgpt "github.com/solywsh/chartgpt"
)

func main() {
	chat := chatgpt.New("YOUR_API_KEY")
	defer chat.Close()
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

