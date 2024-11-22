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
	_email "github.com/lunarianss/Luna/pkg/email"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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
	TenantRepo  repo.TenantRepo
	redis       *redis.Client
	config      *config.Config
	email       *_email.Mail
}

func NewAccountDomain(accountRepo repo.AccountRepo, redis *redis.Client, config *config.Config, email *_email.Mail, tenantRepo repo.TenantRepo) *AccountDomain {
	return &AccountDomain{
		AccountRepo: accountRepo,
		redis:       redis,
		config:      config,
		email:       email,
		TenantRepo:  tenantRepo,
	}
}

func (ad *AccountDomain) GetUserThroughEmails(ctx context.Context, email string) (*model.Account, error) {
	return ad.AccountRepo.GetAccountByEmail(ctx, email)
}

func (ad *AccountDomain) GetEmailTokenData(ctx context.Context, token string) (*EmailTokenData, error) {

	var (
		tokenData EmailTokenData
	)

	tokenKey := ad.GetEmailTokenKey(token, EMAIL_CODE_TOKEN)

	v, err := ad.redis.Get(ctx, tokenKey).Result()

	if errors.Is(err, redis.Nil) {
		return nil, errors.WithCode(code.ErrRedisDataExpire, fmt.Sprintf("redis key %s not found", tokenKey))
	} else if err != nil {
		return nil, errors.WithCode(code.ErrRedisRuntime, err.Error())
	}

	if err := json.Unmarshal([]byte(v), &tokenData); err != nil {
		return nil, errors.WithCode(code.ErrDecodingJSON, err.Error())
	}
	return &tokenData, nil
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
func (ad *AccountDomain) SendEmailHtml(ctx context.Context, language string, email, emailCode string) {
	var templatePath string

	if language == "en-US" {
		templatePath = fmt.Sprintf("%s/%s", ad.config.EmailOptions.TemplateDir, "email_code_login_mail_template_en-US.html")
	} else {
		templatePath = fmt.Sprintf("%s/%s", ad.config.EmailOptions.TemplateDir, "email_code_login_mail_template_zh-CN.html")
	}

	go func() {
		err := ad.email.Send(email, "Email Code", templatePath, map[string]interface{}{
			"Code": emailCode,
		}, "")

		if err != nil {
			log.Errorf("Send email failed: %v", err)
		}
	}()
}

func (ad *AccountDomain) Login(ctx context.Context, account *model.Account, ipAddress string) (*TokenPair, error) {

	if ipAddress != "" {
		account.LastLoginIP = ipAddress
		now := time.Now().UTC().Unix()
		account.LastLoginAt = &now
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
			ExpiresAt: jwtV5.NewNumericDate(time.Now().Add(ad.config.JwtOptions.Timeout)),
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
func (ad *AccountDomain) CreateAccount(ctx context.Context, tx *gorm.DB, email, name, interfaceLanguage, interfaceTheme, password string, isSetup bool) (*model.Account, error) {
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

	return ad.AccountRepo.CreateAccount(ctx, account, tx)
}

func (ad *AccountDomain) DeleteRefreshToken(ctx context.Context, refreshToken string, accountID string) error {

	if err := ad.redis.Del(ctx, ad.GetRefreshTokenKey(refreshToken)).Err(); err != nil {
		return err
	}

	if err := ad.redis.Del(ctx, ad.GetAccountRefreshTokenKey(accountID)).Err(); err != nil {
		return err
	}

	return nil
}

func (ad *AccountDomain) StoreRefreshToken(ctx context.Context, refreshToken string, accountID string) error {
	if err := ad.redis.Set(ctx, ad.GetRefreshTokenKey(refreshToken), accountID, ad.config.JwtOptions.Refresh).Err(); err != nil {
		return err
	}

	if err := ad.redis.Set(ctx, ad.GetAccountRefreshTokenKey(accountID), refreshToken, ad.config.JwtOptions.Refresh).Err(); err != nil {
		return err
	}

	return nil
}

func (ad *AccountDomain) LoadUser(ctx context.Context, userID string) (*model.Account, error) {
	account, err := ad.AccountRepo.GetAccountByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	if account.Status == string(model.BANNED) {
		return nil, errors.WithCode(code.ErrAccountBanned, fmt.Sprintf("account %s, email %s, id %s is already banned", account.Name, account.Email, account.ID))
	}

	tenantJoin, err := ad.TenantRepo.GetCurrentTenantJoinByAccount(ctx, account)

	if err != nil {
		return nil, err
	}

	if tenantJoin.ID == "" {
		tenantJoin, err := ad.TenantRepo.FindTenantJoinByAccount(ctx, account, nil)

		if err != nil {
			return nil, err
		}

		if tenantJoin.ID == "" {
			return nil, err
		}

		tenantJoin.Current = 1

		_, err = ad.TenantRepo.UpdateCurrentTenantAccountJoin(ctx, tenantJoin)

		if err != nil {
			return nil, err
		}
	}

	if account.LastLoginAt != nil {
		now := time.Now().UTC().Unix()
		if now-account.LastActiveAt > 10*60 {
			account.LastActiveAt = now
			if err := ad.AccountRepo.UpdateAccountLastActive(ctx, account); err != nil {
				return nil, err
			}
		}
	}

	return account, nil
}

func (ad *AccountDomain) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {

	v, err := ad.redis.Get(ctx, ad.GetRefreshTokenKey(refreshToken)).Result()

	if errors.Is(err, redis.Nil) {
		return nil, errors.WithCode(code.ErrRefreshTokenNotFound, fmt.Sprintf("refresh token %s not found", refreshToken))
	} else if err != nil {
		return nil, errors.WithCode(code.ErrRedisRuntime, err.Error())
	}

	account, err := ad.LoadUser(ctx, v)

	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.WithCode(code.ErrRecordNotFound, fmt.Sprintf("account record %s not found", v))
	}

	newAccessToken, err := ad.GenerateToken(ctx, account)

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := util.GenerateRefreshToken(64)

	if err != nil {
		return nil, err
	}

	if err := ad.DeleteRefreshToken(ctx, refreshToken, account.ID); err != nil {
		return nil, err
	}

	if err := ad.StoreRefreshToken(ctx, newRefreshToken, account.ID); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (ad *AccountDomain) GetCurrentTenantOfAccount(ctx context.Context, accountID string) (*model.Tenant, *model.TenantAccountJoin, error) {

	accountRecord, err := ad.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, nil, err
	}

	accountJoinRecord, err := ad.TenantRepo.GetCurrentTenantJoinByAccount(ctx, accountRecord)

	if err != nil {
		return nil, nil, err
	}

	tenantRecord, err := ad.TenantRepo.GetTenantByID(ctx, accountJoinRecord.TenantID)

	if err != nil {
		return nil, nil, err
	}

	return tenantRecord, accountJoinRecord, nil
}
