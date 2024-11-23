package controller

import "github.com/lunarianss/Luna/internal/api-server/service"

type WorkspaceController struct {
	tenantService *service.TenantService
}

func NewWorkspaceController(tenantService *service.TenantService) *WorkspaceController {
	return &WorkspaceController{
		tenantService: tenantService,
	}
}
