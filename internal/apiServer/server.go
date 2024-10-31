// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package master

import (
	"github.com/Ryan-eng-del/hurricane/internal/apiServer/config"
	"github.com/Ryan-eng-del/hurricane/internal/pkg/mysql"
	"github.com/Ryan-eng-del/hurricane/internal/pkg/server"
	"github.com/Ryan-eng-del/hurricane/pkg/log"
	"github.com/Ryan-eng-del/hurricane/pkg/shutdown"
)

type MasterApiServer struct {
	APIServer        *server.BaseApiServer
	GracefulShutdown *shutdown.GracefulShutdown
	AppRuntimeConfig *config.Config
}

func (s *MasterApiServer) Run() error {
	_, err := mysql.GetMySQLIns(s.AppRuntimeConfig.MySQLOptions)

	if err != nil {
		return err
	}

	if err := s.GracefulShutdown.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.APIServer.Run()
}

func createMasterApiServer(config *config.Config) (*MasterApiServer, error) {
	apiServerConfig, err := buildApiServerConfig(config)

	gs := shutdown.New()
	gs.AddShutdownManager(shutdown.NewPosixSignalManager())

	if err != nil {
		return nil, err
	}
	apiServer, err := apiServerConfig.NewServer()
	if err != nil {
		return nil, err
	}

	return &MasterApiServer{
		APIServer:        apiServer,
		GracefulShutdown: gs,
		AppRuntimeConfig: config,
	}, nil
}

func buildApiServerConfig(config *config.Config) (serverConfig *server.Config, lastErr error) {
	serverConfig = server.NewConfig()
	if lastErr = config.GenericServerRunOptions.ApplyTo(serverConfig); lastErr != nil {
		return
	}

	if lastErr = config.SecureServing.ApplyTo(serverConfig); lastErr != nil {
		return
	}

	if lastErr = config.FeatureOptions.ApplyTo(serverConfig); lastErr != nil {
		return
	}

	if lastErr = config.InsecureServing.ApplyTo(serverConfig); lastErr != nil {
		return
	}

	return
}
