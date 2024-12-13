package main

import (
	"log"
	"strconv"
)

func main() {

	v := ""
	countFloat, _ := strconv.ParseFloat(v, 64)

	log.Println(countFloat)

}
