// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package master

import (
	"github.com/Ryan-eng-del/hurricane/internal/pkg/server"

	"github.com/gin-gonic/gin"
)

var masterMiddlewares = make([]gin.HandlerFunc, 0)

type MasterMiddleware struct{}

//nolint: gochecknoinits
func init() {
	server.AddMiddlewares(&MasterMiddleware{})
}

func (m *MasterMiddleware) Fetch() []gin.HandlerFunc {
	return masterMiddlewares
}
