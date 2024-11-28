package biz_entity

type EmailTokenData struct {
	Code      string `json:"code"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
