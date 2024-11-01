// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"

	dto "github.com/lunarianss/Luna/internal/api-server/dto/blog"
	"github.com/lunarianss/Luna/internal/pkg/core"
	"github.com/lunarianss/Luna/pkg/log"
)

func (bc *BlogController) List(c *gin.Context) {
	log.InfoL(c, "blog list function called.")
	params := &dto.GetBlogRequest{}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	blogs, count, err := bc.blogService.List(c, params.Page, params.PageSize)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, err, &dto.GetBlogListResponse{
		Count:    count,
		Items:    blogs,
		NextPage: params.Page + 1,
	})
}
