package main

import (
	"fmt"
	"sync"
)

func a() {

	// 模拟 panic
	fmt.Println("Goroutine is going to panic")
	panic("Something went wrong in goroutine")
}

func safeGoroutine(wg *sync.WaitGroup) {
	defer wg.Done()

	// 使用 defer 和 recover 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	a()

}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go safeGoroutine(&wg)

	// 等待 goroutine 完成
	wg.Wait()

	fmt.Println("Main program continues execution")
}
