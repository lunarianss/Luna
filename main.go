// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"

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

func get() (string, error) {
	return "134", nil
}

func DeepCopyUsingJSON[T interface{}](src, dst T) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

type A struct {
	Name string `json:"name"`
}

var a = &A{Name: "test"}

func main() {
	var b A
	// DeepCopyUsingJSON(a, b)
	// a.Name = "cccc"
	fmt.Printf(b == nil)
	fmt.Printf("b %+v", b)

}
