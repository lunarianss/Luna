// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_provider

import (
	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
)

type ProviderConfiguration struct {
	TenantId              string               `json:"tenant_id"`
	Provider              *ProviderEntity      `json:"provider"`
	PreferredProviderType model.ProviderType   `json:"preferred_provider_type"`
	UsingProviderType     model.ProviderType   `json:"using_provider_type"`
	SystemConfiguration   *SystemConfiguration `json:"system_configuration"`
	CustomConfiguration   *CustomConfiguration `json:"custom_configuration"`
	ModelSettings         []*ModelSettings     `json:"model_settings"`
}

func (c *ProviderConfiguration) GetCurrentCredentials(modelType base.ModelType, model string) (map[string]interface{}, error) {

	var credentials map[string]interface{}

	if c.CustomConfiguration.Models != nil {
		for _, modelConfiguration := range c.CustomConfiguration.Models {
			if modelConfiguration.ModelType == string(modelType) && modelConfiguration.Model == model {
				credentials = modelConfiguration.Credentials
				break
			}
		}
	}

	if credentials == nil {
		credentials, _ = c.CustomConfiguration.Provider.Credentials.(map[string]interface{})
	}
	return credentials, nil

}

type ProviderConfigurations struct {
	TenantId       string                            `json:"tenant_id"`
	Configurations map[string]*ProviderConfiguration `json:"configurations"`
}
