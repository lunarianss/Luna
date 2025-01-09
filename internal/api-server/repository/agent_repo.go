package repo_impl

import (
	"context"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"gorm.io/gorm"
)

type AgentRepoImpl struct {
	db *gorm.DB
}

var _ repository.AgentRepo = (*AgentRepoImpl)(nil)

func NewAgentRepoImpl(db *gorm.DB) *AgentRepoImpl {
	return &AgentRepoImpl{db: db}
}

func (ar *AgentRepoImpl) CreateAgentThought(ctx context.Context, agentThought *po_entity.MessageAgentThought) (*po_entity.MessageAgentThought, error) {
	if err := ar.db.Create(agentThought).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return agentThought, nil
}

func (ar *AgentRepoImpl) CreateToolFile(ctx context.Context, toolFile *po_entity.ToolFile) (*po_entity.ToolFile, error) {
	if err := ar.db.Create(toolFile).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return toolFile, nil
}

func (ar *AgentRepoImpl) CreateMessageFile(ctx context.Context, messageFile *po_entity.MessageFile) (*po_entity.MessageFile, error) {
	if err := ar.db.Create(messageFile).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return messageFile, nil
}

func (ar *AgentRepoImpl) GetMessageFileByID(ctx context.Context, fileID string) (*po_entity.MessageFile, error) {
	var messageFile po_entity.MessageFile

	if err := ar.db.First(&messageFile, "id = ?", fileID).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return &messageFile, nil
}

func (ar *AgentRepoImpl) GetToolFileByID(ctx context.Context, toolFileID string) (*po_entity.ToolFile, error) {
	var toolFile po_entity.ToolFile

	if err := ar.db.First(&toolFile, "id = ?", toolFileID).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return &toolFile, nil
}

func (ar *AgentRepoImpl) GetAgentThoughtByID(ctx context.Context, id string) (*po_entity.MessageAgentThought, error) {
	var agentThought po_entity.MessageAgentThought
	if err := ar.db.First(&agentThought, "id = ?", id).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return &agentThought, nil
}

func (ar *AgentRepoImpl) UpdateAgentThought(ctx context.Context, agentThought *po_entity.MessageAgentThought) error {

	if err := ar.db.Model(&po_entity.MessageAgentThought{}).Where("id = ?", agentThought.ID).Updates(agentThought).Error; err != nil {
		return errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return nil
}
