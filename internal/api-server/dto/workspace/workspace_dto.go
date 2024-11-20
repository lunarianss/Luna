package dto

type CurrentTenantInfo struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Plan           string         `json:"plan"`
	Status         string         `json:"status"`
	CreateAt       int64          `json:"create_at"`
	InTrail        bool           `json:"in_trail"`
	TrialEndReason string         `json:"trial_end_reason"`
	Role           string         `json:"role"`
	CustomConfig   map[string]any `json:"custom_config"`
}
