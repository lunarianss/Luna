// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	Mode            string
	Middlewares     []string
	EnableHealth    bool
	EnableProfiling bool
	EnableMetrics   bool
	ReadTimeout     int
	WriteTimeout    int
	MaxHeaderBytes  int
}

func NewConfig() *Config {
	return &Config{
		EnableHealth:    true,
		Mode:            gin.DebugMode,
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
		ReadTimeout:     10,
		WriteTimeout:    10,
		MaxHeaderBytes:  1 << 20,
	}
}

func (c *Config) NewServer() (*BaseApiServer, error) {
	gin.SetMode(c.Mode)

	s := &BaseApiServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		enableHealth:        c.EnableHealth,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
		middlewares:         c.Middlewares,
		ServerOption: &ServerOption{
			ReadTimeout:    c.ReadTimeout,
			WriteTimeout:   c.WriteTimeout,
			MaxHeaderBytes: c.MaxHeaderBytes,
		},
		Engine: gin.New(),
	}

	initGenericAPIServer(s)

	return s, nil
}
