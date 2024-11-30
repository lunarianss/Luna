package biz_entity

type SimpleChatPromptConfig struct {
	HumanPrefix        string   `json:"human_prefix"`
	AssistantPrefix    string   `json:"assistant_prefix"`
	ContextPrompt      string   `json:"context_prompt"`
	HistoriesPrompt    string   `json:"histories_prompt"`
	SystemPromptOrders []string `json:"system_prompt_orders"`
	QueryPrompt        string   `json:"query_prompt"`
	Stops              []string `json:"stops"`
}

type SimpleChatPromptTransformConfig struct {
	PromptTemplate      IPromptTemplateParser `json:"prompt_template"`
	CustomVariableKeys  []string              `json:"custom_variable_keys"`
	SpecialVariableKeys []string              `json:"special_variable_keys"`
	PromptRules         *SimpleChatPromptConfig
}

type PromptTemplateEntity struct {
	PromptType                       string                                  `json:"prompt_type"`
	SimplePromptTemplate             string                                  `json:"simple_prompt_template"`
	AdvancedChatPromptTemplate       *AdvancedChatPromptTemplateEntity       `json:"advanced_chat_prompt_template"`
	AdvancedCompletionPromptTemplate *AdvancedCompletionPromptTemplateEntity `json:"advanced_completion_prompt_template"`
}

type AdvancedCompletionPromptTemplateEntity struct {
	Prompt     string            `json:"prompt"`
	RolePrefix *RolePrefixEntity `json:"role_prefix"`
}

type AdvancedChatMessageEntity struct {
	Text string `json:"text"`
	Role string `json:"role"` // Assuming PromptMessageRole is defined as string
}

type AdvancedChatPromptTemplateEntity struct {
	Messages []*AdvancedChatMessageEntity `json:"messages"`
}
