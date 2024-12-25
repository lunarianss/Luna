package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func (ac *AnnotationController) AnnotationReplyStatus(c *gin.Context) {

	paramsUrl := &dto.ApplyAnnotationStatusRequestUrl{}

	if err := c.ShouldBindUri(paramsUrl); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if paramsUrl.Action == "enable" {
		enableAnnotation, err := ac.annotation.EnableAppAnnotationStatus(c, paramsUrl.AppID, userID, paramsUrl.JobID, paramsUrl.Action)
		if err != nil {
			core.WriteResponse(c, err, nil)
			return
		} else {
			core.WriteResponse(c, nil, enableAnnotation)
		}
	}
}
