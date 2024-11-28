package util

import (
	"reflect"

	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
)

const (
	PATCH_FUNCTION_NAME = "PatchZh"
)

type RecursiveObject func(obj interface{})

func PatchI18nObject(obj interface{}) {
	var recursiveObject RecursiveObject
	recursiveObject = func(obj interface{}) {

		if obj == nil {
			return
		}

		var v reflect.Value

		objValue := reflect.ValueOf(obj)

		if objValue.Kind() == reflect.Ptr {
			v = objValue.Elem()
		} else {
			v = objValue
		}

		if v.Kind() == reflect.Struct {
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)

				if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						continue
					}

					if field.Type() == reflect.TypeOf(&common.I18nObject{}) {
						method := field.MethodByName(PATCH_FUNCTION_NAME)
						if method.IsValid() && method.Type().NumIn() == 0 {
							method.Call(nil)
						}
					} else {
						recursiveObject(field.Interface())
					}
				} else if field.Kind() == reflect.Slice {
					for j := 0; j < field.Len(); j++ {
						recursiveObject(field.Index(j).Interface())
					}
				} else if field.Kind() == reflect.Interface {
					if !field.IsNil() {
						recursiveObject(field.Interface())
					}
				}
			}
		}
	}

	recursiveObject(obj)
}
