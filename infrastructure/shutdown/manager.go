// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package shutdown

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lunarianss/Luna/infrastructure/log"
)

const Name = "PosixSignalManager"

type PosixSignalManager struct {
	signals []os.Signal
}

// NewPosixSignalManager initializes the PosixSignalManager.
// As arguments you can provide os.Signal-s to listen to, if none are given,
// it will default to SIGINT and SIGTERM.
func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt
		sig[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{
		signals: sig,
	}
}

// GetName returns name of this ShutdownManager.
func (posixSignalManager *PosixSignalManager) GetName() string {
	return Name
}

// ShutdownStart does nothing.
func (posixSignalManager *PosixSignalManager) ShutdownStart() error {
	log.Info("开始 shutdown")

	return nil
}

// ShutdownFinish exits the app with os.Exit(0).
func (posixSignalManager *PosixSignalManager) ShutdownFinish() error {
	log.Info("结束 shutdown")
	os.Exit(0)

	return nil
}

func (posixSignalManager *PosixSignalManager) Start(gc GracefulShutdownInterface) error {
	go func() {
		in_signal := make(chan os.Signal, 1)
		signal.Notify(in_signal, posixSignalManager.signals...)

		<-in_signal

		gc.StartShutdown(posixSignalManager)
	}()

	return nil
}
