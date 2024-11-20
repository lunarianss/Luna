package main

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type B struct {
	Age int `json:"age"`
}

type A struct {
	*B
	Name string `json:"name"`
}

func a() {
	_, fullFilePath, _, _ := runtime.Caller(1)
	fmt.Println(fullFilePath)
}

func main() {
	a := &A{
		B: &B{Age: 13},
	}

	c, _ := json.Marshal(a)

	fmt.Println(string(c))
}
