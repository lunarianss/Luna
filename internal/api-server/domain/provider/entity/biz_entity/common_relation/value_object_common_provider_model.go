// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

import (
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type ModelType string

const (
	LLM            ModelType = "llm"
	TEXT_EMBEDDING ModelType = "text-embedding"
	RERANK         ModelType = "rerank"
	SPEECH2TEXT    ModelType = "speech2text"
	MODERATION     ModelType = "moderation"
	TTS            ModelType = "tts"
	TEXT2IMG       ModelType = "text2img"
)

var ModelTypeEnums = []string{"llm", "text-embedding", "rerank", "speech2text", "moderation", "tts", "text2img"}

func (m ModelType) ToOriginModelType() (string, error) {
	var originType string

	switch m {
	case LLM:
		originType = "text-generation"
	case TEXT_EMBEDDING:
		originType = "embeddings"
	case RERANK:
		originType = "reranking"
	case SPEECH2TEXT:
		originType = "speech2text"
	case TTS:
		originType = "tts"
	case MODERATION:
		originType = "moderation"
	case TEXT2IMG:
		originType = "text2img"
	default:
		return "", errors.WithCode(code.ErrToOriginModelType, "unknown model type %s", m)
	}
	return originType, nil
}

type ModelFeature string

const (
	TOOL_CALL        ModelFeature = "tool-call"
	MULTI_TOOL_CALL  ModelFeature = "multi-tool-call"
	AGENT_THOUGHT    ModelFeature = "agent-thought"
	VISION           ModelFeature = "vision"
	STREAM_TOOL_CALL ModelFeature = "stream-tool-call"
)

type FetchFrom string

const (
	PREDEFINED_MODEL_FROM   FetchFrom = "predefined-model"
	CUSTOMIZABLE_MODEL_FROM FetchFrom = "customizable-model"
)

type ModelPropertyKey string

const (
	MODE                      ModelPropertyKey = "mode"
	CONTEXT_SIZE              ModelPropertyKey = "context_size"
	MAX_CHUNKS                ModelPropertyKey = "max_chunks"
	FILE_UPLOAD_LIMIT         ModelPropertyKey = "file_upload_limit"
	SUPPORTED_FILE_EXTENSIONS ModelPropertyKey = "supported_file_extensions"
	MAX_CHARACTERS_PER_CHUNK  ModelPropertyKey = "max_characters_per_chunk"
	DEFAULT_VOICE             ModelPropertyKey = "default_voice"
	VOICES                    ModelPropertyKey = "voices"
	WORD_LIMIT                ModelPropertyKey = "word_limit"
	AUDIO_TYPE                ModelPropertyKey = "audio_type"
	MAX_WORKERS               ModelPropertyKey = "max_workers"
)
