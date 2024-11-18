package service

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/internal/api-server/config"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/pkg/code"
	_email "github.com/lunarianss/Luna/pkg/email"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
)

type AccountService struct {
	AccountDomain *domain.AccountDomain
	email         *_email.Mail
	runtimeConfig *config.Config
}

func NewAccountService(accountDomain *domain.AccountDomain, config *config.Config, email *_email.Mail) *AccountService {
	return &AccountService{
		AccountDomain: accountDomain,
		runtimeConfig: config,
		email:         email,
	}
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

	tokenUUID, emailCode, err := s.AccountDomain.SendEmailCodeLoginEmail(ctx, params.Email, language)
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

func (s *AccountService) EmailCodeValidity(ctx context.Context, email, emailCode, token, tokenType string) (*domain.TokenPair, error) {
	tokenData, err := s.AccountDomain.GetEmailTokenData(ctx, token)

	if err != nil {
		return nil, err
	}

	if tokenData.Code != emailCode {
		return nil, errors.WithCode(code.ErrEmailCode, fmt.Sprintf("email %s, code %s is not valid", email, emailCode))
	}

	if tokenData.Email != email {
		return nil, errors.WithCode(code.ErrEmailCode, "")
	}

	if err := s.AccountDomain.RevokeEmailTokenKey(ctx, token); err != nil {
		return nil, err
	}

	return nil, nil
}
