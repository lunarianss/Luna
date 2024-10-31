// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package master

import (
	"github.com/Ryan-eng-del/hurricane/internal/pkg/server"

	"github.com/gin-gonic/gin"
)

//nolint: gochecknoinits
func init() {
	server.RegisterRoute(&MasterRoute{})
}

type MasterRoute struct{}

func (mr *MasterRoute) Register(r *gin.Engine) {
}
