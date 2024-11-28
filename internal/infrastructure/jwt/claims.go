// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
