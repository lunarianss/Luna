// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	"slices"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	biz_entity_provider_common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
	vtor "github.com/lunarianss/Luna/internal/infrastructure/validation"
)

type modelValidation struct{}

func (bv *modelValidation) Register() error {
	trans, err := vtor.GetGlobalTrans()
	if err != nil {
		return err
	}

	validate, err := vtor.GetGlobalValidate()

	if err != nil {
		return err
	}

	validate.RegisterValidation("valid_model_type", func(fl validator.FieldLevel) bool {
		return slices.Contains(biz_entity_provider_common.ModelTypeEnums, fl.Field().String())
	})

	validate.RegisterTranslation("valid_model_type", trans, func(ut ut.Translator) error {
		return ut.Add("valid_model_type", "{0} is not correct ", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_model_type", fe.Field())
		return t
	})

	return nil
}

func (bv *modelValidation) Module() string {
	return "model"
}
