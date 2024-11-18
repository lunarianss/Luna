package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/config"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/pkg/util"
	_email "github.com/lunarianss/Luna/pkg/email"
	"github.com/lunarianss/Luna/pkg/log"
	"github.com/redis/go-redis/v9"
)

const (
	EMAIL_CODE_TOKEN = "email_code_token"
)

type AccountService struct {
	AccountDomain *domain.AccountDomain
	redis         *redis.Client
	email         *_email.Mail
	runtimeConfig *config.Config
}

func NewAccountService(accountDomain *domain.AccountDomain, redis *redis.Client, config *config.Config, email *_email.Mail) *AccountService {
	return &AccountService{
		AccountDomain: accountDomain,
		redis:         redis,
		runtimeConfig: config,
		email:         email,
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

func (s *AccountService) SendEmailCode(ctx context.Context, params *dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error) {

	var (
		language     string
		templatePath string
	)

	if params.Language == "zh-Hans" {
		language = "zh-Hans"
	} else {
		language = "en-US"
	}

	// account, err := s.AccountDomain.AccountRepo.GetAccountByEmail(ctx, params.Email)

	// if err != nil {
	// 	return nil, err
	// }

	tokenUUID, emailCode, err := s.SendEmailCodeLoginEmail(ctx, params.Email, language)
	if err != nil {
		return nil, err
	}
	if language == "en-US" {
		templatePath = fmt.Sprintf("%s/%s", s.runtimeConfig.EmailOptions.TemplateDir, "email_code_login_mail_template_en-US.html")
	} else {
		templatePath = fmt.Sprintf("%s/%s", s.runtimeConfig.EmailOptions.TemplateDir, "email_code_login_mail_template_zh-CN.html")
	}

	go func() {
		err := s.email.Send(params.Email, "Email Code", templatePath, map[string]interface{}{
			"Code": emailCode,
		}, "")

		if err != nil {
			log.Errorf("Send email failed: %v", err)
		}
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

func (s *AccountService) SendEmailCodeLoginEmail(ctx context.Context, email string, language string) (string, string, error) {

	code := util.GenerateRandomNumber()

	tokenData, tokenUUID, err := s.GetEmailCodeToken(email, EMAIL_CODE_TOKEN, code)

	if err != nil {
		return "", "", err
	}

	tokenKey := s.GetEmailTokenKey(tokenUUID, EMAIL_CODE_TOKEN)

	if err := s.redis.Set(ctx, tokenKey, tokenData, 5*time.Minute).Err(); err != nil {
		return "", "", err
	}

	return tokenUUID, code, nil
}
