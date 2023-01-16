package main

import "fmt"

func calSquare() {
	src := make(chan int)
	dst := make(chan int, 3)
	// 子协程A 发送0~9的数字
	go func() {
		defer close(src)
		for i := 0; i < 10; i++ {
			src <- i
		}
	}()
	// 子协程B 获得协程A发送的数字 然后进行平方操作
	go func() {
		defer close(dst)
		for i := range src {
			dst <- i * i
		}
	}()
	// 主协程 获得B协程计算的结果 然后打印输出
	for i := range dst {
		fmt.Println(i)
	}
}

func main() {
	calSquare()
}
