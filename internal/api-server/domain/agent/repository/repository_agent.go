package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/po_entity"
)

type AgentRepo interface {
	CreateAgentThought(ctx context.Context, agentThought *po_entity.MessageAgentThought) (*po_entity.MessageAgentThought, error)
	CreateToolFile(ctx context.Context, agentThought *po_entity.ToolFile) (*po_entity.ToolFile, error)
	CreateMessageFile(ctx context.Context, agentThought *po_entity.MessageFile) (*po_entity.MessageFile, error)
	UpdateAgentThought(ctx context.Context, agentThought *po_entity.MessageAgentThought) error
	GetAgentThoughtByID(ctx context.Context, id string) (*po_entity.MessageAgentThought, error)
}
