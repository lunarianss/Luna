package chat

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

// ConvertToListMessageDto converts a Message from po_entity to ListChatMessageItem.
func ConvertToListMessageDto(message *po_entity.Message) *dto.ListChatMessageItem {
	return &dto.ListChatMessageItem{
		ID:                      message.ID,
		ConversationID:          message.ConversationID,
		Inputs:                  message.Inputs,
		Query:                   message.Query,
		Message:                 ConvertPromptMessageDto(message.Message), // 假设这两个类型兼容
		MessageTokens:           message.MessageTokens,
		MessageUnitPrice:        message.MessageUnitPrice,
		Answer:                  message.Answer,
		AnswerTokens:            message.AnswerTokens,
		ProviderResponseLatency: message.ProviderResponseLatency,
		TotalPrice:              message.TotalPrice,
		Currency:                message.Currency,
		FromSource:              message.FromSource,
		FromEndUserID:           message.FromEndUserID,
		FromAccountID:           message.FromAccountID,
		CreatedAt:               message.CreatedAt,
		MessagePriceUnit:        message.MessagePriceUnit,
		AnswerPriceUnit:         message.AnswerPriceUnit,
		WorkflowRunID:           message.WorkflowRunID,
		Status:                  message.Status,
		Error:                   message.Error,
		MessageMetadata:         message.MessageMetadata,
		InvokeFrom:              message.InvokeFrom,
		ParentMessageID:         message.ParentMessageID,
		FeedBacks:               make([]string, 0),
		AgentThoughts:           make([]string, 0),
		MessageFiles:            make([]string, 0),
		Metadata:                make(map[string]interface{}),
	}
}

func ConvertPromptMessageDto(messages []*po_entity.PromptMessage) []*dto.PromptMessage {
	var pms []*dto.PromptMessage

	if len(messages) == 0 {
		return make([]*dto.PromptMessage, 0)
	}

	for _, msg := range messages {
		pms = append(pms, &dto.PromptMessage{
			Role:    string(msg.Role),
			Content: msg.Content.(string),
			Name:    msg.Name,
		})
	}

	return pms
}

func ConvertToConversationJoins(conversation *po_entity.Conversation) *dto.ListChatConversationItem {
	return &dto.ListChatConversationItem{
		ID:            conversation.ID,
		Status:        conversation.Status,
		FromSource:    conversation.FromSource,
		FromEndUserID: conversation.FromEndUserID,
		FromAccountID: conversation.FromAccountID,
		Name:          conversation.Name,
		Summary:       conversation.Summary,
		ReadAt:        conversation.ReadAt,
		CreatedAt:     conversation.CreatedAt,
		UpdatedAt:     conversation.UpdatedAt,
	}
}
