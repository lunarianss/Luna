package controller

import (
	"github.com/gin-gonic/gin"
	appDto "github.com/lunarianss/Luna/internal/api-server/dto/app"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/statistic"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

func (sc *StatisticController) ActiveUsers(c *gin.Context) {

	paramsQuery := &dto.StatisticQuery{}
	paramsUrl := &appDto.AppDetailRequest{}

	if err := c.ShouldBindUri(paramsUrl); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}
	if err := c.ShouldBind(paramsQuery); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	userID, err := util.GetUserIDFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	statisticData, err := sc.statisticService.DailyUsers(c, paramsUrl.AppID, userID, paramsQuery.Start, paramsQuery.End)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, statisticData)
}
