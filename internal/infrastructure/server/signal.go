// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"os"
	"os/signal"
	"sync"
)

var (
	setupOnceHandler sync.Once
	stop             chan struct{}
	shoutDownHandler chan os.Signal
)

func SetupSignalHandler() <-chan struct{} {
	setupOnceHandler.Do(func() {
		shoutDownHandler = make(chan os.Signal, 2)
		signal.Notify(shoutDownHandler, shutdownSignals...)

		go func() {
			<-shoutDownHandler
			close(stop)
			<-shoutDownHandler
			os.Exit(1)
		}()
	})

	return stop
}

func RequestShutdown() bool {
	if shoutDownHandler != nil {
		select {
		case shoutDownHandler <- shutdownSignals[0]:
			return true
		default:
		}
	}

	return false
}
