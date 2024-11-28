// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package field

import (
	"fmt"
	"strconv"
)

type Float64 float64

// UnmarshalYAML allows the custom unmarshalling of the float64 value.
func (f *Float64) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	// Try to parse the string as a float64
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("cannot convert %s to float64: %v", s, err)
	}
	*f = Float64(val)
	return nil
}
