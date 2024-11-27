package biz_entity

import (
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/common_relation"
)

type AIModelStaticConfiguration struct {
	*common.ProviderModel `yaml:",inline"`
	ParameterRules        []*ParameterRule `json:"parameter_rules" yaml:"parameter_rules"`
	Pricing               *PriceConfig     `json:"pricing" yaml:"pricing"`
	Position              int              `json:"position" yaml:"position"`
}
