// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	accountDomain "github.com/lunarianss/Luna/internal/api-server/_domain/account/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/_domain/provider/domain_service"
	common "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/common_relation"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/_domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/config"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	model_providers "github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type ModelProviderService struct {
	providerDomain *domain_service.ProviderDomain
	accountDomain  *accountDomain.AccountDomain
	config         *config.Config
}

func NewModelProviderService(providerDomain *domain_service.ProviderDomain, accountDomain *accountDomain.AccountDomain, config *config.Config) *ModelProviderService {
	return &ModelProviderService{providerDomain: providerDomain, accountDomain: accountDomain, config: config}
}

func (mpSrv *ModelProviderService) GetProviderList(ctx context.Context, accountID string, modelType string) ([]*dto.ProviderResponse, error) {
	var customConfigurationStatus dto.CustomConfigurationStatus

	tenantRecord, _, err := mpSrv.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	providerConfigurations, err := mpSrv.providerDomain.GetConfigurations(ctx, tenantRecord.ID)
	if err != nil {
		return nil, err
	}

	providerListResponse := make([]*dto.ProviderResponse, 0, model_providers.PROVIDER_COUNT)

	for _, providerConfiguration := range providerConfigurations.Configurations {

		if modelType != "" {
			if !slices.Contains(providerConfiguration.Provider.SupportedModelTypes, common.ModelType(modelType)) {
				continue
			}
		}

		if providerConfiguration.CustomConfiguration.Provider != nil {
			customConfigurationStatus = dto.ACTIVE
		} else {
			customConfigurationStatus = dto.NO_CONFIGURE
		}

		providerResponse := &dto.ProviderResponse{
			Provider:                 providerConfiguration.Provider.Provider,
			Label:                    providerConfiguration.Provider.Label,
			Description:              providerConfiguration.Provider.Description,
			IconSmall:                providerConfiguration.Provider.IconSmall,
			IconLarge:                providerConfiguration.Provider.IconLarge,
			Background:               providerConfiguration.Provider.Background,
			Help:                     providerConfiguration.Provider.Help,
			Position:                 providerConfiguration.Provider.Position,
			SupportedModelTypes:      providerConfiguration.Provider.SupportedModelTypes,
			ConfigurationMethods:     providerConfiguration.Provider.ConfigurationMethods,
			ProviderCredentialSchema: providerConfiguration.Provider.ProviderCredentialSchema,
			ModelCredentialSchema:    providerConfiguration.Provider.ModelCredentialSchema,
			PreferredProviderType:    providerConfiguration.PreferredProviderType,
			CustomConfiguration: &dto.CustomConfigurationResponse{
				Status: customConfigurationStatus,
			},
			SystemConfiguration: &dto.SystemConfigurationResponse{
				Enabled:             false,
				QuotaConfigurations: make([]*biz_entity_provider_config.QuotaConfiguration, 0),
			},
		}

		providerListResponse = append(providerListResponse, providerResponse)
	}

	for _, providerResponse := range providerListResponse {
		providerResponse.PatchIcon(mpSrv.config)
	}

	sort.Slice(providerListResponse, func(i, j int) bool {
		return providerListResponse[i].Position < providerListResponse[j].Position
	})

	return providerListResponse, nil
}

func (mpSrv *ModelProviderService) GetProviderIconPath(ctx context.Context, provider, iconType, lang string) (string, error) {

	providerPath, err := mpSrv.providerDomain.ProviderRepo.GetProviderPath(ctx, provider)

	if err != nil {
		return "", err
	}

	providerEntity, err := mpSrv.providerDomain.ProviderRepo.GetProviderEntity(ctx, provider)

	if err != nil {
		return "", err
	}

	iconName, err := mpSrv.getIconName(providerEntity, iconType, lang)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", providerPath, model_providers.ASSETS_DIR, iconName), nil
}

func (mpSrv *ModelProviderService) SaveProviderCredentials(ctx context.Context, accountID string, provider string, credentials map[string]interface{}) error {

	tenantRecord, _, err := mpSrv.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return err
	}

	if err := mpSrv.providerDomain.SaveProviderCredentials(ctx, tenantRecord.ID, provider, credentials); err != nil {
		return err
	}

	return nil
}

func (mpSrv *ModelProviderService) getIconName(providerEntity *biz_entity.ProviderStaticConfiguration, iconType, lang string) (string, error) {
	var (
		iconName string
	)

	if iconType == "icon_small" {
		if providerEntity.IconSmall == nil {
			return "", errors.WithCode(code.ErrProviderNotHaveIcon, fmt.Sprintf("provider %s not have a small icon", providerEntity.Provider))
		}

		if strings.ToLower(lang) == "zh_hans" {
			iconName = providerEntity.IconSmall.Zh_Hans
		} else {
			iconName = providerEntity.IconSmall.En_US
		}
	} else {
		if providerEntity.IconLarge == nil {
			return "", errors.WithCode(code.ErrProviderNotHaveIcon, fmt.Sprintf("provider %s not have a large icon", providerEntity.Provider))
		}

		if strings.ToLower(lang) == "zh_hans" {
			iconName = providerEntity.IconLarge.Zh_Hans
		} else {
			iconName = providerEntity.IconLarge.En_US
		}
	}

	return iconName, nil
}
