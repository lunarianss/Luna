// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llm

import (
	"context"

	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers/openai_api_compatible/llm"
	provider_register "github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
)

type tongyiLargeLanguageModel struct {
	llm.IOpenApiCompactLargeLanguage
}

func init() {
	NewTongyiLargeLanguageModel().Register()
}

func NewTongyiLargeLanguageModel() *tongyiLargeLanguageModel {
	return &tongyiLargeLanguageModel{}
}

var _ provider_register.IModelRegistry = (*tongyiLargeLanguageModel)(nil)

func (m *tongyiLargeLanguageModel) Invoke(ctx context.Context, queueManager *biz_entity_chat.StreamGenerateQueue, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, user string, promptMessages []*po_entity.PromptMessage, modelRuntime biz_entity.IAIModelRuntime) {
	credentials = m.addCustomParameters(credentials)
	m.IOpenApiCompactLargeLanguage = llm.NewOpenApiCompactLargeLanguageModel(promptMessages, modelParameters, credentials, model, modelRuntime)
	m.IOpenApiCompactLargeLanguage.Invoke(ctx, queueManager)
}

func (m *tongyiLargeLanguageModel) InvokeNonStream(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, user string, promptMessages []*po_entity.PromptMessage, modelRuntime biz_entity.IAIModelRuntime) (*biz_entity_chat.LLMResult, error) {
	credentials = m.addCustomParameters(credentials)
	m.IOpenApiCompactLargeLanguage = llm.NewOpenApiCompactLargeLanguageModel(promptMessages, modelParameters, credentials, model, modelRuntime)
	return m.IOpenApiCompactLargeLanguage.InvokeNonStream(ctx)
}

func (m *tongyiLargeLanguageModel) Register() {
	provider_register.ModelRuntimeRegistry.RegisterLargeModelInstance(m)
}

func (m *tongyiLargeLanguageModel) RegisterName() string {
	return "tongyi/llm"
}

func (m *tongyiLargeLanguageModel) addCustomParameters(credentials map[string]interface{}) map[string]interface{} {
	credentials["mode"] = "chat"
	credentials["endpoint_url"] = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	credentials["api_key"] = credentials["dashscope_api_key"]
	return credentials
}
