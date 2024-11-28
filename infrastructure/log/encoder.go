// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getDefaultEncoder(options *Options) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	if options.EnableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	if options.Layout != "" {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(options.Layout)
	} else {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	}

	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	if options.Format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}
