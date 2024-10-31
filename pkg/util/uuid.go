// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"strings"

	"github.com/google/uuid"
)

func GetUUID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func GetShortUUID(len int) string {
	return GetUUID()[0:len]
}
