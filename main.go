package main

import (
	"fmt"
	"log"
)

func writeMap(m map[string]interface{}) {
	m["1"] = 2
}

func writeSlice(s *[]string) {
	*s = append(*s, "13")
	log.Println(s)
}

func get() (string, error) {
	return "1", nil
}

func aaa() {
	fmt.Println("aaa")
}
func generateGoRoutine() {

	defer func() {
		fmt.Println("msg ---- 111")
	}()

	aaa()
}

func main() {

	generateGoRoutine()

}
