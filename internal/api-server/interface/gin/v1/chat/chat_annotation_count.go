package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
)

func (cc *ChatController) GetAnnotationCount(c *gin.Context) {
	core.WriteResponse(c, nil, map[string]interface{}{"count": 0})
}
