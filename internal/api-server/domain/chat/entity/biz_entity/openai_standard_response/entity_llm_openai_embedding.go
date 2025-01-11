package biz_entity

type EmbeddingUsage struct {
	Tokens      int
	TotalTokens int
	UnitPrice   float64
	PriceUnit   float64
	TotalPrice  float64
	Currency    string
	Latency     float64
}

type TextEmbeddingResult struct {
	Model      string
	Embeddings [][]float32
	Usage      *EmbeddingUsage
}

type TextEmbeddingLargeModelResult struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
		Object    string    `json:"object"`
	} `json:"data"`
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	ID string `json:"id"`
}
