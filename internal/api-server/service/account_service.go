package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/config"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	tenantDomain "github.com/lunarianss/Luna/internal/api-server/domain/tenant"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/pkg/util"
	_email "github.com/lunarianss/Luna/pkg/email"
	"github.com/lunarianss/Luna/pkg/log"
)

type AccountService struct {
	AccountDomain *domain.AccountDomain
	TenantDomain  *tenantDomain.TenantDomain
	email         *_email.Mail
	runtimeConfig *config.Config
}

func NewAccountService(accountDomain *domain.AccountDomain, config *config.Config, email *_email.Mail, tenantDomain *tenantDomain.TenantDomain) *AccountService {
	return &AccountService{
		AccountDomain: accountDomain,
		runtimeConfig: config,
		email:         email,
		TenantDomain:  tenantDomain,
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

func (s *AccountService) EmailCodeValidity(ctx context.Context, email, emailCode, token string) (*domain.TokenPair, error) {
	tokenData, err := s.AccountDomain.GetEmailTokenData(ctx, token)

	if err != nil {
		return nil, err
	}

	if err := s.AccountDomain.ValidateAndRevokeData(ctx, email, emailCode, token, tokenData); err != nil {
		return nil, err
	}

	account, err := s.AccountDomain.AccountRepo.GetAccountByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if err := s.AccountDomain.RevokeEmailTokenKey(ctx, token); err != nil {
		return nil, err
	}

	if account == nil {
		if _, err := s.CreateAccountAndTenant(ctx, email, email, "zh-Hans", ""); err != nil {
			return nil, err
		}
	}

	tokenPair, err := s.AccountDomain.Login(ctx, account, util.ExtractRemoteIP(ctx.(*gin.Context)))

	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (ad *AccountService) CreateAccountAndTenant(ctx context.Context, email, name, interfaceLanguage, password string) (*model.Account, error) {

	account, err := ad.AccountDomain.CreateAccount(ctx, email, name, interfaceLanguage, "light", password, false)

	if err != nil {
		return nil, err
	}
	if err := ad.TenantDomain.CreateOwnerTenantIfNotExists(ctx, name, account, false); err != nil {
		return nil, err
	}

	return account, nil
}
