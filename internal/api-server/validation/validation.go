// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	"github.com/lunarianss/Luna/internal/infrastructure/validation"
)

// validation unified registration portal
func init() {
	validation.RegisterValidator(&blogValidation{})
	validation.RegisterValidator(&modelValidation{})
}
