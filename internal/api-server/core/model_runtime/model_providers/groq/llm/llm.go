// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llm

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/openai_api_compatible/llm"
	provider_register "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
)

type groqLargeLanguageModel struct {
	llm.IOpenApiCompactLargeLanguage
}

func init() {
	NewGroqLargeLanguageModel().Register()
}

func NewGroqLargeLanguageModel() *groqLargeLanguageModel {
	return &groqLargeLanguageModel{}
}

var _ provider_register.IModelRegistry = (*groqLargeLanguageModel)(nil)

func (m *groqLargeLanguageModel) Invoke(ctx context.Context, queueManager *biz_entity_chat.StreamGenerateQueue, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, user string, promptMessages []*po_entity.PromptMessage, modelRuntime biz_entity.IAIModelRuntime) {
	credentials = m.addCustomParameters(credentials)
	m.IOpenApiCompactLargeLanguage = llm.NewOpenApiCompactLargeLanguageModel(promptMessages, modelParameters, credentials, model, modelRuntime)
	m.IOpenApiCompactLargeLanguage.Invoke(ctx, queueManager)
}

func (m *groqLargeLanguageModel) InvokeNonStream(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, user string, promptMessages []*po_entity.PromptMessage, modelRuntime biz_entity.IAIModelRuntime) (*biz_entity_chat.LLMResult, error) {
	credentials = m.addCustomParameters(credentials)
	m.IOpenApiCompactLargeLanguage = llm.NewOpenApiCompactLargeLanguageModel(promptMessages, modelParameters, credentials, model, modelRuntime)
	return m.IOpenApiCompactLargeLanguage.InvokeNonStream(ctx)
}

func (m *groqLargeLanguageModel) Register() {

	provider_register.ModelRuntimeRegistry.RegisterLargeModelInstance(m)
}
func (m *groqLargeLanguageModel) RegisterName() string {
	return "groq/llm"
}

func (m *groqLargeLanguageModel) addCustomParameters(credentials map[string]interface{}) map[string]interface{} {
	credentials["mode"] = "chat"
	credentials["endpoint_url"] = "https://api.groq.com/openai/v1"
	return credentials
}
