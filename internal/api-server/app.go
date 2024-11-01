// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package master

import (
	"github.com/lunarianss/Hurricane/internal/api-server/config"
	"github.com/lunarianss/Hurricane/internal/api-server/options"
	"github.com/lunarianss/Hurricane/pkg/app"
	"github.com/lunarianss/Hurricane/pkg/log"
)

// nolint: lll
const commandDesc = `Hurricane, a command and rich functional web develop template.

Find more hurricane information at:
    https://github.com/lunarianss/hurricane/`

// NewApp creates an App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("Hurricane Distributed CronTab Application",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.New(opts.Log)
		defer log.Sync()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
