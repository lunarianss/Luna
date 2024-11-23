package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type B struct {
	Age int `json:"age"`
}

type A struct {
	*B
	Name string `json:"name"`
	Ps   []*B   `json:"ps"`
}

func a() {
	_, fullFilePath, _, _ := runtime.Caller(1)
	fmt.Println(fullFilePath)
}

func main() {
	a := &A{}

	fmt.Println(a.Ps == nil)

	a.Ps = nil

	for _, p := range a.Ps {
		fmt.Println(p)
	}

	log.Printf("%#+v", errors.WithCode(code.ErrAccountBanned, ""))
}
