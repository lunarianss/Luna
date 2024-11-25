package domain

import "github.com/lunarianss/Luna/internal/api-server/repo"

type ChatDomain struct {
	MessageRepo repo.MessageRepo
}

func NewChatDomain(messageRepo repo.MessageRepo) *ChatDomain {
	return &ChatDomain{
		MessageRepo: messageRepo,
	}
}
