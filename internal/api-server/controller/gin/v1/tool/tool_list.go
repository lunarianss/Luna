package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (tc *ToolController) List(c *gin.Context) {
	core.WriteResponse(c, nil, []any{})
}
