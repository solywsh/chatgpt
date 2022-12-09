# chatgpt

> chartgpt client for golang

## Usege

```go
package main

import (
	"fmt"
	chatgpt "github.com/solywsh/chatgpt"
)

func main() {
	chat := New("", "", 10*time.Second)
	defer chat.Close()
	select {
	case <-chat.GetTimeOutChan():
		fmt.Println("time out")
	}
	question := "中国在欧洲\n"
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

