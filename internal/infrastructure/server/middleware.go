// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import "github.com/gin-gonic/gin"

type Middleware interface {
	Fetch() []gin.HandlerFunc
}

var injectMiddlewares = []gin.HandlerFunc{}

func AddMiddlewares(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		injectMiddlewares = append(injectMiddlewares, middleware.Fetch()...)
	}
}
