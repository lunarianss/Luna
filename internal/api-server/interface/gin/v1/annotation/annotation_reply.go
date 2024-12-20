package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func (ac *AnnotationController) AnnotationReply(c *gin.Context) {
	params := &dto.ApplyAnnotationRequestBody{}
	paramsUrl := &dto.ApplyAnnotationRequestUrl{}

	if err := c.ShouldBindUri(paramsUrl); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	if err := c.ShouldBind(params); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if paramsUrl.Action == "enable" {
		enableAnnotation, err := ac.annotation.EnableAppAnnotation(c, paramsUrl.AppID, userID, params)
		if err != nil {
			core.WriteResponse(c, err, nil)
			return
		} else {
			core.WriteResponse(c, nil, enableAnnotation)
		}
	}
}
