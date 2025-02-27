// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/google/uuid"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	assembler "github.com/lunarianss/Luna/internal/api-server/assembler/chat"
	accountDomain "github.com/lunarianss/Luna/internal/api-server/domain/account/domain_service"
	appDomain "github.com/lunarianss/Luna/internal/api-server/domain/app/domain_service"
	chatDomain "github.com/lunarianss/Luna/internal/api-server/domain/chat/domain_service"
	biz_entity_app_chat_annotation "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/chat_annotation"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	datasetDomain "github.com/lunarianss/Luna/internal/api-server/domain/dataset/domain_service"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	dto_app "github.com/lunarianss/Luna/internal/api-server/dto/app"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/api-server/event"
	"github.com/lunarianss/Luna/internal/api-server/event/event_handler"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AnnotationService struct {
	appDomain      *appDomain.AppDomain
	providerDomain *domain_service.ProviderDomain
	accountDomain  *accountDomain.AccountDomain
	chatDomain     *chatDomain.ChatDomain
	datasetDomain  *datasetDomain.DatasetDomain
	redis          *redis.Client

	mq rocketmq.Producer
}

func NewAnnotationService(appDomain *appDomain.AppDomain, providerDomain *domain_service.ProviderDomain, accountDomain *accountDomain.AccountDomain, chatDomain *chatDomain.ChatDomain, redis *redis.Client, mq rocketmq.Producer, datasetDomain *datasetDomain.DatasetDomain) *AnnotationService {
	return &AnnotationService{
		appDomain:      appDomain,
		providerDomain: providerDomain,
		accountDomain:  accountDomain,
		chatDomain:     chatDomain,
		redis:          redis,
		mq:             mq,
		datasetDomain:  datasetDomain,
	}
}

func (as *AnnotationService) InsertAnnotationFromMessage(ctx context.Context, accountID string, appID string, args *dto.InsertAnnotationFormMessage) (*dto.MessageAnnotation, error) {

	var (
		bizMessageAnnotation *biz_entity_app_chat_annotation.BizMessageAnnotation
	)

	accountRecord, err := as.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, _, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	app, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	if args.MessageID != "" {
		message, err := as.chatDomain.MessageRepo.GetMessageByApp(ctx, args.MessageID, app.ID)

		if err != nil {
			return nil, err
		}

		bizMessageAnnotation, err = as.chatDomain.AnnotationRepo.GetMessageAnnotation(ctx, message.ID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				annotation := &po_entity.MessageAnnotation{
					AppID:          app.ID,
					ConversationID: message.ConversationID,
					MessageID:      message.ID,
					Content:        args.Answer,
					Question:       args.Question,
					AccountID:      accountID,
				}
				bizMessageAnnotation, err = as.chatDomain.AnnotationRepo.CreateMessageAnnotation(ctx, annotation)

				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			bizMessageAnnotation.Content = args.Answer
			bizMessageAnnotation.Question = args.Question
			if err := as.chatDomain.AnnotationRepo.UpdateMessageAnnotation(ctx, bizMessageAnnotation.ID, bizMessageAnnotation.Content, bizMessageAnnotation.Question); err != nil {
				return nil, err
			}
		}
	} else {
		annotation := &po_entity.MessageAnnotation{
			AppID:     appID,
			Content:   args.Answer,
			Question:  args.Question,
			AccountID: accountID,
		}
		bizMessageAnnotation, err = as.chatDomain.AnnotationRepo.CreateMessageAnnotation(ctx, annotation)
		if err != nil {
			return nil, err
		}
	}

	annotationSetting, err := as.chatDomain.AnnotationRepo.GetAnnotationSetting(ctx, app.ID, nil)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assembler.ConvertToAnnotation(bizMessageAnnotation), nil
		} else {
			return nil, err
		}
	} else {
		addAnnotationBody := &event_handler.AddAnnotationTask{
			AnnotationID:        bizMessageAnnotation.ID,
			Question:            bizMessageAnnotation.Question,
			TenantID:            tenant.ID,
			AppID:               appID,
			CollectionBindingID: annotationSetting.CollectionBindingID,
			AccountID:           accountID,
		}

		marshalMessageBody, err := json.Marshal(addAnnotationBody)

		if err != nil {
			return nil, err
		}

		message := &primitive.Message{
			Topic: event.AnnotationTopic,
			Body:  marshalMessageBody,
		}

		message.WithTag(event.AddAnnotationTag)
		sendResult, err := as.mq.SendSync(ctx, message)

		if err != nil {
			return nil, errors.WithCode(code.ErrMQSend, "mq send sync error when send enable app annotation: %s", err.Error())
		}

		log.Infof("MQ-Send-Result %s", sendResult.String())
	}

	return nil, nil
}

func (as *AnnotationService) EnableAppAnnotationStatus(ctx context.Context, appID, accountID, jobID, action string) (*dto_app.ApplyAnnotationStatusResponse, error) {

	accountRecord, err := as.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, "tenant %s don't have the permission to enable app annotation", tenant.Name)
	}

	appEnableJobKey := fmt.Sprintf("%s_app_annotation_job_%s", action, jobID)

	val, err := as.redis.Get(ctx, appEnableJobKey).Result()

	if err != nil {
		return nil, errors.WithSCode(code.ErrNotFoundJobID, err.Error())
	}

	var errorMsg string
	if val == "error" {
		appAnnotationErrorKey := fmt.Sprintf("%s_app_annotation_error_%s", action, jobID)
		errorMsg, err = as.redis.Get(ctx, appAnnotationErrorKey).Result()
		if err != nil {
			errorMsg = "Occurred internal server when get error info"
			log.Errorf("occurred error %s when get key from redis", err.Error())
		}
	}

	return &dto_app.ApplyAnnotationStatusResponse{
		JobID:        jobID,
		JobStatus:    val,
		ErrorMessage: errorMsg,
	}, nil
}

func (as *AnnotationService) GetAnnotationSetting(ctx context.Context, accountID string, appID string) (*dto_app.AnnotationSettingResponse, error) {
	tenant, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, "You don't have the permission for %s", tenant.Name)
	}

	app, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithSCode(code.ErrResourceNotFound, err.Error())
		}
		return nil, err
	}

	annotationSetting, err := as.chatDomain.AnnotationRepo.GetAnnotationSetting(ctx, app.ID, nil)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto_app.AnnotationSettingResponse{
				Enabled: false,
			}, nil
		} else {
			return nil, err
		}
	}

	return &dto_app.AnnotationSettingResponse{
		Enabled:        true,
		ID:             annotationSetting.ID,
		ScoreThreshold: annotationSetting.ScoreThreshold,
		EmbeddingModel: &dto_app.AnnotationSettingEmbeddingModel{
			EmbeddingProviderName: annotationSetting.CollectionBindingDetail.ProviderName,
			EmbeddingModelName:    annotationSetting.CollectionBindingDetail.ModelName,
		},
	}, nil
}

func (as *AnnotationService) ListAnnotations(ctx context.Context, appID, accountID string, args *dto_app.ListAnnotationsArgs) (*dto_app.ListAnnotationResponse, error) {

	var (
		annotationItems []*dto_app.ListAnnotationItem
	)
	accountRecord, err := as.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, "tenant %s don't have the permission to enable app annotation", tenant.Name)
	}

	app, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	annotations, count, err := as.chatDomain.AnnotationRepo.FindAppAnnotationsInLog(ctx, app.ID, args.Page, args.Limit, args.Keyword)

	if err != nil {
		return nil, err
	}

	for _, annotation := range annotations {
		annotationItems = append(annotationItems, &dto_app.ListAnnotationItem{
			ID:        annotation.ID,
			Question:  annotation.Question,
			Answer:    annotation.Content,
			HitCount:  annotation.HitCount,
			CreatedAt: annotation.CreatedAt,
		})
	}

	if len(annotationItems) == 0 {
		annotationItems = make([]*dto_app.ListAnnotationItem, 0)
	}

	return &dto_app.ListAnnotationResponse{
		Data:    annotationItems,
		Page:    args.Page,
		Limit:   args.Limit,
		Total:   count,
		HasMore: len(annotationItems) == args.Limit,
	}, nil
}

func (as *AnnotationService) ListHitAnnotations(ctx context.Context, appID, annotationID, accountID string, args *dto_app.ListHitAnnotationsArgs) (*dto_app.ListHitAnnotationResponse, error) {

	var (
		annotationItems []*dto_app.ListHitAnnotationItem
	)
	accountRecord, err := as.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, "tenant %s don't have the permission to enable app annotation", tenant.Name)
	}

	app, err := as.appDomain.AppRepo.GetTenantApp(ctx, appID, tenant.ID)

	if err != nil {
		return nil, err
	}

	hitAnnotations, count, err := as.chatDomain.AnnotationRepo.FindAppHitAnnotationsInLog(ctx, app.ID, annotationID, args.Page, args.Limit)

	if err != nil {
		return nil, err
	}

	for _, hitAnnotation := range hitAnnotations {

		annotationItems = append(annotationItems, &dto_app.ListHitAnnotationItem{
			ID:        hitAnnotation.ID,
			Source:    hitAnnotation.Source,
			Score:     hitAnnotation.Score,
			Question:  hitAnnotation.Question,
			Match:     hitAnnotation.AnnotationQuestion,
			Response:  hitAnnotation.AnnotationContent,
			CreatedAt: hitAnnotation.CreatedAt,
		})
	}

	if len(annotationItems) == 0 {
		annotationItems = make([]*dto_app.ListHitAnnotationItem, 0)
	}

	return &dto_app.ListHitAnnotationResponse{
		Data:    annotationItems,
		Page:    args.Page,
		Limit:   args.Limit,
		Total:   count,
		HasMore: len(annotationItems) == args.Limit,
	}, nil
}

func (as *AnnotationService) EnableAppAnnotation(ctx context.Context, appID, accountID string, args *dto_app.ApplyAnnotationRequestBody) (*dto_app.ApplyAnnotationResponse, error) {

	accountRecord, err := as.accountDomain.AccountRepo.GetAccountByID(ctx, accountID)

	if err != nil {
		return nil, err
	}

	tenant, tenantJoin, err := as.accountDomain.GetCurrentTenantOfAccount(ctx, accountRecord.ID)

	if err != nil {
		return nil, err
	}

	if !tenantJoin.IsEditor() {
		return nil, errors.WithCode(code.ErrForbidden, "tenant %s don't have the permission to enable app annotation", tenant.Name)
	}

	enableAnnotationKey := fmt.Sprintf("enable_app_annotation_%s", appID)

	v, err := as.redis.Get(ctx, enableAnnotationKey).Result()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, err
		}
	}

	if v != "" {
		return dto_app.NewApplyAnnotationProcessing(v), nil
	}

	jobID := uuid.NewString()

	enableAppAnnotationJobKey := fmt.Sprintf("enable_app_annotation_job_%s", jobID)

	_, err = as.redis.SetNX(ctx, enableAppAnnotationJobKey, "waiting", time.Duration(0)).Result()

	if err != nil {
		return nil, err
	}

	enableAnnotationMessageBody := &event_handler.EnableAnnotationReplyTask{
		JobID:                 jobID,
		AppID:                 appID,
		AccountID:             accountID,
		TenantID:              tenant.ID,
		ScoreThreshold:        args.ScoreThreshold,
		EmbeddingProviderName: args.EmbeddingProviderName,
		EmbeddingModelName:    args.EmbeddingModelName,
	}

	marshalMessageBody, err := json.Marshal(enableAnnotationMessageBody)

	if err != nil {
		return nil, err
	}

	message := &primitive.Message{
		Topic: event.AnnotationTopic,
		Body:  marshalMessageBody,
	}

	message.WithTag(event.EnableAnnotationReplyTag)
	sendResult, err := as.mq.SendSync(ctx, message)

	if err != nil {
		return nil, errors.WithCode(code.ErrMQSend, "mq send sync error when send enable app annotation: %s", err.Error())
	}

	log.Infof("MQ-Send-Result %s", sendResult.String())

	return dto_app.NewApplyAnnotationWaiting(jobID), nil
}
