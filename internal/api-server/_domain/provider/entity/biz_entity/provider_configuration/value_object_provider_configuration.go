package biz_entity

import (
	common "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/common_relation"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type ModelStatus string

const (
	ACTIVE         ModelStatus = "active"
	NO_CONFIGURE   ModelStatus = "no-configure"
	QUOTA_EXCEEDED ModelStatus = "quota-exceeded"
	NO_PERMISSION  ModelStatus = "no-permission"
	DISABLED       ModelStatus = "disabled"
)

type QuotaUnit string

const (
	TIMES   QuotaUnit = "times"
	TOKENS  QuotaUnit = "tokens"
	CREDITS QuotaUnit = "credits"
)

type ModelSettings struct {
	Model     string
	ModelType common.ModelType
	Enabled   bool
}

type SystemConfiguration struct {
	Enabled             bool
	CurrentQuotaType    model.ProviderQuotaType
	QuotaConfigurations []*QuotaConfiguration
	Credentials         interface{}
}

type RestrictModels struct {
	Model         string
	BaseModelName string
	ModelType     string
}

type QuotaConfiguration struct {
	QuotaType      model.ProviderQuotaType
	QuotaUnit      QuotaUnit
	QuotaLimit     int
	QuotaUsed      int
	IsValid        bool
	RestrictModels []*RestrictModels
}

type CustomConfiguration struct {
	Provider *CustomProviderConfiguration
	Models   []*CustomModelConfiguration
}

type CustomModelConfiguration struct {
	Model       string
	ModelType   string
	Credentials map[string]interface{}
}

type CustomProviderConfiguration struct {
	Credentials interface{}
}
