// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

type Option interface {
	apply(*Options)
}

type optionFunc func(*Options)

func (f optionFunc) apply(o *Options) {
	f(o)
}

func NewWithOptions(options ...Option) {
	option := NewOption()
	for _, o := range options {
		o.apply(option)
	}
	installOptions(option)
}

func WithEnableColor() Option {
	return optionFunc(func(o *Options) {
		o.EnableColor = true
	})
}

func WithEnableFile() Option {
	return optionFunc(func(o *Options) {
		o.EnableFile = true
	})
}

func WithDebugMode() Option {
	return optionFunc(func(o *Options) {
		o.DebugMode = true
	})
}

func WithLayout(layout string) Option {
	return optionFunc(func(o *Options) {
		o.Layout = layout
	})
}

func WithMaxSize(maxSize int) Option {
	return optionFunc(func(o *Options) {
		o.MaxSize = maxSize
	})
}

func WithMaxAge(maxAge int) Option {
	return optionFunc(func(o *Options) {
		o.MaxAge = maxAge
	})
}

func WithBackups(backups int) Option {
	return optionFunc(func(o *Options) {
		o.MaxBackups = backups
	})
}

func WithFormat(format string) Option {
	return optionFunc(func(o *Options) {
		o.Format = format
	})
}

func WithInfoFilePath(path string) Option {
	return optionFunc(func(o *Options) {
		o.InfoFilePath = path
	})
}

func WithDebugFilePath(path string) Option {
	return optionFunc(func(o *Options) {
		o.DebugFilePath = path
	})
}

func WithErrorFilePath(path string) Option {
	return optionFunc(func(o *Options) {
		o.ErrorFilePath = path
	})
}
