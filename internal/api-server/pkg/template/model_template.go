package template

import (
	"encoding/json"

	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

// App holds the basic app configuration.
type App struct {
	Mode       string `json:"mode"`
	EnableSite bool   `json:"enable_site"`
	EnableAPI  bool   `json:"enable_api"`
}

// Model holds the model-specific configuration.
type Model struct {
	Provider         string                 `json:"provider"`
	Name             string                 `json:"name"`
	Mode             string                 `json:"mode"`
	CompletionParams map[string]interface{} `json:"completion_params"`
}

// ModelConfig holds the model and additional configurations.
type ModelConfig struct {
	Model         Model  `json:"model"`
	UserInputForm string `json:"user_input_form"`
	PrePrompt     string `json:"pre_prompt"`
}

// AppTemplate holds the template for each mode.
type AppTemplate struct {
	App         App          `json:"app"`
	ModelConfig *ModelConfig `json:"model_config,omitempty"`
}

// Initialize default app templates
var DefaultAppTemplates = map[model.AppMode]AppTemplate{
	model.WORKFLOW: {
		App: App{
			Mode:       string(model.WORKFLOW),
			EnableSite: true,
			EnableAPI:  true,
		},
	},
	model.COMPLETION: {
		App: App{
			Mode:       string(model.COMPLETION),
			EnableSite: true,
			EnableAPI:  true,
		},
		ModelConfig: &ModelConfig{
			Model: Model{
				Provider:         "openai",
				Name:             "gpt-4o",
				Mode:             "chat",
				CompletionParams: map[string]interface{}{},
			},
			UserInputForm: string(json.RawMessage(`[{
				"paragraph": {
						"label": "Query",
						"variable": "query",
						"required": true,
						"default": ""
				}
		}]`)),
			PrePrompt: "{{query}}",
		},
	},
	model.CHAT: {
		App: App{
			Mode:       string(model.CHAT),
			EnableSite: true,
			EnableAPI:  true,
		},
		ModelConfig: &ModelConfig{
			Model: Model{
				Provider:         "openai",
				Name:             "gpt-4o",
				Mode:             "chat",
				CompletionParams: map[string]interface{}{},
			},
		},
	},
	model.ADVANCED_CHAT: {
		App: App{
			Mode:       string(model.ADVANCED_CHAT),
			EnableSite: true,
			EnableAPI:  true,
		},
	},
	model.AGENT_CHAT: {
		App: App{
			Mode:       string(model.AGENT_CHAT),
			EnableSite: true,
			EnableAPI:  true,
		},
		ModelConfig: &ModelConfig{
			Model: Model{
				Provider:         "openai",
				Name:             "gpt-4o",
				Mode:             "chat",
				CompletionParams: map[string]interface{}{},
			},
		},
	},
}
