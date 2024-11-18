package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type LunaClaims struct {
	AccountId   int32
	NickName    string
	AuthorityId int32
	jwt.RegisteredClaims
}

var JWTIns *JWT

type JWT struct {
	SignKey []byte
}

func NewJWT(signKey string) *JWT {
	JWTIns = &JWT{
		SignKey: []byte(signKey),
	}
	return JWTIns
}

func GetJWTIns() *JWT {
	return JWTIns
}

func (j *JWT) GenerateJWT(claims LunaClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(j.SignKey)
	if err != nil || tokenStr == "" {
		return "", errors.WithCode(code.ErrTokenGenerate, fmt.Sprintf("use claims %+v generate token with sign key %s", claims, j.SignKey))
	}
	return tokenStr, nil
}

func (j *JWT) ParseJWT(tokenStr string) (*LunaClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &LunaClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.WithCode(code.ErrTokenMethodErr, fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		return j.SignKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.WithCode(code.ErrTokenInvalid, err.Error())
	}

	if claims, ok := token.Claims.(*LunaClaims); ok {
		return claims, nil
	} else {
		return nil, errors.WithCode(code.ErrTokenInvalid, fmt.Sprintf("token %s can not be parse as a LunaClaims", tokenStr))
	}
}

// func (j *JWT) RefreshJWT(tokenStr string) (string, error) {
// 	jwt.WithTimeFunc(func() time.Time {
// 		return time.Unix(0, 0)
// 	})

// 	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return j.SignKey, nil
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
// 		jwt.WithTimeFunc(func() time.Time {
// 			return time.Now()
// 		})
// 		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
// 		return j.GenerateJWT(*claims)
// 	}

// 	return "", TokenInvalid

// }