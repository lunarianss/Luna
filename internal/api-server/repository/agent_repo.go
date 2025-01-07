package repo_impl

import (
	"context"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"gorm.io/gorm"
)

type AgentRepoImpl struct {
	db *gorm.DB
}

func NewAgentRepoImpl(db *gorm.DB) *AgentRepoImpl {
	return &AgentRepoImpl{db: db}
}

func (ar *AgentRepoImpl) CreateAgentThought(ctx context.Context, agentThought *po_entity.MessageAgentThought) (*po_entity.MessageAgentThought, error) {
	if err := ar.db.Create(agentThought).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return agentThought, nil
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
