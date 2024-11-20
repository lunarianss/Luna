package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *SetupController) ValidateSetup(c *gin.Context) {

	c.JSON(http.StatusOK, map[string]any{
		"step": "finished",
	})
}
