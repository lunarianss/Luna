// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Hurricane/internal/pkg/code"
	"github.com/lunarianss/Hurricane/internal/pkg/core"
	"github.com/lunarianss/Hurricane/pkg/errors"
	"github.com/lunarianss/Hurricane/pkg/log"
)

func (bc *BlogController) Delete(c *gin.Context) {
	log.InfoL(c, "blog delete function called.")

	blogIdStr := c.Param("blogId")
	blogId, err := strconv.ParseInt(blogIdStr, 10, 64)

	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrRestfulId, err.Error()), nil)
		return
	}

	err = bc.blogService.Delete(c, blogId)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
