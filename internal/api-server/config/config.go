// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import "github.com/lunarianss/Hurricane/internal/api-server/options"

// Config is the running configuration structure of the hurricane service.
type Config struct {
	*options.Options
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given hurricane pump command line or configuration file option.
func CreateConfigFromOptions(options *options.Options) (*Config, error) {
	return &Config{options}, nil
}
