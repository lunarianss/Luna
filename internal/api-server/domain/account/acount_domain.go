package account

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	jwtV5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/jwt"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	EMAIL_CODE_TOKEN             = "email_code_token"
	REFRESH_TOKEN_PREFIX         = "refresh_token"
	ACCOUNT_REFRESH_TOKEN_PREFIX = "account_refresh_token"
)

type EmailTokenData struct {
	Code      string `json:"code"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccountDomain struct {
	AccountRepo repo.AccountRepo
	redis       *redis.Client
	config      *config.Config
}

func NewAccountDomain(accountRepo repo.AccountRepo, redis *redis.Client) *AccountDomain {
	return &AccountDomain{
		AccountRepo: accountRepo,
		redis:       redis,
	}
}

func (ad *AccountDomain) GetUserThroughEmails(ctx context.Context, email string) (*model.Account, error) {
	return ad.AccountRepo.GetAccountByEmail(ctx, email)
}

func (ad *AccountDomain) GetEmailTokenData(ctx context.Context, token string) (*EmailTokenData, error) {

	var (
		tokenData *EmailTokenData
	)

	tokenKey := ad.GetEmailTokenKey(token, EMAIL_CODE_TOKEN)

	v, err := ad.redis.Get(ctx, tokenKey).Result()

	if errors.Is(err, redis.Nil) {
		return nil, errors.WithCode(code.ErrRedisDataExpire, fmt.Sprintf("redis key %s not found", tokenKey))
	} else if err != nil {
		return nil, errors.WithCode(code.ErrRedisRuntime, err.Error())
	}

	if err := json.Unmarshal([]byte(v), tokenData); err != nil {
		return nil, errors.WithCode(code.ErrDecodingJSON, err.Error())
	}
	return tokenData, nil
}

func (ad *AccountDomain) ValidateAndRevokeData(ctx context.Context, email, emailCode, token string, tokenData *EmailTokenData) error {
	if tokenData.Code != emailCode {
		return errors.WithCode(code.ErrEmailCode, fmt.Sprintf("email %s, code %s is not valid", email, emailCode))
	}

	if tokenData.Email != email {
		return errors.WithCode(code.ErrEmailCode, "")
	}

	if err := ad.RevokeEmailTokenKey(ctx, token); err != nil {
		return err
	}

	return nil
}

func (ad *AccountDomain) ConstructEmailCodeToken(email string, tokenType string, code string) (string, string, error) {

	token := uuid.NewString()
	tokenData := &EmailTokenData{
		Email:     email,
		TokenType: tokenType,
		Code:      code,
	}

	tokenByte, err := json.Marshal(tokenData)

	if err != nil {
		return "", "", nil
	}

	return string(tokenByte), token, nil
}

func (ad *AccountDomain) GetEmailTokenKey(token, tokenType string) string {
	return fmt.Sprintf("%s:token:%s", tokenType, token)
}

func (ad *AccountDomain) GetRefreshTokenKey(refreshToken string) string {
	return fmt.Sprintf("%s:%s", REFRESH_TOKEN_PREFIX, refreshToken)
}

func (ad *AccountDomain) GetAccountRefreshTokenKey(accountID string) string {
	return fmt.Sprintf("%s:%s", ACCOUNT_REFRESH_TOKEN_PREFIX, accountID)
}

func (ad *AccountDomain) RevokeEmailTokenKey(ctx context.Context, token string) error {
	tokenKey := ad.GetEmailTokenKey(token, EMAIL_CODE_TOKEN)
	if err := ad.redis.Del(ctx, tokenKey).Err(); err != nil {
		return errors.WithCode(code.ErrRedisRuntime, err.Error())
	}

	return nil
}

func (ad *AccountDomain) SendEmailCodeLoginEmail(ctx context.Context, email string, language string) (string, string, error) {

	code := util.GenerateRandomNumber()

	tokenData, tokenUUID, err := ad.ConstructEmailCodeToken(email, EMAIL_CODE_TOKEN, code)

	if err != nil {
		return "", "", err
	}

	tokenKey := ad.GetEmailTokenKey(tokenUUID, EMAIL_CODE_TOKEN)

	if err := ad.redis.Set(ctx, tokenKey, tokenData, 5*time.Minute).Err(); err != nil {
		return "", "", err
	}

	return tokenUUID, code, nil
}

func (ad *AccountDomain) Login(ctx context.Context, account *model.Account, ipAddress string) (*TokenPair, error) {

	if ipAddress != "" {
		account.LastLoginIP = ipAddress
		account.LastLoginAt = time.Now().UTC().Unix()
		if err := ad.AccountRepo.UpdateAccountIpAddress(ctx, account); err != nil {
			return nil, err
		}
	}

	if account.Status == string(model.PENDING) {
		account.Status = string(model.ACTIVE)

		if err := ad.AccountRepo.UpdateAccountStatus(ctx, account); err != nil {
			return nil, err
		}
	}

	accessToken, err := ad.GenerateToken(ctx, account)

	if err != nil {
		return nil, err
	}

	refreshToken, err := util.GenerateRefreshToken(64)

	if err != nil {
		return nil, err
	}

	if err := ad.StoreRefreshToken(ctx, refreshToken, account.ID); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (ad *AccountDomain) GenerateToken(ctx context.Context, account *model.Account) (string, error) {

	var (
		jwtToken string
		err      error
	)

	claims := jwt.LunaClaims{
		RegisteredClaims: jwtV5.RegisteredClaims{
			ExpiresAt: jwtV5.NewNumericDate(time.Now().Add(ad.config.JwtOptions.Timeout * time.Hour)),
			IssuedAt:  jwtV5.NewNumericDate(time.Now()),
			NotBefore: jwtV5.NewNumericDate(time.Now()),
			Issuer:    ad.config.JwtOptions.Realm,
			Subject:   "Admin",
			Audience:  []string{"not yet"},
		},
		AccountId: account.ID,
	}

	jwt := jwt.GetJWTIns()

	if jwt == nil {
		return jwtToken, errors.WithCode(code.ErrTokenInsNotFound, "")
	}

	jwtToken, err = jwt.GenerateJWT(claims)

	if err != nil {
		return jwtToken, err
	}
	return jwtToken, nil
}
func (ad *AccountDomain) CreateAccount(ctx context.Context, email, name, interfaceLanguage, interfaceTheme, password string, isSetup bool) (*model.Account, error) {
	// todo 补充密码和 system feature
	var timezone string
	timezone, ok := util.LanguageMapping[interfaceLanguage]

	if !ok {
		timezone = "UTC"
	}

	account := &model.Account{
		Email:             email,
		Name:              name,
		InterfaceLanguage: interfaceLanguage,
		InterfaceTheme:    interfaceTheme,
		Timezone:          timezone,
	}

	return ad.AccountRepo.CreateAccount(ctx, account)
}

func (ad *AccountDomain) StoreRefreshToken(ctx context.Context, refreshToken string, accountID string) error {
	if err := ad.redis.Set(ctx, ad.GetRefreshTokenKey(refreshToken), accountID, ad.config.JwtOptions.Refresh*24*time.Hour).Err(); err != nil {
		return err
	}

	if err := ad.redis.Set(ctx, ad.GetAccountRefreshTokenKey(accountID), refreshToken, ad.config.JwtOptions.Refresh*24*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}
