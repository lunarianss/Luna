// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"github.com/lunarianss/Luna/internal/api-server/options"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/infrastructure/errors"
)

// Config is the running configuration structure of the Luna service.
type Config struct {
	*options.Options
}

var lunaRuntimeConfiguration *Config

// CreateConfigFromOptions creates a running configuration instance based
// on a given Luna pump command line or configuration file option.
func CreateConfigFromOptions(options *options.Options) (*Config, error) {
	lunaRuntimeConfiguration = &Config{options}
	return lunaRuntimeConfiguration, nil
}

func GetLunaRuntimeConfig() (*Config, error) {
	if lunaRuntimeConfiguration == nil {
		return nil, errors.WithCode(code.ErrRunTimeConfig, "luna runtime configuration is nil")
	}
	return lunaRuntimeConfiguration, nil
}
