package main

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// 自定义校验：检查时间格式是否符合指定的格式
func validateDatetime(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02 15:04", fl.Field().String())
	return err == nil
}

// 自定义校验：检查整数是否在指定的范围内
func validateIntRange(min, max int) validator.Func {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().Int()
		return val >= int64(min) && val <= int64(max)
	}
}

// 请求参数结构体
type RequestParams struct {
	Keyword          string `json:"keyword" validate:"required"`
	Start            string `json:"start" validate:"datetime"`
	End              string `json:"end" validate:"datetime"`
	AnnotationStatus string `json:"annotation_status" validate:"oneof=annotated not_annotated all"`
	MessageCountGte  int    `json:"message_count_gte" validate:"omitempty,min=1,max=99999"`
	Page             int    `json:"page" validate:"omitempty,min=1,max=99999"`
	Limit            int    `json:"limit" validate:"omitempty,min=1,max=100"`
	SortBy           string `json:"sort_by" validate:"omitempty,oneof=created_at -created_at updated_at -updated_at"`
}

func main() {
	// 创建一个 validator 实例
	validate := validator.New()

	// 注册自定义校验规则
	validate.RegisterValidation("datetime", validateDatetime)
	validate.RegisterValidation("int_range", validateIntRange(1, 99999))

	// 模拟的请求参数
	params := RequestParams{
		Keyword:          "search",
		Start:            "2023-10-01 14:30",
		End:              "2023-10-02 14:30",
		AnnotationStatus: "annotated",
		MessageCountGte:  100,
		Page:             1,
		Limit:            20,
		SortBy:           "updated_at",
	}

	// 校验请求参数
	err := validate.Struct(params)
	if err != nil {
		// 输出错误信息
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("Validation failed:", err.Namespace(), err.Field(), err.Tag())
		}
	} else {
		// 校验成功
		fmt.Println("Validation passed!")
	}
}
