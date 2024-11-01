// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	defaultFormat      = "text"
	defaultEnableColor = false
	defaultEnableFile  = false
	consoleFormat      = "text"
	jsonFormat         = "json"
)

var (
	defaultDebugFilePath = getDefaultFilePath("hurricane.deb.log")
	defaultInfoFilePath  = getDefaultFilePath("hurricane.inf.log")
	defaultErrorFilePath = getDefaultFilePath("hurricane.err.log")
)

func getDefaultFilePath(name string) string {
	pwd, _ := os.Getwd()

	return filepath.Join(pwd, name)
}

type Options struct {
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `mapstructure:"max-size"        json:"max-size"`
	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `mapstructure:"max-age"         json:"max-age"`
	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups    int    `mapstructure:"max-backups"     json:"max-backups"`
	DebugFilePath string `mapstructure:"debug-file-path" json:"debug-file-path"`
	Format        string `mapstructure:"format"          json:"format"`
	Layout        string `mapstructure:"layout"          json:"layout"`
	EnableFile    bool   `mapstructure:"enable-file"     json:"enable-file"`
	InfoFilePath  string `mapstructure:"info-file-path"  json:"info-file-path"`
	ErrorFilePath string `mapstructure:"error-file-path" json:"error-file-path"`
	EnableColor   bool   `mapstructure:"enable-color"    json:"enable-color"`
	DebugMode     bool   `mapstructure:"debug-mode"      json:"debug-mode"`
}

// Validate validate the options fields.
func (o *Options) Validate() []error {
	var errs []error

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.DebugFilePath, "log.debug-file-path", o.DebugFilePath, "debug out file path")
	fs.StringVar(&o.InfoFilePath, "log.info-file-path", o.InfoFilePath, "info out file path")
	fs.StringVar(&o.ErrorFilePath, "log.error-file-path", o.ErrorFilePath, "error out file path")

	fs.BoolVar(&o.EnableFile, "log.enable-file", o.EnableFile, "Disable output of caller information in the log.")
	fs.BoolVar(&o.EnableColor, "log.enable-color", o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.BoolVar(&o.DebugMode, "log.debug-mode", o.DebugMode, "Enable output ansi colors in plain format logs.")

	fs.IntVar(&o.MaxAge, "log.max-age", o.MaxAge, "Max age of the log")
	fs.IntVar(&o.MaxBackups, "log.max-backups", o.MaxAge, "Max backups of the log")
	fs.IntVar(&o.MaxSize, "log.max-size", o.MaxSize, "Max sizes of the log")

	fs.StringVar(&o.Format, "log.format",
		o.Format, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Layout, "log.layout",
		o.Layout, "Disable the log to record a stack trace for all messages at or above panic level.")
}

// nolint: errchkjson
func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func NewOption() *Options {
	return &Options{
		Format:      defaultFormat,
		EnableColor: defaultEnableColor,
		EnableFile:  defaultEnableFile,
	}
}

func getDefaultFileOptions(opts *Options) *Options {
	if opts.DebugMode && opts.DebugFilePath == "" {
		opts.DebugFilePath = defaultDebugFilePath
	}

	if opts.ErrorFilePath == "" {
		opts.ErrorFilePath = defaultErrorFilePath
	}

	if opts.InfoFilePath == "" {
		opts.InfoFilePath = defaultInfoFilePath
	}

	return opts
}

func setDisableFileLogger(options *Options) []zapcore.Core {
	var (
		zapCores = make([]zapcore.Core, 0, 1)
		minLevel = zapcore.DebugLevel
	)

	if !options.DebugMode {
		minLevel = zapcore.InfoLevel
	}

	stdCore := zapcore.NewCore(getDefaultEncoder(options), zapcore.Lock(os.Stdout), minLevel)

	return append(zapCores, stdCore)
}

func setFileLogger(options *Options) []zapcore.Core {
	zapCores := make([]zapcore.Core, 0, 4)
	options = getDefaultFileOptions(options)

	if options.DebugMode {
		zapCores = append(zapCores, getDebugCore(options))
	}

	return append(zapCores, getInfoCore(options), getWarnCore(options))
}
