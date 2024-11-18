package service

import (
	"context"

	"github.com/gin-gonic/gin"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account"
	tenantDomain "github.com/lunarianss/Luna/internal/api-server/domain/tenant"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"gorm.io/gorm"
)

type AccountService struct {
	AccountDomain *domain.AccountDomain
	TenantDomain  *tenantDomain.TenantDomain
	db            *gorm.DB
}

func NewAccountService(accountDomain *domain.AccountDomain, tenantDomain *tenantDomain.TenantDomain, db *gorm.DB) *AccountService {
	return &AccountService{
		AccountDomain: accountDomain,
		TenantDomain:  tenantDomain,
		db:            db,
	}
}

func (s *AccountService) SendEmailCode(ctx context.Context, params *dto.SendEmailCodeRequest) (*dto.SendEmailCodeResponse, error) {

	var (
		language string
	)

	if params.Language == "zh-Hans" {
		language = "zh-Hans"
	} else {
		language = "en-US"
	}

	tokenUUID, emailCode, err := s.AccountDomain.SendEmailCodeLoginEmail(ctx, params.Email, language)

	if err != nil {
		return nil, err
	}

	go s.AccountDomain.SendEmailHtml(ctx, language, params.Email, emailCode)

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

	if account.ID == "" {
		if account, err = s.CreateAccountAndTenant(ctx, email, email, "zh-Hans", ""); err != nil {
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

	tx := ad.db.Begin()
	account, err := ad.AccountDomain.CreateAccountTx(ctx, tx, email, name, interfaceLanguage, "light", password, false)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = ad.TenantDomain.CreateOwnerTenantIfNotExistsTx(ctx, tx, name, account, false)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return account, tx.Commit().Error
}

func (ad *AccountService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {

	return ad.AccountDomain.RefreshToken(ctx, refreshToken)
}
