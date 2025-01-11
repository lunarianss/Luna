// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package speech2text

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/openai_api_compatible/speech2text"
	provider_register "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"

	biz_entity_openai_standard_response "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/openai_standard_response"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
)

type groqAudioLargeLanguageModel struct {
	speech2text.IOpenAudioApiCompactLargeLanguage
}

func init() {
	NewGroqLargeLanguageModel().Register()
}

func NewGroqLargeLanguageModel() *groqAudioLargeLanguageModel {
	return &groqAudioLargeLanguageModel{}
}

var _ provider_register.IAudioModelRegistry = (*groqAudioLargeLanguageModel)(nil)

func (m *groqAudioLargeLanguageModel) Invoke(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, user, filename string, fileContent []byte, modelRuntime biz_entity.IAIModelRuntime) (*biz_entity_openai_standard_response.Speech2TextResp, error) {
	credentials = m.addCustomParameters(credentials)
	m.IOpenAudioApiCompactLargeLanguage = speech2text.NewOpenAudioApiCompactLargeLanguage(fileContent, nil, credentials, model, filename, modelRuntime)
	return m.IOpenAudioApiCompactLargeLanguage.Invoke(ctx)
}

func (m *groqAudioLargeLanguageModel) Register() {
	provider_register.AudioModelRuntimeRegistry.RegisterLargeModelInstance(m)
}

func (m *groqAudioLargeLanguageModel) RegisterName() string {
	return "groq/speech2text"
}

func (m *groqAudioLargeLanguageModel) addCustomParameters(credentials map[string]interface{}) map[string]interface{} {
	credentials["mode"] = "chat"
	credentials["endpoint_url"] = "https://api.groq.com/openai/v1"
	return credentials
}
