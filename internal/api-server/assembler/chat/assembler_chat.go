package assembler

import (
	po_account "github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
)

func ConvertToAnnotation(annotation *biz_entity.BizMessageAnnotation) *dto.MessageAnnotation {
	if annotation == nil {
		return nil
	}

	return &dto.MessageAnnotation{
		ID:             annotation.ID,
		AppID:          annotation.AppID,
		ConversationID: annotation.ConversationID,
		MessageID:      annotation.MessageID,
		Question:       annotation.Question,
		Content:        annotation.Content,
		HitCount:       annotation.HitCount,
		AccountID:      annotation.AccountID,
		CreatedAt:      annotation.CreatedAt,
		UpdatedAt:      annotation.UpdatedAt,
		Account:        ConvertToAnnotationAccount(annotation.Account),
	}
}

func ConvertToAnnotationAccount(a *po_account.Account) (s *dto.AnnotationAccount) {
	return &dto.AnnotationAccount{
		ID:                a.ID,
		Name:              a.Name,
		Email:             a.Email,
		Avatar:            a.Avatar,
		InterfaceLanguage: a.InterfaceLanguage,
		InterfaceTheme:    a.InterfaceTheme,
		Timezone:          a.Timezone,
		LastLoginIP:       a.LastLoginIP,
		LastLoginAt:       a.LastLoginAt,
		CreatedAt:         a.CreatedAt,
	}
}

// ConvertToListMessageDto converts a Message from po_entity to ListChatMessageItem.
func ConvertToListMessageDto(message *po_entity.Message, annotation *biz_entity.BizMessageAnnotation, history *po_entity.AppAnnotationHitHistory, account *po_account.Account) *dto.ListChatMessageItem {
	messageItem := &dto.ListChatMessageItem{
		ID:                      message.ID,
		ConversationID:          message.ConversationID,
		Inputs:                  message.Inputs,
		Query:                   message.Query,
		Message:                 ConvertPromptMessageDto(message.Message),
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
		Annotation:              ConvertToAnnotation(annotation),
	}

	if account != nil && history != nil {
		messageItem.AnnotationHistory = &dto.AnnotationHistory{
			Annotation: history.AnnotationID,
			AnnotationCreateAccount: &dto.AnnotationCreateAccount{
				ID:    account.ID,
				Name:  account.Name,
				Email: account.Email,
			},
			CreatedAt: history.CreatedAt,
		}
	}

	return messageItem
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

func CovertToServiceChatCompletionResponse(message *po_entity.Message, cID string, llmResult *biz_entity.LLMResult) *dto.ServiceChatCompletionResponse {
	return &dto.ServiceChatCompletionResponse{
		MessageID:      message.ID,
		CreatedAt:      message.CreatedAt,
		ConversationID: cID,
		Mode:           "chat",
		Answer:         llmResult.Message.Content.(string),
		Metadata: &dto.ServiceChatCompletionMetaDataResponse{
			RetrieverResources: make([]any, 0),
			Usage: &dto.Usage{
				PromptTokens:        llmResult.Usage.PromptTokens,
				PromptUnitPrice:     llmResult.Usage.PromptUnitPrice,
				PromptPriceUnit:     llmResult.Usage.PromptPriceUnit,
				PromptPrice:         llmResult.Usage.PromptPrice,
				CompletionTokens:    llmResult.Usage.CompletionTokens,
				CompletionUnitPrice: llmResult.Usage.CompletionUnitPrice,
				CompletionPriceUnit: llmResult.Usage.PromptPriceUnit,
				CompletionPrice:     llmResult.Usage.CompletionPrice,
				TotalTokens:         llmResult.Usage.TotalTokens,
				Currency:            llmResult.Usage.Currency,
				Latency:             llmResult.Usage.Latency,
			},
		},
	}
}
