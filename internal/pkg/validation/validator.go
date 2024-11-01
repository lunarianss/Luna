// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/lunarianss/Hurricane/pkg/errors"
)

var globalTrans ut.Translator
var validators []Validator
var globalValidate *validator.Validate

var transNilError = errors.New("global trans is nil")
var validateNilError = errors.New("global validator is nil")

type Validator interface {
	Register() error
	Module() string
}

func GetGlobalTrans() (ut.Translator, error) {
	if globalTrans == nil {
		return nil, transNilError
	}
	return globalTrans, nil
}

func GetGlobalValidate() (*validator.Validate, error) {
	if globalValidate == nil {
		return nil, validateNilError
	}
	return globalValidate, nil
}

func InitAppValidator() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		globalValidate = v
		en := en.New()
		uni := ut.New(en, en)
		trans, _ := uni.GetTranslator("en")
		globalTrans = trans
		v.SetTagName("validate")

		if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
			return err
		}

		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			return field.Tag.Get("comment")
		})

		for _, validator := range validators {
			if err := validator.Register(); err != nil {
				return errors.WithMessage(err, fmt.Sprintf("register module %s validator error", validator.Module()))
			}
		}
	}
	return nil
}

func RegisterValidator(validator Validator) {
	validators = append(validators, validator)
}

func TranslateValidate(err error) string {
	errs := err.(validator.ValidationErrors)
	sliceErrs := []string{}
	for _, e := range errs {
		if trans, err := GetGlobalTrans(); err != nil {
			return transNilError.Error()
		} else {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
	}
	return strings.Join(sliceErrs, ",\n")
}
