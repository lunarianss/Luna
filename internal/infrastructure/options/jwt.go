// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"time"

	"github.com/spf13/pflag"
)

// JwtOptions contains configuration items related to API server features.
type JwtOptions struct {
	Realm   string        `json:"realm"       mapstructure:"realm"`
	Key     string        `json:"key"         mapstructure:"key"`
	Timeout time.Duration `json:"timeout"     mapstructure:"timeout"`
	Refresh time.Duration `json:"max-refresh" mapstructure:"refresh"`
}

// NewJwtOptions creates a JwtOptions object with default parameters.
func NewJwtOptions() *JwtOptions {
	return &JwtOptions{
		Realm:   "luna jwt",
		Timeout: 24 * time.Hour,
		Refresh: 30 * 24 * time.Hour,
	}
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (s *JwtOptions) Validate() []error {
	var errs []error

	return errs
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (s *JwtOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.StringVar(&s.Realm, "jwt.realm", s.Realm, "Realm name to display to the user.")
	fs.StringVar(&s.Key, "jwt.key", s.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&s.Timeout, "jwt.timeout", s.Timeout, "JWT token timeout.")

	fs.DurationVar(&s.Refresh, "jwt.refresh", s.Refresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")
}
