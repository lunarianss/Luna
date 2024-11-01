// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package master

import "github.com/lunarianss/Hurricane/internal/api-server/config"

func Run(config *config.Config) error {
	server, err := createMasterApiServer(config)
	if err != nil {
		return err
	}

	return server.Run()
}
