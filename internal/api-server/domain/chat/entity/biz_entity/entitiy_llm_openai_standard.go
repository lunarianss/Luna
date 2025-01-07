package biz_entity

// Define the structure of the response
type OpenaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenaiChoice struct {
	Index        int            `json:"index"`
	Message      *OpenaiMessage `json:"message"`
	Logprobs     string         `json:"logprobs"` // `null` values can be handled with *interface{}
	FinishReason string         `json:"finish_reason"`
}

type OpenaiUsage struct {
	QueueTime        float64 `json:"queue_time"`
	PromptTokens     int     `json:"prompt_tokens"`
	PromptTime       float64 `json:"prompt_time"`
	CompletionTokens int     `json:"completion_tokens"`
	CompletionTime   float64 `json:"completion_time"`
	TotalTokens      int     `json:"total_tokens"`
	TotalTime        float64 `json:"total_time"`
}

type OpenaiResponse struct {
	ID                string          `json:"id"`
	Object            string          `json:"object"`
	Created           int             `json:"created"`
	Model             string          `json:"model"`
	Choices           []*OpenaiChoice `json:"choices"`
	Usage             *OpenaiUsage    `json:"usage"`
	SystemFingerprint string          `json:"system_fingerprint"`
}

type ToolCallFunction struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ToolCall struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Function *ToolCallFunction `json:"function"`
}

type ToolCallStreamFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCallStream struct {
	ID       string                  `json:"id"`
	Type     string                  `json:"type"`
	Function *ToolCallStreamFunction `json:"function"`
}
