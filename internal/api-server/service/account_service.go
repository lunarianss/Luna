package service

import (
	"context"

	"github.com/gin-gonic/gin"
	domain "github.com/lunarianss/Luna/internal/api-server/_domain/account/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/_domain/account/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/_domain/account/entity/po_entity"
	accountDto "github.com/lunarianss/Luna/internal/api-server/dto/account"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/pkg/util"
	"gorm.io/gorm"
)

type AccountService struct {
	accountDomain *domain.AccountDomain
	tenantDomain  *domain.TenantDomain
	db            *gorm.DB
}

func NewAccountService(accountDomain *domain.AccountDomain, tenantDomain *domain.TenantDomain, db *gorm.DB) *AccountService {
	return &AccountService{
		accountDomain: accountDomain,
		tenantDomain:  tenantDomain,
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

	tokenUUID, emailCode, err := s.accountDomain.SendEmailCodeLoginEmail(ctx, params.Email, language)

	if err != nil {
		return nil, err
	}

	go s.accountDomain.SendEmailHtml(ctx, language, params.Email, emailCode)

	return &dto.SendEmailCodeResponse{
		Data: tokenUUID,
	}, nil
}

func (s *AccountService) EmailCodeValidity(ctx context.Context, email, emailCode, token string) (*biz_entity.TokenPair, error) {
	tokenData, err := s.accountDomain.GetEmailTokenData(ctx, token)

	if err != nil {
		return nil, err
	}

	if err := s.accountDomain.ValidateAndRevokeData(ctx, email, emailCode, token, tokenData); err != nil {
		return nil, err
	}

	account, err := s.accountDomain.AccountRepo.GetAccountByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if account.ID == "" {
		if account, err = s.CreateAccountAndTenant(ctx, email, email, "zh-Hans", ""); err != nil {
			return nil, err
		}
	}

	tokenPair, err := s.accountDomain.Login(ctx, account, util.ExtractRemoteIP(ctx.(*gin.Context)))

	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (ad *AccountService) CreateAccountAndTenant(ctx context.Context, email, name, interfaceLanguage, password string) (*po_entity.Account, error) {

	tx := ad.db.Begin()
	account, err := ad.accountDomain.CreateAccount(ctx, tx, email, name, interfaceLanguage, "light", password, false)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = ad.tenantDomain.CreateOwnerTenantIfNotExists(ctx, tx, account, false)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return account, tx.Commit().Error
}

func (ad *AccountService) RefreshToken(ctx context.Context, refreshToken string) (*biz_entity.TokenPair, error) {

	return ad.accountDomain.RefreshToken(ctx, refreshToken)
}

func (ad *AccountService) GetAccountProfile(ctx context.Context, accountID string) (*accountDto.GetAccountProfileResp, error) {
	accountRecord, err := ad.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	resp := accountDto.AccountConvertToProfile(accountRecord)

	if accountRecord.Password != "" && accountRecord.PasswordSalt != "" {
		resp.IsPasswordSet = true
	}

	return resp, nil
}
