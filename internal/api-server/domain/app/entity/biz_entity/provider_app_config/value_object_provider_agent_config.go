package biz_entity

type Strategy string

var (
	CHAIN_OF_THOUGHT Strategy = "chain-of-thought"
	FUNCTION_CALLING Strategy = "function-calling"
)

type AgentToolProviderType string

var (
	BUILTIN                AgentToolProviderType = "builtin"
	API                    AgentToolProviderType = "api"
	WORKFLOW_PROVIDER_TYPE AgentToolProviderType = "workflow"
)

type AgentPromptEntity struct {
	FirstPrompt   string
	NextIteration string
}

type AgentToolEntity struct {
	ProviderType   AgentToolProviderType
	ProviderID     string
	ToolName       string
	ToolParameters map[string]any
}
