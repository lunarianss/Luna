// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"context"
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log    *zap.SugaredLogger
	StdLog *zap.SugaredLogger
)

var _logOnce sync.Once

// Defines common log fields.
const (
	KeyRequestID string = "requestID"
	KeyUsername  string = "username"
)

// StdInfoLogger returns logger of standard library which writes to supplied zap
// logger at info level.
func StdInfoLogger() *log.Logger {
	if StdLog == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(StdLog.Desugar(), zapcore.InfoLevel); err == nil {
		return l
	}
	return nil
}

func setGlobalZapLogger(zapCores []zapcore.Core) {
	_logOnce.Do(func() {
		if Log == nil {
			zapCore := zapcore.NewTee(zapCores...)
			Logger := zap.New(zapCore, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.PanicLevel))
			zap.ReplaceGlobals(Logger)
			Log = Logger.Sugar()
		}
	})
}

func New(opt *Options) {
	installOptions(opt)
}

func installOptions(option *Options) {
	var (
		zapCores   = make([]zapcore.Core, 0, 1)
		enableFile = option.EnableFile
	)

	if enableFile {
		zapCores = setFileLogger(option)
	}

	zapCores = append(zapCores, setDisableFileLogger(option)...)
	setGlobalZapLogger(zapCores)
}

func Sync() {
	if err := Log.Sync(); err != nil {
		Log.Errorf("LogSyncError: %s", err.Error())
	}
}

func getLogger() *zap.SugaredLogger {
	if Log != nil {
		return Log
	}

	return StdLog
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	getLogger().Debug(args)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	getLogger().Info(args)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	getLogger().Warn(args)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	getLogger().Error(args)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) {
	getLogger().DPanic(args)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	getLogger().Panic(args)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	getLogger().Fatal(args)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	getLogger().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	getLogger().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	getLogger().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	getLogger().Errorf(template, args...)
}

func With(args ...interface{}) *zap.SugaredLogger {
	return getLogger().With(args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	getLogger().DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	getLogger().Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	getLogger().Fatalf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//
//	s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) {
	getLogger().Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	getLogger().Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	getLogger().Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	getLogger().Errorw(msg, keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) {
	getLogger().DPanicw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	getLogger().Panicw(msg, keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	getLogger().Fatalw(msg, keysAndValues...)
}

// Info uses fmt.Sprint to construct and log a message.
func InfoL(ctx context.Context, args ...interface{}) {
	var withCtxLogger *zap.SugaredLogger = nil
	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		withCtxLogger = getLogger().With(zap.Any(KeyRequestID, requestID))
	}
	if username := ctx.Value(KeyUsername); username != nil {
		if withCtxLogger == nil {
			withCtxLogger = getLogger().With(zap.Any(KeyUsername, username))
		} else {
			withCtxLogger = withCtxLogger.With(zap.Any(KeyUsername, username))
		}
	}

	if withCtxLogger == nil {
		getLogger().With().Info(args)
	} else {
		withCtxLogger.Info(args)
	}
}
