package token_buffer_memory

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	biz_entity_chat_prompt_message "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/chat_prompt_message"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

const UUID_NIL = "00000000-0000-0000-0000-000000000000"

type ITokenBufferMemory interface {
	GetHistoryPromptMessage(ctx context.Context, maxTokenLimit int, messageLimit int) ([]*biz_entity_chat_prompt_message.PromptMessage, error)
}

type tokenBufferMemory struct {
	conversation        *po_entity.Conversation
	modelRegistryCaller model_registry.IModelRegistryCall
	chatDomain          *domain_service.ChatDomain
}

func NewTokenBufferMemory(conversation *po_entity.Conversation, modelRegistryCaller model_registry.IModelRegistryCall, chatDomain *domain_service.ChatDomain) *tokenBufferMemory {
	return &tokenBufferMemory{
		conversation:        conversation,
		modelRegistryCaller: modelRegistryCaller,
		chatDomain:          chatDomain,
	}
}

func (s *tokenBufferMemory) GetHistoryPromptMessage(ctx context.Context, maxTokenLimit int, messageLimit int) ([]*biz_entity_chat_prompt_message.PromptMessage, error) {

	var (
		promptMessages []*biz_entity_chat_prompt_message.PromptMessage
	)

	if messageLimit != 0 && messageLimit > 0 {
		if messageLimit > 500 {
			messageLimit = 500
		}
	} else {
		messageLimit = 500
	}

	messages, err := s.chatDomain.MessageRepo.FindHistoryPromptMessage(ctx, s.conversation.ID, messageLimit)

	if err != nil {
		return nil, err
	}

	messages = s.extractThreadMessage(messages)

	if messages[0].Answer == "" {
		messages = messages[1:]
	}

	messages = util.SliceReverse(messages)

	for _, message := range messages {
		promptMessages = append(promptMessages, biz_entity_chat_prompt_message.NewUserMessage(message.Query), biz_entity_chat_prompt_message.NewAssistantMessage(message.Answer))
	}

	return promptMessages, nil
}

func (s *tokenBufferMemory) extractThreadMessage(messages []*po_entity.Message) []*po_entity.Message {

	var threadMessages []*po_entity.Message

	var nextMessage *string

	for _, threadMessage := range messages {

		if threadMessage.ParentMessageID == "" {
			threadMessages = append(threadMessages, threadMessage)
			break
		}

		if nextMessage == nil {
			threadMessages = append(threadMessages, threadMessage)
			nextMessage = &threadMessage.ParentMessageID
		} else {
			if *nextMessage == threadMessage.ID || *nextMessage == UUID_NIL {
				threadMessages = append(threadMessages, threadMessage)
				nextMessage = &threadMessage.ParentMessageID
			}
		}
	}

	return threadMessages
}
