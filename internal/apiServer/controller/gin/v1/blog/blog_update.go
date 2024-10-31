// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Hurricane/internal/apiServer/dto/blog"
	"github.com/lunarianss/Hurricane/internal/pkg/code"
	"github.com/lunarianss/Hurricane/internal/pkg/core"
	"github.com/lunarianss/Hurricane/pkg/errors"
	"github.com/lunarianss/Hurricane/pkg/log"
)

func (bc *BlogController) Update(c *gin.Context) {
	log.InfoL(c, "blog update function called.")

	blogIdStr := c.Param("blogId")
	blogId, err := strconv.ParseInt(blogIdStr, 10, 64)

	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrRestfulId, err.Error()), nil)
		return
	}

	params := &dto.UpdateBlogRequest{}

	if err := c.ShouldBindJSON(&params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	blog, err := bc.blogService.Update(c, blogId, params)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, blog)
}