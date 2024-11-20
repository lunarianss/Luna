package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (wc *WorkspaceController) GetTenantCurrentWorkspace(c *gin.Context) {

	accountID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	currentTenant, err := wc.TenantService.GetTenantCurrentWorkspace(c, accountID)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, currentTenant)
}
