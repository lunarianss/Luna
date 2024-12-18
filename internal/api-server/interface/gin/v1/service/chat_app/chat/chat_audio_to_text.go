// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"io"
	"path/filepath"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/core"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

const FILE_SIZE = 30
const FILE_SIZE_LIMIT = FILE_SIZE * 1024 * 1024

var ALLOWED_EXTENSIONS = []string{".mp3", ".mp4", ".mpeg", ".mpga", ".m4a", ".wav", ".webm", ".amr"}
var ALLOWED_MIMETYPE = []string{"audio/mp3", "audio/mp4", "audio/mpeg", "audio/mpga", "audio/m4a", "audio/wav", "audio/webm", "audio/amr"}

func (ac *ServiceChatController) AudioToChatMessage(c *gin.Context) {

	var audioFileContent []byte

	audioFile, err := c.FormFile("file")

	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrRunTimeCaller, err.Error()), nil)
		return
	}

	audioFileOpen, err := audioFile.Open()

	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrRunTimeCaller, err.Error()), nil)
		return
	}

	defer audioFileOpen.Close()

	audioFileContent, err = io.ReadAll(audioFileOpen)

	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrRunTimeCaller, err.Error()), nil)
		return
	}

	if audioFile.Size == 0 {
		core.WriteResponse(c, errors.WithCode(code.ErrAudioFileEmpty, ""), nil)
		return
	}

	if audioFile.Size > FILE_SIZE_LIMIT {
		core.WriteResponse(c, errors.WithCode(code.ErrAudioFileToLarge, ""), nil)
		return
	}

	fileExtension := filepath.Ext(audioFile.Filename)

	if !slices.Contains(ALLOWED_EXTENSIONS, fileExtension) {
		core.WriteResponse(c, errors.WithCode(code.ErrAudioType, ""), nil)
		return
	}

	contentType := audioFile.Header.Get("Content-Type")

	if contentType == "" || !slices.Contains(ALLOWED_MIMETYPE, contentType) {
		core.WriteResponse(c, errors.WithCode(code.ErrAudioType, ""), nil)
		return
	}

	appID, _, endUserID, err := util.GetWebAppFromGin(c)

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if transcription, err := ac.serviceChatService.AudioToText(c, audioFileContent, audioFile.Filename, appID, endUserID); err != nil {
		core.WriteResponse(c, err, nil)
	} else {
		core.WriteResponse(c, nil, transcription)
	}
}
