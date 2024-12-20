package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
)

type AnnotationRepo interface {
	// Create
	CreateMessageAnnotation(ctx context.Context, annotation *po_entity.MessageAnnotation) (*po_entity.MessageAnnotation, error)
	// CreateAppAnnotationSetting(ctx context.Context, setting *po_entity.AppAnnotationSetting) (*po_entity.AppAnnotationSetting, error)
	// Update
	UpdateMessageAnnotation(ctx context.Context, annotation *po_entity.MessageAnnotation) error
	// Get
	GetMessageAnnotation(ctx context.Context, messageID string) (*po_entity.MessageAnnotation, error)
	GetAnnotationSetting(ctx context.Context, appID string) (*po_entity.AppAnnotationSetting, error)
}
