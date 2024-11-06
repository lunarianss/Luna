package dto

// Create App Input Dto
type CreateAppRequest struct {
	Name           string `json:"name" validate:"required"`
	Mode           string `json:"mode" validate:"required"`
	Icon           string `json:"icon" validate:"required"`
	Description    string `json:"description"`
	IconType       string `json:"icon_type"`
	IconBackground string `json:"icon_background"`
	ApiRph         int    `json:"api_rph"`
	ApiRpm         int    `json:"api_rpm"`
}
