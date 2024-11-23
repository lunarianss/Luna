package jwt

import "github.com/golang-jwt/jwt/v5"

type LunaPassportClaims struct {
	AppID     string
	AppCode   string
	EndUserID string
	jwt.RegisteredClaims
}

type LunaClaims struct {
	AccountId   string
	NickName    string
	AuthorityId string
	jwt.RegisteredClaims
}
