package assembler

import dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"

func ConvertToCreateChatMessageBody(msg *dto.ServiceCreateChatMessageBody) *dto.CreateChatMessageBody {
	return &dto.CreateChatMessageBody{
		ResponseMode:                 msg.ResponseMode,
		ConversationID:               msg.ConversationID,
		Query:                        msg.Query,
		Files:                        msg.Files,
		Inputs:                       msg.Inputs,
		ModelConfig:                  msg.ModelConfig,
		ParentMessageId:              msg.ParentMessageId,
		AutoGenerateConversationName: msg.AutoGenerateConversationName,
	}
}
