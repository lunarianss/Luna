package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/file"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
)

func (fc *FileController) PreviewFile(c *gin.Context) {
	query := &dto.PreviewFileQuery{}
	url := &dto.PreviewFileUri{}

	if err := c.ShouldBind(query); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	if err := c.ShouldBindUri(url); err != nil {
		core.WriteBindErrResponse(c, err)
		return
	}

	filename := url.Filename

	fileID := strings.Split(filename, ".")[0]

	if err := fc.fileService.PreviewFile(c, fileID, query); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
}
