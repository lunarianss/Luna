// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package master

import (
	"github.com/lunarianss/Luna/internal/api-server/config"

	_ "github.com/lunarianss/Luna/internal/api-server/facade"
	_ "github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers"
	_ "github.com/lunarianss/Luna/internal/api-server/validation"

	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/infrastructure/shutdown"
	"github.com/lunarianss/Luna/internal/infrastructure/email"
	"github.com/lunarianss/Luna/internal/infrastructure/jwt"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
	"github.com/lunarianss/Luna/internal/infrastructure/redis"
	"github.com/lunarianss/Luna/internal/infrastructure/server"
	"github.com/lunarianss/Luna/internal/infrastructure/validation"
)

type LunaApiServer struct {
	APIServer        *server.BaseApiServer
	GracefulShutdown *shutdown.GracefulShutdown
	AppRuntimeConfig *config.Config
}

func (s *LunaApiServer) Run() error {
	// Register the module of master router and validator
	if err := validation.InitAppValidator(); err != nil {
		return err
	}

	if _, err := redis.GetRedisIns(s.AppRuntimeConfig.RedisOptions); err != nil {
		return err
	}

	if _, err := mysql.GetMySQLIns(s.AppRuntimeConfig.MySQLOptions); err != nil {
		return err
	}

	if _, err := email.GetEmailSMTPIns(s.AppRuntimeConfig.EmailOptions); err != nil {
		return err
	}

	_ = jwt.NewJWT(s.AppRuntimeConfig.JwtOptions.Key)

	if err := s.APIServer.InitRouter(s.APIServer.Engine); err != nil {
		return err
	}

	if err := s.GracefulShutdown.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.APIServer.Run()
}

func createLunaApiServer(config *config.Config) (*LunaApiServer, error) {
	apiServerConfig, err := buildApiServerConfig(config)

	gs := shutdown.New()
	gs.AddShutdownManager(shutdown.NewPosixSignalManager())

	gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		log.Info("call shutdown callback")
		log.Info("finish shutdown callback")
		return nil
	}))

	if err != nil {
		return nil, err
	}
	apiServer, err := apiServerConfig.NewServer()
	if err != nil {
		return nil, err
	}

	return &LunaApiServer{
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
