// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	domain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/account/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	accountDto "github.com/lunarianss/Luna/internal/api-server/dto/account"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/auth"
	"github.com/lunarianss/Luna/internal/api-server/event"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	"gorm.io/gorm"
)

type AccountService struct {
	accountDomain *domain.AccountDomain
	tenantDomain  *domain.TenantDomain
	db            *gorm.DB
	mqProducer    rocketmq.Producer
}

func NewAccountService(accountDomain *domain.AccountDomain, tenantDomain *domain.TenantDomain, db *gorm.DB, mqProducer rocketmq.Producer) *AccountService {
	return &AccountService{
		accountDomain: accountDomain,
		tenantDomain:  tenantDomain,
		db:            db,
		mqProducer:    mqProducer,
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

	sendCodeMessage := event.SendEmailCodeMessage{
		Language:  language,
		Email:     params.Email,
		EmailCode: emailCode,
	}

	messageBodyData, err := json.Marshal(sendCodeMessage)

	if err != nil {
		return nil, errors.WithCode(code.ErrEncodingJSON, "sendCodeMessage json encoding error %s", err.Error())
	}

	message := &primitive.Message{
		Topic: event.AuthTopic,
		Body:  messageBodyData,
	}

	message = message.WithTag(event.SendEmailCodeTag)

	sendResult, err := s.mqProducer.SendSync(ctx, message)

	if err != nil {
		return nil, errors.WithCode(code.ErrMQSend, "mq send sync error when send email code: %s", err.Error())
	}

	log.Infof("MQ-Send-Result %s", sendResult.String())

	return &dto.SendEmailCodeResponse{
		Data:   tokenUUID,
		Result: "success",
	}, nil
}

func (s *AccountService) EmailCodeValidity(ctx context.Context, email, emailCode, token string) (*biz_entity.ValidateTokenResponse, error) {
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

	return &biz_entity.ValidateTokenResponse{
		Data:   tokenPair,
		Result: "success",
	}, nil
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

func (ad *AccountService) RefreshToken(ctx context.Context, refreshToken string) (*accountDto.RefreshTokenResponse, error) {

	tokenPair, err := ad.accountDomain.RefreshToken(ctx, refreshToken)

	if err != nil {
		return nil, err
	}

	return &accountDto.RefreshTokenResponse{
		Data: &accountDto.TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: refreshToken,
		},
		Result: "success",
	}, nil
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
