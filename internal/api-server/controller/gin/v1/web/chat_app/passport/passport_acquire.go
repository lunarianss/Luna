package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/core"
	"github.com/lunarianss/Luna/pkg/errors"
)

func (pc *PassportController) Acquire(c *gin.Context) {

	appCode := c.GetHeader("X-App-Code")

	if appCode == "" {
		core.WriteResponse(c, errors.WithCode(code.ErrAppCodeNotFound, "app code not exist in X-App-Code http header"), nil)
	}

	passport, err := pc.passportService.AcquirePassport(c, appCode)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, passport)
}
