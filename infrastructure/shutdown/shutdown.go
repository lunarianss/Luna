// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package shutdown

import (
	"sync"
)

type ErrorHandler interface {
	OnError(err error)
}

type ErrorFunc func(error)

func (e ErrorFunc) OnError(err error) {
	e(err)
}

type ShutdownCallback interface {
	OnShutdown(string) error
}

type ShutdownFunc func(string) error

func (f ShutdownFunc) OnShutdown(manager string) error {
	return f(manager)
}

type ShutdownManager interface {
	GetName() string
	Start(gs GracefulShutdownInterface) error
	ShutdownStart() error
	ShutdownFinish() error
}

type GracefulShutdownInterface interface {
	AddShutdownCallback(shutdownCallback ShutdownCallback)
	ReportError(err error)
	StartShutdown(sm ShutdownManager)
}

type GracefulShutdown struct {
	callbacks   []ShutdownCallback
	errorHandle ErrorHandler
	managers    []ShutdownManager
}

// New initializes GracefulShutdown.
func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks: make([]ShutdownCallback, 0, 10),
		managers:  make([]ShutdownManager, 0, 3),
	}
}

func (gs *GracefulShutdown) AddShutdownCallback(shutdownCallback ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, shutdownCallback)
}

func (gs *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	gs.errorHandle = errorHandler
}

// ReportError is a function that can be used to report errors to
// ErrorHandler. It is used in ShutdownManagers.
func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandle != nil {
		gs.errorHandle.OnError(err)
	}
}

func (gs *GracefulShutdown) StartShutdown(sm ShutdownManager) {
	gs.ReportError(sm.ShutdownStart())
	var wg sync.WaitGroup

	for _, callback := range gs.callbacks {
		wg.Add(1)
		go func(callback ShutdownCallback) {
			defer wg.Done()
			gs.ReportError(callback.OnShutdown(sm.GetName()))
		}(callback)
	}

	wg.Wait()
	gs.ReportError(sm.ShutdownFinish())
}

func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}

	return nil
}

// AddShutdownManager adds a ShutdownManager that will listen to shutdown requests.
func (gs *GracefulShutdown) AddShutdownManager(manager ShutdownManager) {
	gs.managers = append(gs.managers, manager)
}
