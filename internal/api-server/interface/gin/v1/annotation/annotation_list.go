package controller

import (
	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func (ac *AnnotationController) AnnotationList(c *gin.Context) {

	paramsUrl := &dto.AppDetailRequest{}
	params := &dto.ListAnnotationsArgs{}

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

	annotations, err := ac.annotation.ListAnnotations(c, paramsUrl.AppID, userID, params)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, annotations)

}
