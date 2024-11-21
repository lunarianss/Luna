// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/lunarianss/Luna/internal/pkg/middleware"
	"github.com/lunarianss/Luna/pkg/log"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
)

type BaseApiServer struct {
	middlewares                                  []string
	SecureServingInfo                            *SecureServingInfo
	InsecureServingInfo                          *InsecureServingInfo
	enableHealth, enableMetrics, enableProfiling bool
	InsecureServer, SecureServer                 *http.Server
	*ServerOption
	*gin.Engine
}

func initGenericAPIServer(s *BaseApiServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

// InstallAPIs install generic apis.
func (s *BaseApiServer) InstallAPIs() {
	// install healthz handler
	if s.enableHealth {
		s.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})
	}

	// install metric handler
	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	// install pprof handler
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	s.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"AppVersion": "v1"})
	})
}

// Setup do some setup work for gin engine.
func (s *BaseApiServer) Setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// InstallMiddlewares install generic middlewares.
func (s *BaseApiServer) InstallMiddlewares() {
	// necessary middlewares
	s.Use(middleware.Recovery())
	s.Use(middleware.Context())
	s.Use(middleware.RequestID())

	// install custom middlewares
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)

			continue
		}

		log.Infof("install middleware: %s", m)
		s.Use(mw)
	}

	for _, injectMiddleware := range injectMiddlewares {
		s.Use(injectMiddleware)
	}
}

// Run spawns the http server. It only returns when the port cannot be listened on initially.
func (s *BaseApiServer) Run() error {
	// For scalability, use custom HTTP configuration mode here
	s.InsecureServer = &http.Server{
		Addr:           s.InsecureServingInfo.Address(),
		Handler:        s,
		ReadTimeout:    time.Duration(s.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.WriteTimeout) * time.Second,
		MaxHeaderBytes: s.MaxHeaderBytes,
	}

	// For scalability, use custom HTTP configuration mode here
	s.SecureServer = &http.Server{
		Addr:           s.SecureServingInfo.Address(),
		Handler:        s,
		ReadTimeout:    time.Duration(s.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.WriteTimeout) * time.Second,
		MaxHeaderBytes: s.MaxHeaderBytes,
	}

	var eg errgroup.Group
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	eg.Go(func() error {
		log.Infof(
			"Start to listening the incoming requests on http address: %s",
			color.HiGreenString(s.InsecureServingInfo.Address()),
		)
		if err := s.InsecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address)

		return nil
	})

	eg.Go(func() error {
		key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
		if cert == "" || key == "" || s.SecureServingInfo.BindPort == 0 {
			return nil
		}

		log.Infof(
			"Start to listening the incoming requests on https address: %s",
			color.HiGreenString(s.SecureServer.Addr),
		)

		if err := s.SecureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.SecureServingInfo.Address())

		return nil
	})

	// Ping the server to make sure the router is working.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if s.enableHealth {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func (s *BaseApiServer) ping(ctx context.Context) (err error) {
	address := s.InsecureServingInfo.Address()
	url := fmt.Sprintf("http://%s/healthz", address)
	if strings.Contains(address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(address, ":")[1])
	}

	for {
		// Change NewRequest to NewRequestWithContext and pass context it
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		// Ping the server by sending a GET request to `/healthz`.
		resp, err := http.DefaultClient.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("The router has been deployed successfully.")
			resp.Body.Close()

			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
}
