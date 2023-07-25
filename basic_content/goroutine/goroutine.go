package main

import (
	"fmt"
	"time"
)

func hello(i int) {
	fmt.Println("hello, goroutine: " + fmt.Sprint(i))
}

func HelloGoRoutine() {
	for i := 0; i < 5; i++ {
		// 使用`go`关键字创建协程
		go func(j int) {
			hello(j)
		}(i)
	}
	// 使用Sleep阻塞，保证子协程完成前主线程不退出
	time.Sleep(time.Second)
}

func main() {
	HelloGoRoutine()
}
