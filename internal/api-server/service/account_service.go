package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"github.com/redis/go-redis/v9"
)

const (
	EMAIL_CODE_TOKEN = "email_code_token"
)

type AccountService struct {
	AccountDomain *domain.AccountDomain
	redis         *redis.Client
}

func NewAccountService(accountDomain *domain.AccountDomain, redis *redis.Client) *AccountService {
	return &AccountService{
		AccountDomain: accountDomain,
		redis:         redis,
	}
}

type EmailTokenData struct {
	Code      string `json:"code"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
}

func (s *AccountService) GetUserThroughEmails(ctx context.Context, email string) (*model.Account, error) {
	return s.AccountDomain.AccountRepo.GetAccountByEmail(ctx, email)
}

func (s *AccountService) SetEmailCode(ctx context.Context, params *dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error) {

	var (
		language string
	)

	if params.Language == "zh-Hans" {
		language = "zh-Hans"
	}

	language = "en-US"

	account, err := s.AccountDomain.AccountRepo.GetAccountByEmail(ctx, params.Email)

	if account.ID != "" {

	}

	if err != nil {
		return nil, err
	}

	tokenUUID, err := s.SendEmailCodeLoginEmail(ctx, params.Email, language)
	if err != nil {
		return nil, err
	}

	go func() {
		
	}()

	return &dto.SendEmailCodeResponse{
		Data: tokenUUID,
	}, nil
}

func (s *AccountService) GetEmailCodeToken(email string, tokenType string, code string) (string, string, error) {

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

func (s *AccountService) GetEmailTokenKey(token, tokenType string) string {
	return fmt.Sprintf("%s:token:%s", tokenType, token)
}

func (s *AccountService) SendEmailCodeLoginEmail(ctx context.Context, email string, language string) (string, error) {

	code := util.GenerateRandomNumber()

	tokenData, tokenUUID, err := s.GetEmailCodeToken(email, EMAIL_CODE_TOKEN, code)

	if err != nil {
		return "", err
	}

	tokenKey := s.GetEmailTokenKey(tokenUUID, EMAIL_CODE_TOKEN)

	if err := s.redis.Set(ctx, tokenKey, tokenData, 5*time.Minute).Err(); err != nil {
		return "", err
	}

	return tokenUUID, redis.Nil
}
