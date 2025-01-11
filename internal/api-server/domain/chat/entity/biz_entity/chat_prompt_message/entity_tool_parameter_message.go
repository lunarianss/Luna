package biz_entity

type PromptMessageToolProperty struct {
	Type        string   `json:"string,omitempty"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

type PromptMessageToolProperties map[string]*PromptMessageToolProperty

type PromptMessageToolParameter struct {
	Type       string                      `json:"type"`
	Properties PromptMessageToolProperties `json:"properties"`
	Required   []string                    `json:"required,omitempty"`
}

func NewPromptMessageToolParameter() *PromptMessageToolParameter {
	return &PromptMessageToolParameter{
		Required:   make([]string, 0),
		Type:       "object",
		Properties: make(PromptMessageToolProperties, 0),
	}
}

type PromptMessageTool struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Parameters  *PromptMessageToolParameter `json:"parameters"`
}

type PromptMessageFunction struct {
	Type     string             `json:"type"`
	Function *PromptMessageTool `json:"function"`
}

func NewFunctionTools(function *PromptMessageTool) *PromptMessageFunction {
	return &PromptMessageFunction{
		Type:     "function",
		Function: function,
	}
}
