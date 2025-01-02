package biz_entity

type AgentEntity struct {
	Provider     string
	Model        string
	Strategy     Strategy
	Prompt       *AgentPromptEntity
	Tools        []*AgentToolEntity
	MaxIteration int
}

type AgentChatAppConfig struct {
	*EasyUIBasedAppConfig
	*AgentEntity
}
