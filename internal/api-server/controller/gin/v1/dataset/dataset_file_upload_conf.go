package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/pkg/core"
)

func (dc *DatasetController) GetFileUploadConfiguration(c *gin.Context) {
	config, err := dc.datasetService.GetFileUploadConfiguration(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, config)
}
