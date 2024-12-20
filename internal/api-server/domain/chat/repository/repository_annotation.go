package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
)

type AnnotationRepo interface {
	// Create
	CreateMessageAnnotation(ctx context.Context, annotation *po_entity.MessageAnnotation) (*biz_entity.BizMessageAnnotation, error)
	// CreateAppAnnotationSetting(ctx context.Context, setting *po_entity.AppAnnotationSetting) (*po_entity.AppAnnotationSetting, error)
	// Update
	UpdateMessageAnnotation(ctx context.Context, id, answer, question string) error
	// Get
	GetMessageAnnotation(ctx context.Context, messageID string) (*biz_entity.BizMessageAnnotation, error)
	GetAnnotationSetting(ctx context.Context, appID string) (*po_entity.AppAnnotationSetting, error)
}
