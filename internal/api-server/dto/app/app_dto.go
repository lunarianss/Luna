package dto

// Create App Input Dto
type CreateAppRequest struct {
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description"`
	Mode           string `json:"mode" validate:"required"`
	IconType       string `json:"icon_type"`
	Icon           string `json:"icon" validate:"required"`
	IconBackground string `json:"icon_background"`
}
