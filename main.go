// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func main() {
	_, file, _, _ := runtime.Caller(0)

	fmt.Println(file)

	fmt.Println(filepath.Dir(file))

	fmt.Println(filepath.Base(filepath.Dir(file)))

}
