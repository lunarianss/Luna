// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model_providers

import (
	// llm
	// groq/llm
	_ "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/groq/llm"
	// tongyi/llm
	_ "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/tongyi/llm"
	// zhipuai/llm
	_ "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/zhipuai/llm"

	// speech2text
	// groq/speech2text
	_ "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/groq/speech2text"
	// tenant/speech2text
	_ "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/tencent/speech2text"

	// tts
	// tongyi/tts
	_ "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/tongyi/tts"

	// embedding
	// tongyi/embedding
	_ "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/tongyi/text_embedding"
)
