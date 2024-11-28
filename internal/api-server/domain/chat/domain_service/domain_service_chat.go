// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package domain_service

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/repository"
)

type ChatDomain struct {
	MessageRepo repository.MessageRepo
}

func NewChatDomain(messageRepo repository.MessageRepo) *ChatDomain {
	return &ChatDomain{
		MessageRepo: messageRepo,
	}
}
