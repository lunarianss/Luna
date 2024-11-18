package main

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	var duration time.Duration

	// 绑定命令行标志 --duration 到变量 duration
	pflag.DurationVar(&duration, "duration", 0, "Specify a duration (e.g., 1h, 30m, 2d)")
	pflag.Parse()

	// 打印解析后的值
	fmt.Printf("Parsed duration: %v\n", duration)

	// 示例：获取总小时数
	fmt.Printf("Total hours: %.2f\n", duration.Hours())
}
