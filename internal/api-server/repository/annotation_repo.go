package repo_impl

import (
	"context"
	"fmt"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	po_account "github.com/lunarianss/Luna/internal/api-server/domain/account/entity/po_entity"
	biz_entity_app_chat_annotation "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/chat_annotation"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/repository"
	po_dataset "github.com/lunarianss/Luna/internal/api-server/domain/dataset/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/mysql"
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

func (ap *AnnotationRepoImpl) GetMessageAnnotation(ctx context.Context, messageID string) (*biz_entity_app_chat_annotation.BizMessageAnnotation, error) {
	var (
		ma      po_entity.MessageAnnotation
		account po_account.Account
	)
	if err := ap.db.First(&ma, "message_id = ?", messageID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get annotation by messageID-[%s] error: %s", messageID, err.Error())
	}

	if err := ap.db.First(&account, "id = ?", ma.AccountID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get account by accountID-[%s] error: %s", ma.AccountID, err.Error())
	}

	return biz_entity_app_chat_annotation.ConvertToBizMessageAnnotation(&ma, &account), nil
}

func (ap *AnnotationRepoImpl) GetMessageAnnotationHistory(ctx context.Context, messageID string) (*po_entity.AppAnnotationHitHistory, error) {
	var (
		ma po_entity.AppAnnotationHitHistory
	)

	if err := ap.db.First(&ma, "message_id = ?", messageID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get annotation by message-[%s] error: %s", messageID, err.Error())
	}
	return &ma, nil
}

func (ap *AnnotationRepoImpl) GetAnnotationByID(ctx context.Context, id string) (*po_entity.MessageAnnotation, error) {
	var (
		ma po_entity.MessageAnnotation
	)

	if err := ap.db.First(&ma, "id = ?", id).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get annotation by id-[%s] error: %s", id, err.Error())
	}

	return &ma, nil
}

func (ap *AnnotationRepoImpl) FindAppHitAnnotationsInLog(ctx context.Context, appID string, annotationID string, page, pageSize int) ([]*po_entity.AppAnnotationHitHistory, int64, error) {

	var (
		hitAnnotations []*po_entity.AppAnnotationHitHistory
		count          int64
	)

	db := ap.db.Model(&po_entity.AppAnnotationHitHistory{}).Where("app_id = ? AND annotation_id = ?", appID, annotationID).Count(&count).Scopes(mysql.Paginate(page, pageSize))

	if err := db.Order("id DESC").Find(&hitAnnotations).Error; err != nil {
		return nil, count, errors.WithSCode(code.ErrDatabase, err.Error())
	}

	return hitAnnotations, count, nil
}

func (ap *AnnotationRepoImpl) FindAppAnnotationsInLog(ctx context.Context, appID string, page, pageSize int, keyword string) ([]*po_entity.MessageAnnotation, int64, error) {
	var (
		annotations []*po_entity.MessageAnnotation
		count       int64
	)

	db := ap.db.Model(&po_entity.MessageAnnotation{}).Scopes(mysql.Paginate(page, pageSize)).Where("app_id = ?", appID)

	if keyword != "" {
		keyWordFilter := fmt.Sprintf("%%%s%%", keyword)
		db = db.Where("question LIKE ? OR content LIKE ?", keyWordFilter, keyWordFilter)
	}

	if err := db.Order("id DESC").Count(&count).Find(&annotations).Error; err != nil {
		return nil, count, errors.WithSCode(code.ErrDatabase, err.Error())
	}

	return annotations, count, nil

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

func (ap *AnnotationRepoImpl) GetAnnotationSettingWithCreate(ctx context.Context, appID string, scoreThreshold float32, bindingID string, accountID string, tx *gorm.DB) (*po_entity.AppAnnotationSetting, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ap.db
	}
	var (
		annotationSetting   *biz_entity_app_chat_annotation.AnnotationSettingWithBinding
		poAnnotationSetting *po_entity.AppAnnotationSetting
	)

	annotationSetting, err := ap.GetAnnotationSetting(ctx, appID, dbIns)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if annotationSetting == nil {
		setting := &po_entity.AppAnnotationSetting{
			AppID:               appID,
			ScoreThreshold:      scoreThreshold,
			CollectionBindingID: bindingID,
			CreatedUserID:       accountID,
			UpdatedUserID:       accountID,
		}

		poAnnotationSetting, err = ap.CreateAppAnnotationSetting(ctx, setting, dbIns)

		if err != nil {
			return nil, err
		}
	}

	return poAnnotationSetting, nil
}

func (ap *AnnotationRepoImpl) CreateMessageAnnotation(ctx context.Context, annotation *po_entity.MessageAnnotation) (*biz_entity_app_chat_annotation.BizMessageAnnotation, error) {
	var (
		account po_account.Account
	)
	if err := ap.db.Create(annotation).Error; err != nil {
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}

	if err := ap.db.First(&account, "id = ?", annotation.AccountID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get account by accountID-[%s] error: %s", annotation.AccountID, err.Error())
	}

	return biz_entity_app_chat_annotation.ConvertToBizMessageAnnotation(annotation, &account), nil
}

func (ap *AnnotationRepoImpl) UpdateMessageAnnotation(ctx context.Context, id, answer, question string) error {
	log.Info(id, answer, question)
	if err := ap.db.Model(&po_entity.MessageAnnotation{}).Select("content", "question").Where("id = ?", id).Updates(map[string]interface{}{"content": answer, "question": question}).Error; err != nil {
		return errors.WithSCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (ap *AnnotationRepoImpl) GetAnnotationSetting(ctx context.Context, appID string, tx *gorm.DB) (*biz_entity_app_chat_annotation.AnnotationSettingWithBinding, error) {

	var dbIns *gorm.DB

	if tx != nil {
		dbIns = tx
	} else {
		dbIns = ap.db
	}

	var ma po_entity.AppAnnotationSetting
	if err := dbIns.First(&ma, "app_id = ?", appID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get annotation setting by appID-[%s] error: %s", appID, err.Error())
	}

	var binding po_dataset.DatasetCollectionBinding

	if err := dbIns.First(&binding, "id = ?", ma.CollectionBindingID).Error; err != nil {
		return nil, errors.WrapC(err, code.ErrDatabase, "Get collection binding by ID-[%s] not error: %s", ma.CollectionBindingID, err.Error())
	}

	return biz_entity_app_chat_annotation.ConvertPoAnnotationSetting(&ma, &binding), nil
}

func (ap *AnnotationRepoImpl) CreateAppAnnotationHistory(ctx context.Context, history *po_entity.AppAnnotationHitHistory) (*po_entity.AppAnnotationHitHistory, error) {

	tx := ap.db.Begin()
	if err := tx.Model(&po_entity.MessageAnnotation{}).
		Where("id = ?", history.AnnotationID).
		UpdateColumn("hit_count", gorm.Expr("hit_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return nil, errors.WithSCode(code.ErrDatabase, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return history, nil
}
