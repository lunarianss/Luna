// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package biz_entity

type EmailTokenData struct {
	Code      string `json:"code"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
}

type ValidateTokenResponse struct {
	Result string     `json:"result"`
	Data   *TokenPair `json:"data"`
}
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
