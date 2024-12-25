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
