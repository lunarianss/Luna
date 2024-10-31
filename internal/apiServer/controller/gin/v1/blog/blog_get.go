// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"strconv"

	"github.com/Ryan-eng-del/hurricane/internal/pkg/code"
	"github.com/Ryan-eng-del/hurricane/internal/pkg/core"
	"github.com/Ryan-eng-del/hurricane/pkg/errors"
	"github.com/Ryan-eng-del/hurricane/pkg/log"
	"github.com/gin-gonic/gin"
)

func (bc *BlogController) Get(c *gin.Context) {
	log.InfoL(c, "blog get function called.")

	blogIdStr := c.Param("blogId")
	blogId, err := strconv.ParseInt(blogIdStr, 10, 64)

	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrRestfulId, err.Error()), nil)
		return
	}

	blog, err := bc.blogService.Get(c, blogId)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, blog)
}
