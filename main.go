package main

import (
	"fmt"
	"time"
)

func producer(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch) // 关闭通道，通知消费者数据已发送完毕
	fmt.Println("关闭了通道")
}

func consumer(ch chan int, done chan struct{}) {
	for {
		select {
		case value, ok := <-ch:
			time.Sleep(2 * time.Second)
			if !ok {
				// 通道已关闭且无数据，退出消费者循环5
				fmt.Println("Channel closed, consumer exiting")
				done <- struct{}{} // 通知主协程消费者已完成
				return
			}
			time.Sleep(2 * time.Second)
			fmt.Println("Received:", value)
		}
	}
}

func main() {
	ch := make(chan int)
	done := make(chan struct{})

	// 启动生产者
	go producer(ch)

	// 启动消费者
	go consumer(ch, done)

	// 等待消费者处理完所有数据
	<-done
	fmt.Println("All data processed, main exiting")
}
