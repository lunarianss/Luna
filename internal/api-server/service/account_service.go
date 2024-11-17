package service

import (
	"context"
	"encoding/json"
	"fmt"

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
	Redis         *redis.Client
}

func NewAccountService(accountDomain *domain.AccountDomain, redis *redis.Client) *AccountService {
	return &AccountService{
		AccountDomain: accountDomain,
		Redis:         redis,
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

	if err != nil {
		return nil, err
	}

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
		return "", nil
	}

	return string(tokenByte), token, nil
}

func (s *AccountService) GetEmailTokenKey(token, tokenType string) string {
	fmt.Sprintf("%s:token:%s", tokenType, token)
}

func (s *AccountService) SendEmailCodeLoginEmail(email string, language string) (string, error) {

	code := util.GenerateRandomNumber()

	tokenData, tokenUUID, err := s.GetEmailCodeToken(email, EMAIL_CODE_TOKEN, code)

	if err != nil {
		return "", err
	}

	tokenKey := s.GetEmailTokenKey(tokenUUID, EMAIL_CODE_TOKEN)

}
