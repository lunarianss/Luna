package main

import "fmt"

type A struct {
	Name string `json:"name"`
}

func main() {

	var b *A
	var c A
	// a := A{
	// 	Name: "cyan",
	// }
	fmt.Println(b == nil, &c == nil)

}
