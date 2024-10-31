// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

import (
	"net/http"

	"github.com/Ryan-eng-del/hurricane/internal/pkg/code"
	"github.com/Ryan-eng-del/hurricane/internal/pkg/validation"
	"github.com/Ryan-eng-del/hurricane/pkg/errors"
	"github.com/Ryan-eng-del/hurricane/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

func GetSuccessResponse() map[string]string {
	return map[string]string{
		"message": "success",
	}
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%#+v", err)
		coder := errors.ParseCode(err)
		c.JSON(coder.HTTPStatus(), ErrResponse{Code: coder.Code(), Reference: coder.Reference(), Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func WriteBindErrResponse(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		WriteResponse(c, errors.WithCode(code.ErrValidation, validation.TranslateValidate(errs)), nil)
	} else {
		WriteResponse(c, errors.WithCode(code.ErrBind, err.Error(), nil), nil)
	}
}
