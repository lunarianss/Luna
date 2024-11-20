package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (wc *WorkspaceController) List(c *gin.Context) {
	accountID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	tenantsInfo, err := wc.TenantService.GetJoinTenants(c, accountID)

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
