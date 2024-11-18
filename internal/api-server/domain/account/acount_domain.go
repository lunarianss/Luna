package account

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/api-server/repo"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	EMAIL_CODE_TOKEN = "email_code_token"
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
