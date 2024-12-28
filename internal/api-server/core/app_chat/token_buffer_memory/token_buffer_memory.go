package token_buffer_memory

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

const UUID_NIL = "00000000-0000-0000-0000-000000000000"

type ITokenBufferMemory interface {
	GetHistoryPromptMessage(ctx context.Context, maxTokenLimit int, messageLimit int) ([]*po_entity.PromptMessage, error)
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

func (s *tokenBufferMemory) GetHistoryPromptMessage(ctx context.Context, maxTokenLimit int, messageLimit int) ([]*po_entity.PromptMessage, error) {

	var (
		promptMessages []*po_entity.PromptMessage
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
		promptMessages = append(promptMessages, po_entity.NewUserMessage(message.Query), po_entity.NewAssistantMessage(message.Answer))
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
