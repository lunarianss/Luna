package main

import "fmt"

type A struct {
}

func main() {
	var b = make([]A, 0)
	fmt.Println(b == nil, len(b))
}
