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

func (m *groqLargeLanguageModel) Invoke(ctx context.Context, queueManager *biz_entity_chat.StreamGenerateQueue, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, stream bool, user string, promptMessages []*po_entity.PromptMessage, modelRuntime biz_entity.IAIModelRuntime) {
	credentials = m.addCustomParameters(credentials)
	m.IOpenApiCompactLargeLanguage = llm.NewOpenApiCompactLargeLanguageModel(promptMessages, modelParameters, credentials, queueManager, model, stream, modelRuntime)
	m.IOpenApiCompactLargeLanguage.Invoke(ctx)
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
