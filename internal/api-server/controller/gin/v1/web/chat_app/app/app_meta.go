package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ac *WebAppController) AppMeta(c *gin.Context) {

	c.JSON(http.StatusOK, map[string]map[string]any{
		"tool_icons": {},
	})
}
