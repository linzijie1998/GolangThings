package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	x    int64      // 共享内存
	lock sync.Mutex // 互斥锁
)

func addWithLock() {
	for i := 0; i < 5000; i++ {
		// 临界区保护
		lock.Lock() // 加锁
		x += 1
		lock.Unlock() // 解锁
	}
}

func addWithoutLock() {
	for i := 0; i < 5000; i++ {
		// 临界区无保护
		x += 1
	}
}

func add() {
	x = 0
	for i := 0; i < 5; i++ {
		go addWithoutLock()
	}
	time.Sleep(time.Second)
	fmt.Println("WithoutLock: ", x)
	x = 0
	for i := 0; i < 5; i++ {
		go addWithLock()
	}
	time.Sleep(time.Second)
	fmt.Println("WithLock: ", x)
}

func main() {
	add()
}
