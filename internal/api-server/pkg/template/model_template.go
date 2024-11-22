package template

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

// Model holds the model-specific configuration.
type Model struct {
	Provider         string                 `json:"provider"`
	Name             string                 `json:"name"`
	Mode             string                 `json:"mode"`
	CompletionParams map[string]interface{} `json:"completion_params"`
}

// ModelConfig holds the model and additional configurations.
type ModelConfig struct {
	Model         Model                               `json:"model"`
	UserInputForm []map[string]map[string]interface{} `json:"user_input_form"`
	PrePrompt     string                              `json:"pre_prompt"`
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
			Model: Model{
				Provider:         "openai",
				Name:             "gpt-4o",
				Mode:             "chat",
				CompletionParams: map[string]interface{}{},
			},
			PrePrompt: "{{query}}",
			UserInputForm: []map[string]map[string]interface{}{
				{
					"params": {
						"label":    "Query",
						"variable": "query",
						"required": true,
						"default":  "",
					},
				},
			},
		},
	},
	CHAT: {
		App: App{
			Mode:       string(CHAT),
			EnableSite: 1,
			EnableAPI:  1,
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
			Model: Model{
				Provider:         "openai",
				Name:             "gpt-4o",
				Mode:             "chat",
				CompletionParams: map[string]interface{}{},
			},
		},
	},
}
