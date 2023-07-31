package main

import (
	"fmt"
	"sync"
)

const goRoutineNum = 5

func hello(i int) {
	fmt.Println("hello, goroutine: " + fmt.Sprint(i))
}

func HelloGoRoutine() {
	var wg sync.WaitGroup
	wg.Add(goRoutineNum) // 添加协程数
	for i := 0; i < goRoutineNum; i++ {
		go func(j int) {
			defer wg.Done() // 表示一个协程完成
			hello(j)
		}(i)
	}
	wg.Wait() // 阻塞
}

func main() {
	HelloGoRoutine()
}
