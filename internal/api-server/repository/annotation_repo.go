package repo_impl

import (
	"context"

	"github.com/lunarianss/Luna/infrastructure/errors"
	po_account "github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
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

func (ap *AnnotationRepoImpl) GetMessageAnnotation(ctx context.Context, messageID string) (*biz_entity.BizMessageAnnotation, error) {
	var (
		ma      po_entity.MessageAnnotation
		account po_account.Account
	)
	if err := ap.db.First(&ma, "message_id = ?", messageID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get annotation by messageID-[%s] not exists", messageID)
	}

	if err := ap.db.First(&account, "id = ?", ma.AccountID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get account by accountID-[%s] not exists", ma.AccountID)
	}

	return biz_entity.ConvertToBizMessageAnnotation(&ma, &account), nil
}

func (ap *AnnotationRepoImpl) CreateMessageAnnotation(ctx context.Context, annotation *po_entity.MessageAnnotation) (*biz_entity.BizMessageAnnotation, error) {
	var (
		account po_account.Account
	)
	if err := ap.db.Create(annotation).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	if err := ap.db.First(&account, "id = ?", annotation.AccountID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get account by accountID-[%s] not exists", annotation.AccountID)
	}

	return biz_entity.ConvertToBizMessageAnnotation(annotation, &account), nil
}

func (ap *AnnotationRepoImpl) UpdateMessageAnnotation(ctx context.Context, id, answer, question string) error {
	if err := ap.db.Model(&po_entity.MessageAnnotation{}).Select("content", "question").Where("id = ?", id).Updates(map[string]string{"content": answer, "question": question}).Error; err != nil {
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
