package site

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/pkg/util"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (s *WebSiteController) Retrieve(c *gin.Context) {
	appID, appCode, endUserID, err := util.GetWebAppFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	wenSite, err := s.webSiteService.GetSiteByWebToken(c, appID, endUserID, appCode)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, wenSite)
}
