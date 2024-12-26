package dto

type ApplyAnnotationRequestUrl struct {
	Action string `json:"action" uri:"action" validate:"required"`
	AppID  string `json:"appID" uri:"appID" validate:"required"`
}

type ApplyAnnotationStatusRequestUrl struct {
	Action string `json:"action" uri:"action" validate:"required"`
	AppID  string `json:"appID" uri:"appID" validate:"required"`
	JobID  string `json:"jobID" uri:"jobID" validate:"required"`
}

type ApplyAnnotationRequestBody struct {
	ScoreThreshold        float32 `json:"score_threshold" validate:"required"`
	EmbeddingProviderName string  `json:"embedding_provider_name" validate:"required"`
	EmbeddingModelName    string  `json:"embedding_model_name" validate:"required"`
}

type ApplyAnnotationResponse struct {
	JobID     string `json:"job_id"`
	JobStatus string `json:"job_status"`
}

type ApplyAnnotationStatusResponse struct {
	JobID        string `json:"job_id"`
	JobStatus    string `json:"job_status"`
	ErrorMessage string `json:"error_msg"`
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
