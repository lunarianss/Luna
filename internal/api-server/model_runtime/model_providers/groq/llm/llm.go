// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llm

import "context"

type GroqLargeLanguageModel struct {
}

func (m *GroqLargeLanguageModel) Invoke(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, stop []string, stream bool, user string) error {

	credentials = m.addCustomParameters(credentials)
}

func (m *GroqLargeLanguageModel) addCustomParameters(credentials map[string]interface{}) map[string]interface{} {
	credentials["mode"] = "chat"
	credentials["endpoint_url"] = "https://api.groq.com/openai/v1"
	return credentials
}
