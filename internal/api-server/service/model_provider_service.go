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

	domain "github.com/lunarianss/Luna/internal/api-server/domain/model-provider"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/provider"
	"github.com/lunarianss/Luna/internal/api-server/model-runtime/entities"
	model_providers "github.com/lunarianss/Luna/internal/api-server/model-runtime/model-providers"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type ModelProviderService struct {
	ModelProviderDomain *domain.ModelProviderDomain
}

func NewModelProviderService(modelProviderDomain *domain.ModelProviderDomain) *ModelProviderService {
	return &ModelProviderService{ModelProviderDomain: modelProviderDomain}
}

func (mpSrv *ModelProviderService) GetProviderList(ctx context.Context, tenantId string, modelType string) ([]*dto.ProviderResponse, error) {
	var customConfigurationStatus dto.CustomConfigurationStatus

	providerConfigurations, err := mpSrv.ModelProviderDomain.GetConfigurations(ctx, tenantId)
	if err != nil {
		return nil, err
	}

	providerListResponse := make([]*dto.ProviderResponse, 0, model_providers.PROVIDER_COUNT)

	for _, providerConfiguration := range providerConfigurations.Configurations {

		if modelType != "" {
			if !slices.Contains(providerConfiguration.Provider.SupportedModelTypes, entities.ModelType(modelType)) {
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
		}

		if err := providerResponse.PatchIcon(); err != nil {
			return nil, err
		}

		providerListResponse = append(providerListResponse, providerResponse)

	}

	sort.Slice(providerListResponse, func(i, j int) bool {
		return providerListResponse[i].Position < providerListResponse[j].Position
	})

	return providerListResponse, nil
}

func (mpSrv *ModelProviderService) GetProviderIconPath(ctx context.Context, provider, iconType, lang string) (string, error) {

	providerPath, err := mpSrv.ModelProviderDomain.ModelProviderRepo.GetProviderPath(ctx, provider)

	if err != nil {
		return "", err
	}

	providerEntity, err := mpSrv.ModelProviderDomain.ModelProviderRepo.GetProviderEntity(ctx, provider)

	if err != nil {
		return "", err
	}

	iconName, err := mpSrv.getIconName(providerEntity, iconType, lang)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", providerPath, model_providers.ASSETS_DIR, iconName), nil
}

func (mpSrv *ModelProviderService) SaveProviderCredentials(ctx context.Context, tenantId string, provider string, credentials map[string]interface{}) error {
	providerConfigurations, err := mpSrv.ModelProviderDomain.GetConfigurations(ctx, tenantId)

	if err != nil {
		return err
	}

	providerConfiguration, ok := providerConfigurations.Configurations[provider]

	if !ok {
		return errors.WithCode(code.ErrProviderMapModel, fmt.Sprintf("when create %s provider credential for provider", provider))
	}

	if err := mpSrv.ModelProviderDomain.AddOrUpdateCustomProviderCredentials(ctx, providerConfiguration, credentials); err != nil {
		return err
	}
	return nil
}

func (mpSrv *ModelProviderService) getIconName(providerEntity *entities.ProviderEntity, iconType, lang string) (string, error) {
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
