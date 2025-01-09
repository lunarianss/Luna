package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/po_entity"
)

type AgentRepo interface {
	CreateAgentThought(ctx context.Context, agentThought *po_entity.MessageAgentThought) (*po_entity.MessageAgentThought, error)
	CreateToolFile(ctx context.Context, agentThought *po_entity.ToolFile) (*po_entity.ToolFile, error)
	CreateMessageFile(ctx context.Context, agentThought *po_entity.MessageFile) (*po_entity.MessageFile, error)
	GetMessageFileByID(ctx context.Context, fileID string) (*po_entity.MessageFile, error)
	GetToolFileByID(ctx context.Context, toolFileID string) (*po_entity.ToolFile, error)
	UpdateAgentThought(ctx context.Context, agentThought *po_entity.MessageAgentThought) error
	GetAgentThoughtByID(ctx context.Context, id string) (*po_entity.MessageAgentThought, error)
}
