package domain_service

import (
	"github.com/lunarianss/Luna/internal/api-server/_domain/chat/repository"
)

type ChatDomain struct {
	MessageRepo repository.MessageRepo
}

func NewChatDomain(messageRepo repository.MessageRepo) *ChatDomain {
	return &ChatDomain{
		MessageRepo: messageRepo,
	}
}
