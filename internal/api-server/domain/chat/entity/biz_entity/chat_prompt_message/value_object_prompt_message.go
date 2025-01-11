package biz_entity

type PromptMessageRole string

const (
	SYSTEM    PromptMessageRole = "system"
	USER      PromptMessageRole = "user"
	ASSISTANT PromptMessageRole = "assistant"
	TOOL      PromptMessageRole = "tool"
)

type PromptMessageContentType string

const (
	TEXT  PromptMessageContentType = "text"
	IMAGE PromptMessageContentType = "image"
	AUDIO PromptMessageContentType = "audio"
)

type PromptMessageContent struct {
	Type PromptMessageContentType `json:"type"`
	Data any                      `json:"data"`
}
