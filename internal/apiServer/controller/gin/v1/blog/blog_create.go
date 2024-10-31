// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	dto "github.com/Ryan-eng-del/hurricane/internal/apiServer/dto/blog"
	"github.com/Ryan-eng-del/hurricane/internal/pkg/core"
	"github.com/Ryan-eng-del/hurricane/pkg/log"
	"github.com/gin-gonic/gin"
)

func (bc *BlogController) Create(c *gin.Context) {
	log.InfoL(c, "blog get function called.")
	params := &dto.CreateBlogRequest{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	blog, err := bc.blogService.Create(c, params)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, blog)
}
