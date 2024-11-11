// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llm

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime/model_providers/openai_api_compatible/llm"
)

type GroqLargeLanguageModel struct {
	*llm.OpenApiCompactLargeLanguageModel
}

func (m *GroqLargeLanguageModel) Invoke(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, stream bool, user string, promptMessages []*message.PromptMessage) error {
	credentials = m.addCustomParameters(credentials)
	m.OpenApiCompactLargeLanguageModel.Invoke(ctx, promptMessages, modelParameters, credentials)
	return nil
}

func (m *GroqLargeLanguageModel) addCustomParameters(credentials map[string]interface{}) map[string]interface{} {
	credentials["mode"] = "chat"
	credentials["endpoint_url"] = "https://api.groq.com/openai/v1"
	return credentials
}
