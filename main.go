package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func main() {
	// 定义一个 []float32 类型的数据
	data := []float32{1.1, 2.2, 3.3, 4.4, 5.5}

	// 使用 bytes.Buffer 创建一个 buffer 来保存编码后的数据
	var buf bytes.Buffer

	// 创建一个 gob 编码器
	encoder := gob.NewEncoder(&buf)

	// 编码 []float32 数据到 buffer
	err := encoder.Encode(data)
	if err != nil {
		fmt.Println("编码错误:", err)
		return
	}

	// 打印编码后的字节流
	fmt.Println("编码后的数据:", buf.Bytes())

	// 创建一个解码器
	decoder := gob.NewDecoder(&buf)

	// 定义一个变量来接收解码后的数据
	var decodedData []float32

	// 解码字节流到 decodedData
	err = decoder.Decode(&decodedData)
	if err != nil {
		fmt.Println("解码错误:", err)
		return
	}

	// 打印解码后的数据
	fmt.Println("解码后的数据:", decodedData)
}
