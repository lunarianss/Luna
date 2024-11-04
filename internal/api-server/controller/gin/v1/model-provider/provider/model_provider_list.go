// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/lunarianss/Luna/internal/pkg/core"
	"github.com/lunarianss/Luna/pkg/log"
)

func (bc *ModelProviderController) List(c *gin.Context) {
	log.InfoL(c, "model provider list function called.")
	providerLists, err := bc.modelProviderService.GetProviderList("9ecdc361-cbc1-4c9b-8fb9-827dff4c145a", "")

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, providerLists)
}
