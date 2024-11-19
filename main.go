package main

import (
	"fmt"
	"runtime"
)

type A struct {
	Name string `json:"name"`
}

func a() {
	_, fullFilePath, _, _ := runtime.Caller(1)
	fmt.Println(fullFilePath)
}

func main() {
	a()
}
