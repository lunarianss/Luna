// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
)

func (fc *FeatureController) GetSystemConfigs(c *gin.Context) {
	systemConfig, err := fc.featureService.GetSystemConfig()

	if err != nil {
		core.WriteResponse(c, err, nil)
	}

	core.WriteResponse(c, nil, systemConfig)
}
