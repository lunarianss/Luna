package repository

import (
	"context"

	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"gorm.io/gorm"
)

type AnnotationRepo interface {
	// Create
	CreateMessageAnnotation(ctx context.Context, annotation *po_entity.MessageAnnotation) (*biz_entity.BizMessageAnnotation, error)
	CreateAppAnnotationSetting(ctx context.Context, setting *po_entity.AppAnnotationSetting, tx *gorm.DB) (*po_entity.AppAnnotationSetting, error)
	CreateAppAnnotationHistory(ctx context.Context, history *po_entity.AppAnnotationHitHistory) (*po_entity.AppAnnotationHitHistory, error)
	// Update
	UpdateMessageAnnotation(ctx context.Context, id, answer, question string) error
	// Get
	GetMessageAnnotation(ctx context.Context, messageID string) (*biz_entity.BizMessageAnnotation, error)
	GetMessageAnnotationHistory(ctx context.Context, messageID string) (*po_entity.AppAnnotationHitHistory, error)
	GetAnnotationByID(ctx context.Context, id string) (*po_entity.MessageAnnotation, error)
	GetAnnotationSetting(ctx context.Context, appID string, tx *gorm.DB) (*biz_entity.AnnotationSettingWithBinding, error)
	GetAnnotationSettingWithCreate(ctx context.Context, appID string, scoreThreshold float32, bindingID string, accountID string, tx *gorm.DB) (*po_entity.AppAnnotationSetting, error)
	// Find
	FindAppAnnotations(ctx context.Context, appID string) ([]*po_entity.MessageAnnotation, error)
	FindAppAnnotationsInLog(ctx context.Context, appID string, page, pageSize int, keyword string) ([]*po_entity.MessageAnnotation, int64, error)
	FindAppHitAnnotationsInLog(ctx context.Context, appID string, annotationID string, page, pageSize int) ([]*po_entity.AppAnnotationHitHistory, int64, error)
}
