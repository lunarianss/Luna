// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

import common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"

type AIModelStaticConfiguration struct {
	*common.ProviderModel `yaml:",inline"`
	ParameterRules        []*ParameterRule `json:"parameter_rules" yaml:"parameter_rules"`
	Pricing               *PriceConfig     `json:"pricing" yaml:"pricing"`
	Position              int              `json:"position" yaml:"position"`
}
