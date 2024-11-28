// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getDebugCore(opts *Options) zapcore.Core {
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.DebugLevel
	})
	writeSyncDebug := getLogWriter(opts)(opts.DebugFilePath)

	return zapcore.NewCore(getDefaultEncoder(opts), writeSyncDebug, debugPriority)
}

func getInfoCore(opts *Options) zapcore.Core {
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.WarnLevel
	})
	writeSyncDebug := getLogWriter(opts)(opts.InfoFilePath)

	return zapcore.NewCore(getDefaultEncoder(opts), writeSyncDebug, infoPriority)
}

func getWarnCore(opts *Options) zapcore.Core {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	writeSyncDebug := getLogWriter(opts)(opts.ErrorFilePath)

	return zapcore.NewCore(getDefaultEncoder(opts), writeSyncDebug, highPriority)
}

func getLogWriter(opts *Options) func(filename string) zapcore.WriteSyncer {
	return func(filename string) zapcore.WriteSyncer {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
		}

		return zapcore.AddSync(lumberJackLogger)
	}
}
