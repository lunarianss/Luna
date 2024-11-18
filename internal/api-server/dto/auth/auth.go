package dto

type SendEmailCodeRequest struct {
	Email    string `json:"email" validate:"required"`
	Language string `json:"language" validate:"required"`
}

type SendEmailCodeResponse struct {
	Data string `json:"data"`
}

type EmailCodeValidityRequest struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required"`
	Token string `json:"token" validate:"required"`
}