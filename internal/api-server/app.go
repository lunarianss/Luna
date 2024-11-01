// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package master

import (
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/options"
	"github.com/lunarianss/Luna/pkg/app"
	"github.com/lunarianss/Luna/pkg/log"
)

// nolint: lll
const commandDesc = `Luna🌛 is an open-source platform for building AI applications ⚡️, combine LLMOps to streamline the development of generative AI solutions.
Find more Luna information at:
    https://github.com/lunarianss/Luna/`

// NewApp creates an App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("Luna AI Application",
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
