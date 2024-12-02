package main

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

func main() {
	text := "Hello, world!"
	text1 := "你好，世界!"
	// var a float64 = 1.2

	// c1, err := strconv.ParseInt(a, 10, 64)
	// c, _ := strconv.ParseFloat(a, 64)
	// if err != nil {
	// 	log.Println(err)
	// }

	const (
		MODEL_O200K_BASE  string = "o200k_base"
		MODEL_CL100K_BASE string = "cl100k_base"
		MODEL_P50K_BASE   string = "p50k_base"
		MODEL_P50K_EDIT   string = "p50k_edit"
		MODEL_R50K_BASE   string = "r50k_base"
	)
	// encoding := "gpt-4o"

	// if you don't want download dictionary at runtime, you can use offline loader
	tke, err := tiktoken.GetEncoding(MODEL_CL100K_BASE)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
	}

	token := tke.Encode(text, nil, nil)
	token1 := tke.Encode(text1, nil, nil)

	//tokens
	fmt.Printf("你好，世界! %d\n", len(token))
	fmt.Printf("text %d\n", len(token1))
}
