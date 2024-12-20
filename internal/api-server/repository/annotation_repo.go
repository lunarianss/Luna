package repo_impl

import (
	"context"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/repository"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"gorm.io/gorm"
)

type AnnotationRepoImpl struct {
	db *gorm.DB
}

var _ repository.AnnotationRepo = (*AnnotationRepoImpl)(nil)

func NewAnnotationRepoImpl(db *gorm.DB) *AnnotationRepoImpl {
	return &AnnotationRepoImpl{db}
}

func (ap *AnnotationRepoImpl) GetMessageAnnotation(ctx context.Context, messageID string) (*po_entity.MessageAnnotation, error) {
	var ma po_entity.MessageAnnotation
	if err := ap.db.First(&ma, "message_id = ?", messageID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get annotation by messageID-[%s] not exists", messageID)
	}
	return &ma, nil
}

func (ap *AnnotationRepoImpl) CreateMessageAnnotation(ctx context.Context, annotation *po_entity.MessageAnnotation) (*po_entity.MessageAnnotation, error) {
	if err := ap.db.Create(annotation).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return annotation, nil
}

func (ap *AnnotationRepoImpl) UpdateMessageAnnotation(ctx context.Context, annotation *po_entity.MessageAnnotation) error {
	if err := ap.db.Model(&po_entity.MessageAnnotation{}).Select("content", "question").Where("id = ?", annotation.ID).Updates(annotation).Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (ap *AnnotationRepoImpl) GetAnnotationSetting(ctx context.Context, appID string) (*po_entity.AppAnnotationSetting, error) {
	var ma po_entity.AppAnnotationSetting
	if err := ap.db.First(&ma, "app_id = ?", appID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get annotation setting by appID-[%s] not exists", appID)
	}
	return &ma, nil
}
