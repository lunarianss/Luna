// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/pkg/core"
	"github.com/lunarianss/Luna/internal/pkg/util"
)

func (wc *WorkspaceController) List(c *gin.Context) {
	accountID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	tenantsInfo, err := wc.tenantService.GetJoinTenants(c, accountID)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if len(tenantsInfo) == 1 {
		core.WriteResponse(c, nil, tenantsInfo[0])
		return
	}

	core.WriteResponse(c, nil, tenantsInfo)
}
