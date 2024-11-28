// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

import (
	common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"
)

type ParameterRuleItemTemplate struct {
	Label     *common.I18nObject `json:"label"`               // Label for the parameter
	Type      string             `json:"type"`                // Parameter type (e.g., "float", "int", etc.)
	Help      *common.I18nObject `json:"help"`                // Help information
	Required  bool               `json:"required"`            // Whether the parameter is required
	Default   any                `json:"default,omitempty"`   // Default value of the parameter
	Min       float64            `json:"min,omitempty"`       // Minimum value
	Max       float64            `json:"max,omitempty"`       // Maximum value
	Precision int                `json:"precision,omitempty"` // Precision for float types
	Options   []string           `json:"options,omitempty"`   // Options for string or categorical types
}

type ParameterRuleTemplate map[string]ParameterRuleItemTemplate

var ParameterRuleTemplates = ParameterRuleTemplate{
	"temperature": {
		Label: &common.I18nObject{
			En_US:   "Temperature",
			Zh_Hans: "温度",
		},
		Type: "float",
		Help: &common.I18nObject{
			En_US:   "Controls randomness. Lower temperature results in less random completions. As the temperature approaches zero, the model will become deterministic and repetitive. Higher temperature results in more random completions.",
			Zh_Hans: "温度控制随机性。较低的温度会导致较少的随机完成。随着温度接近零，模型将变得确定性和重复性。较高的温度会导致更多的随机完成。",
		},
		Required:  false,
		Default:   0.0,
		Min:       0.0,
		Max:       1.0,
		Precision: 2,
	},
	"top_p": {
		Label: &common.I18nObject{
			En_US:   "Top P",
			Zh_Hans: "Top P",
		},
		Type: "float",
		Help: &common.I18nObject{
			En_US:   "Controls diversity via nucleus sampling: 0.5 means half of all likelihood-weighted options are considered.",
			Zh_Hans: "通过核心采样控制多样性：0.5表示考虑了一半的所有可能性加权选项。",
		},
		Required:  false,
		Default:   1.0,
		Min:       0.0,
		Max:       1.0,
		Precision: 2,
	},
	"top_k": {
		Label: &common.I18nObject{
			En_US:   "Top K",
			Zh_Hans: "Top K",
		},
		Type: "int",
		Help: &common.I18nObject{
			En_US:   "Limits the number of tokens to consider for each step by keeping only the k most likely tokens.",
			Zh_Hans: "通过只保留每一步中最可能的 k 个标记来限制要考虑的标记数量。",
		},
		Required:  false,
		Default:   50,
		Min:       1,
		Max:       100,
		Precision: 0,
	},
	"presence_penalty": {
		Label: &common.I18nObject{
			En_US:   "Presence Penalty",
			Zh_Hans: "存在惩罚",
		},
		Type: "float",
		Help: &common.I18nObject{
			En_US:   "Applies a penalty to the log-probability of tokens already in the text.",
			Zh_Hans: "对文本中已有的标记的对数概率施加惩罚。",
		},
		Required:  false,
		Default:   0.0,
		Min:       0.0,
		Max:       1.0,
		Precision: 2,
	},
	"frequency_penalty": {
		Label: &common.I18nObject{
			En_US:   "Frequency Penalty",
			Zh_Hans: "频率惩罚",
		},
		Type: "float",
		Help: &common.I18nObject{
			En_US:   "Applies a penalty to the log-probability of tokens that appear in the text.",
			Zh_Hans: "对文本中出现的标记的对数概率施加惩罚。",
		},
		Required:  false,
		Default:   0.0,
		Min:       0.0,
		Max:       1.0,
		Precision: 2,
	},
	"max_tokens": {
		Label: &common.I18nObject{
			En_US:   "Max Tokens",
			Zh_Hans: "最大标记",
		},
		Type: "int",
		Help: &common.I18nObject{
			En_US:   "Specifies the upper limit on the length of generated results. If the generated results are truncated, you can increase this parameter.",
			Zh_Hans: "指定生成结果长度的上限。如果生成结果截断，可以调大该参数。",
		},
		Required:  false,
		Default:   64,
		Min:       1,
		Max:       2048,
		Precision: 0,
	},
	"response_format": {
		Label: &common.I18nObject{
			En_US:   "Response Format",
			Zh_Hans: "回复格式",
		},
		Type: "string",
		Help: &common.I18nObject{
			En_US:   "Set a response format, ensure the output from llm is a valid code block as possible, such as JSON, XML, etc.",
			Zh_Hans: "设置一个返回格式，确保llm的输出尽可能是有效的代码块，如JSON、XML等",
		},
		Required: false,
		Options:  []string{"JSON", "XML"},
	},
	"json_schema": {
		Label: &common.I18nObject{
			En_US: "JSON Schema",
		},
		Type: "text",
		Help: &common.I18nObject{
			En_US:   "Set a response json schema will ensure LLM to adhere it.",
			Zh_Hans: "设置返回的json schema，llm将按照它返回",
		},
		Required: false,
	},
}
