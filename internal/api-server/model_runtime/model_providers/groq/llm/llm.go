// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llm

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers/openai_api_compatible/llm"
	provider_register "github.com/lunarianss/Luna/internal/api-server/model_runtime/model_registry"
)

type GroqLargeLanguageModel struct {
	*llm.OpenApiCompactLargeLanguageModel
}

func init() {
	NewGroqLargeLanguageModel().Register()
}

func NewGroqLargeLanguageModel() *GroqLargeLanguageModel {
	return &GroqLargeLanguageModel{}
}

var _ provider_register.IModelRegistry = (*GroqLargeLanguageModel)(nil)

func (m *GroqLargeLanguageModel) Invoke(ctx context.Context, queueManager *model_runtime.StreamGenerateQueue, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, stream bool, user string, promptMessages []*message.PromptMessage) {
	credentials = m.addCustomParameters(credentials)
	m.OpenApiCompactLargeLanguageModel = &llm.OpenApiCompactLargeLanguageModel{
		Stream: stream,
		Model:  model,
	}

	m.OpenApiCompactLargeLanguageModel.Invoke(ctx, promptMessages, modelParameters, credentials, queueManager)
}

func (m *GroqLargeLanguageModel) Register() {

	provider_register.ModelRuntimeRegistry.RegisterLargeModelInstance(m)
}
func (m *GroqLargeLanguageModel) RegisterName() string {
	return "groq/llm"
}

func (m *GroqLargeLanguageModel) addCustomParameters(credentials map[string]interface{}) map[string]interface{} {
	credentials["mode"] = "chat"
	credentials["endpoint_url"] = "https://api.groq.com/openai/v1"
	return credentials
}
