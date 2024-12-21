package main

import (
	"fmt"
	"log"
)

func main() {
	err := fmt.Errorf("1233")

	log.Println(err.Error())
}
