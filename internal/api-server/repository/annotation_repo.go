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

func (ap *AnnotationRepoImpl) CreateAppAnnotationSetting(ctx context.Context, setting *po_entity.AppAnnotationSetting, tx *gorm.DB) (*po_entity.AppAnnotationSetting, error) {
	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ap.db
	}
	if err := dbIns.Create(setting).Error; err != nil {
		return nil, err
	}

	return setting, nil
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

func (ap *AnnotationRepoImpl) FindAppAnnotations(ctx context.Context, appID string) ([]*po_entity.MessageAnnotation, error) {

	var (
		annotations []*po_entity.MessageAnnotation
	)

	if err := ap.db.Where("app_id = ?", appID).Find(&annotations).Error; err != nil {
		return nil, err
	}

	return annotations, nil
}

func (ap *AnnotationRepoImpl) GetAnnotationSettingWithCreate(ctx context.Context, appID string, scoreThreshold float64, bindingID string, accountID string, tx *gorm.DB) (*po_entity.AppAnnotationSetting, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ap.db
	}
	var (
		annotationSetting *po_entity.AppAnnotationSetting
	)

	annotationSetting, err := ap.GetAnnotationSetting(ctx, appID, dbIns)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if annotationSetting.ID == "" {
		setting := &po_entity.AppAnnotationSetting{
			AppID:               appID,
			ScoreThreshold:      scoreThreshold,
			CollectionBindingID: bindingID,
			CreatedUserID:       accountID,
			UpdatedUserID:       accountID,
		}

		annotationSetting, err = ap.CreateAppAnnotationSetting(ctx, setting, dbIns)

		if err != nil {
			return nil, err
		}
	}
	return annotationSetting, nil
}

func (ap *AnnotationRepoImpl) CreateMessageAnnotation(ctx context.Context, annotation *po_entity.MessageAnnotation) (*biz_entity.BizMessageAnnotation, error) {
	var (
		account po_account.Account
	)
	if err := ap.db.Create(annotation).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}

	if err := ap.db.First(&account, "id = ?", annotation.AccountID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get account by accountID-[%s] not exists", annotation.AccountID)
	}

	return biz_entity.ConvertToBizMessageAnnotation(annotation, &account), nil
}

func (ap *AnnotationRepoImpl) UpdateMessageAnnotation(ctx context.Context, id, answer, question string) error {
	if err := ap.db.Model(&po_entity.MessageAnnotation{}).Select("content", "question").Where("id = ?", id).Updates(map[string]string{"content": answer, "question": question}).Error; err != nil {
		return errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (ap *AnnotationRepoImpl) GetAnnotationSetting(ctx context.Context, appID string, tx *gorm.DB) (*po_entity.AppAnnotationSetting, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ap.db
	}

	var ma po_entity.AppAnnotationSetting
	if err := dbIns.First(&ma, "app_id = ?", appID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get annotation setting by appID-[%s] not exists", appID)
	}
	return &ma, nil
}
