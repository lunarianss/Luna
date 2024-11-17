package dto

type SendEmailCodeRequest struct {
	Email    string `json:"email" validate:"required"`
	Language string `json:"language" validate:"required"`
}

type SendEmailCodeResponse struct {
	Data string `json:"data"`
}
