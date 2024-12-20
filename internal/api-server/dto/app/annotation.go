package dto

type ApplyAnnotationRequestUrl struct {
	Action string `json:"action" uri:"action"`
	AppID  string `json:"appID" uri:"appID"`
}

type ApplyAnnotationRequestBody struct {
	ScoreThreshold        float32 `json:"score_threshold"`
	EmbeddingProviderName string  `json:"embedding_provider_name"`
	EmbeddingModelName    string  `json:"embedding_model_name"`
}

type ApplyAnnotationResponse struct {
	JobID     string `json:"job_id"`
	JobStatus string `json:"job_status"`
}

func NewApplyAnnotationProcessing(jobID string) *ApplyAnnotationResponse {
	return &ApplyAnnotationResponse{
		JobID:     jobID,
		JobStatus: "processing",
	}
}

func NewApplyAnnotationWaiting(jobID string) *ApplyAnnotationResponse {
	return &ApplyAnnotationResponse{
		JobID:     jobID,
		JobStatus: "waiting",
	}
}
