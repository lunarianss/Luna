// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

type AppMode string

const (
	COMPLETION    AppMode = "completion"
	WORKFLOW      AppMode = "workflow"
	CHAT          AppMode = "chat"
	ADVANCED_CHAT AppMode = "advanced-chat"
	AGENT_CHAT    AppMode = "agent-chat"
	CHANNEL       AppMode = "channel"
)

// App holds the basic app configuration.
type App struct {
	Mode       string `json:"mode"`
	EnableSite int    `json:"enable_site"`
	EnableAPI  int    `json:"enable_api"`
}

// ModelConfig holds the model and additional configurations.
type ModelConfig struct {
	Model         ModelInfo       `json:"model"`
	UserInputForm []UserInputForm `json:"user_input_form"`
	PrePrompt     string          `json:"pre_prompt"`
}

// AppTemplate holds the template for each mode.
type AppTemplate struct {
	App         App          `json:"app"`
	ModelConfig *ModelConfig `json:"model_config,omitempty"`
}

// Initialize default app templates
var DefaultAppTemplates = map[AppMode]AppTemplate{
	WORKFLOW: {
		App: App{
			Mode:       string(WORKFLOW),
			EnableSite: 1,
			EnableAPI:  1,
		},
	},
	COMPLETION: {
		App: App{
			Mode:       string(COMPLETION),
			EnableSite: 1,
			EnableAPI:  1,
		},
		ModelConfig: &ModelConfig{
			Model: ModelInfo{
				Provider:         "openai",
				Name:             "gpt-4o",
				Mode:             "chat",
				CompletionParams: map[string]interface{}{},
			},
			PrePrompt: "{{query}}",
			// UserInputForm: []map[string]map[string]interface{}{
			// 	{
			// 		"params": {
			// 			"label":    "Query",
			// 			"variable": "query",
			// 			"required": true,
			// 			"default":  "",
			// 		},
			// 	},
			// },
		},
	},
	CHAT: {
		App: App{
			Mode:       string(CHAT),
			EnableSite: 1,
			EnableAPI:  1,
		},
		ModelConfig: &ModelConfig{
			Model: ModelInfo{
				Provider:         "openai",
				Name:             "gpt-4o",
				Mode:             "chat",
				CompletionParams: map[string]interface{}{},
			},
		},
	},
	ADVANCED_CHAT: {
		App: App{
			Mode:       string(ADVANCED_CHAT),
			EnableSite: 1,
			EnableAPI:  1,
		},
	},
	AGENT_CHAT: {
		App: App{
			Mode:       string(AGENT_CHAT),
			EnableSite: 1,
			EnableAPI:  1,
		},
		ModelConfig: &ModelConfig{
			Model: ModelInfo{
				Provider:         "openai",
				Name:             "gpt-4o",
				Mode:             "chat",
				CompletionParams: map[string]interface{}{},
			},
		},
	},
}
