package event

import "github.com/lunarianss/Luna/internal/infrastructure/server"

func init() {
	server.RegisterConsumer(&AuthEvent{})
}
