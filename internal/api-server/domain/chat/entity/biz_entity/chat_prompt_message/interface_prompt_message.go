package biz_entity

type IPromptMessage interface {
	GetRole() string
	GetContent() string
	GetName() string
	ConvertToRequestData() (map[string]interface{}, error)
}
