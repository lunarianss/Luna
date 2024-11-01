// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	vtor "github.com/lunarianss/Hurricane/internal/pkg/validation"
)

type blogValidation struct{}

func (bv *blogValidation) Register() error {
	trans, err := vtor.GetGlobalTrans()
	if err != nil {
		return err
	}

	validate, err := vtor.GetGlobalValidate()

	if err != nil {
		return err
	}

	validate.RegisterValidation("valid_username", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "admin"
	})

	validate.RegisterTranslation("valid_username", trans, func(ut ut.Translator) error {
		return ut.Add("valid_username", "{0} is not correct", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_username", fe.Field())
		return t
	})

	return nil
}

func (bv *blogValidation) Module() string {
	return "blog"
}
