// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package entities

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
)

type StreamEvent string

const (
	StreamEventPing                   StreamEvent = "ping"
	StreamEventError                  StreamEvent = "error"
	StreamEventMessage                StreamEvent = "message"
	StreamEventMessageEnd             StreamEvent = "message_end"
	StreamEventTTSMessage             StreamEvent = "tts_message"
	StreamEventTTSMessageEnd          StreamEvent = "tts_message_end"
	StreamEventMessageFile            StreamEvent = "message_file"
	StreamEventMessageReplace         StreamEvent = "message_replace"
	StreamEventAgentThought           StreamEvent = "agent_thought"
	StreamEventAgentMessage           StreamEvent = "agent_message"
	StreamEventWorkflowStarted        StreamEvent = "workflow_started"
	StreamEventWorkflowFinished       StreamEvent = "workflow_finished"
	StreamEventNodeStarted            StreamEvent = "node_started"
	StreamEventNodeFinished           StreamEvent = "node_finished"
	StreamEventParallelBranchStarted  StreamEvent = "parallel_branch_started"
	StreamEventParallelBranchFinished StreamEvent = "parallel_branch_finished"
	StreamEventIterationStarted       StreamEvent = "iteration_started"
	StreamEventIterationNext          StreamEvent = "iteration_next"
	StreamEventIterationCompleted     StreamEvent = "iteration_completed"
	StreamEventTextChunk              StreamEvent = "text_chunk"
	StreamEventTextReplace            StreamEvent = "text_replace"
)

type ChatAppTaskState struct {
	Metadata  any
	LLMResult *biz_entity.LLMResult
}

type IStreamResponse interface {
	GetEvent() StreamEvent
	GetTaskID() string
}

// Base StreamResponse struct
type StreamResponse struct {
	Event  StreamEvent `json:"event"`
	TaskID string      `json:"task_id"`
}

func (s *StreamResponse) GetEvent() StreamEvent {
	return s.Event
}

func (s *StreamResponse) GetTaskID() string {
	return s.TaskID
}

// ErrorStreamResponse entity
type ErrorStreamResponse struct {
	*StreamResponse
	Err     string `json:"err"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
}

// MessageStreamResponse entity
type MessageStreamResponse struct {
	*StreamResponse
	ID                   string   `json:"id"`
	Answer               string   `json:"answer"`
	FromVariableSelector []string `json:"from_variable_selector,omitempty"`
}

// MessageAudioStreamResponse entity
type MessageAudioStreamResponse struct {
	*StreamResponse
	Audio string `json:"audio"`
}

// MessageAudioEndStreamResponse entity
type MessageAudioEndStreamResponse struct {
	*StreamResponse
	Audio string `json:"audio"`
}

// MessageEndStreamResponse entity
type MessageEndStreamResponse struct {
	*StreamResponse
	ID       string                   `json:"id"`
	Metadata map[string]interface{}   `json:"metadata"`
	Files    []map[string]interface{} `json:"files"`
}

// MessageFileStreamResponse entity
type MessageFileStreamResponse struct {
	*StreamResponse
	ID        string `json:"id"`
	Type      string `json:"type"`
	BelongsTo string `json:"belongs_to"`
	URL       string `json:"url"`
}

// MessageReplaceStreamResponse entity
type MessageReplaceStreamResponse struct {
	*StreamResponse
	Answer string `json:"answer"`
}

type ChatBotAppEndStreamResponse struct {
	*MessageEndStreamResponse
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
}

type ChatBotAppErrStreamResponse struct {
	*ErrorStreamResponse
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
}

func NewChatBotAppErrStreamResponse(cID, mID string, createAt int64, streamResp *ErrorStreamResponse) *ChatBotAppErrStreamResponse {
	return &ChatBotAppErrStreamResponse{
		ConversationID:      cID,
		MessageID:           mID,
		CreatedAt:           createAt,
		ErrorStreamResponse: streamResp,
	}
}

func NewChatBotAppEndStreamResponse(cID, mID string, createAt int64, streamResp *MessageEndStreamResponse) *ChatBotAppEndStreamResponse {
	return &ChatBotAppEndStreamResponse{
		ConversationID:           cID,
		MessageID:                mID,
		CreatedAt:                createAt,
		MessageEndStreamResponse: streamResp,
	}
}

type ChatBotAppStreamResponse struct {
	*MessageStreamResponse
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
}

func NewChatBotAppStreamResponse(cID, mID string, createAt int64, streamResp *MessageStreamResponse) *ChatBotAppStreamResponse {
	return &ChatBotAppStreamResponse{
		ConversationID:        cID,
		MessageID:             mID,
		CreatedAt:             createAt,
		MessageStreamResponse: streamResp,
	}
}
