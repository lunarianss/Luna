// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package route

import "github.com/gin-gonic/gin"

type staticRoute struct {
}

func (r *staticRoute) Register(g *gin.Engine) error {
	v1 := g.Group("/v1")
	v1.Static("/static", "./assets")
	return nil
}

func (r *staticRoute) GetModule() string {
	return "static"
}
