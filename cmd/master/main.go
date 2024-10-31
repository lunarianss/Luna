// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"time"

	master "github.com/lunarianss/Hurricane/internal/apiServer"

	_ "github.com/lunarianss/Hurricane/third_party/forked/automaxprocs"
)

func main() {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	master.NewApp("hurricane").Run()
}
