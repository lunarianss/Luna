// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package field

import (
	"database/sql/driver"

	"github.com/Ryan-eng-del/hurricane/internal/pkg/code"
	"github.com/Ryan-eng-del/hurricane/pkg/errors"
)

type BitBool int

func (b *BitBool) Value() (driver.Value, error) {
	result := make([]uint8, 0, 1)
	if *b == 1 {
		result = append(result, 1)
	} else {
		result = append(result, 0)
	}
	return result, nil
}

func (b *BitBool) Scan(v interface{}) error {
	if bytes, ok := v.([]uint8); ok {
		if bytes[0] == 0 {
			*b = 0
		} else {
			*b = 1
		}
		return nil
	}
	return errors.WithCode(code.ErrScanToField, "can not convert %v to int", v)
}
