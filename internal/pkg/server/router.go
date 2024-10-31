// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Hurricane/pkg/errors"
)

type Router interface {
	Register(r *gin.Engine) error
	GetModule() string
}

var routers []Router

func RegisterRoute(rs ...Router) {
	routers = append(routers, rs...)
}

func (s *BaseApiServer) InitRouter(r *gin.Engine) error {
	for _, router := range routers {
		if err := router.Register(r); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("route module %s error", router.GetModule()))
		}
	}
	return nil
}
