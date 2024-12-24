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
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/provider/domain_service"
	dto_app "github.com/lunarianss/Luna/internal/api-server/dto/app"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/chat"
	"github.com/lunarianss/Luna/internal/api-server/event"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AnnotationService struct {
	appDomain      *appDomain.AppDomain
	providerDomain *domain_service.ProviderDomain
	accountDomain  *accountDomain.AccountDomain
	chatDomain     *chatDomain.ChatDomain

	redis *redis.Client

	mq rocketmq.Producer
}

func NewAnnotationService(appDomain *appDomain.AppDomain, providerDomain *domain_service.ProviderDomain, accountDomain *accountDomain.AccountDomain, chatDomain *chatDomain.ChatDomain, redis *redis.Client, mq rocketmq.Producer) *AnnotationService {
	return &AnnotationService{
		appDomain:      appDomain,
		providerDomain: providerDomain,
		accountDomain:  accountDomain,
		chatDomain:     chatDomain,
		redis:          redis,
		mq:             mq,
	}
}

func (as *AnnotationService) InsertAnnotationFromMessage(ctx context.Context, accountID string, appID string, args *dto.InsertAnnotationFormMessage) (*dto.MessageAnnotation, error) {

	var (
		bizMessageAnnotation *biz_entity.BizMessageAnnotation
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

	_, err = as.chatDomain.AnnotationRepo.GetAnnotationSetting(ctx, app.ID, nil)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assembler.ConvertToAnnotation(bizMessageAnnotation), nil
		} else {
			return nil, err
		}
	} else {
		// todo rocketmq 异步消息
	}

	return nil, nil

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

	enableAnnotationMessageBody := event.EnableAnnotationReplyTask{
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

	sendResult, err := as.mq.SendSync(ctx, message)

	if err != nil {
		return nil, errors.WithCode(code.ErrMQSend, "mq send sync error when send enable app annotation: %s", err.Error())
	}

	log.Infof("MQ-Send-Result %s", sendResult.String())

	return dto_app.NewApplyAnnotationWaiting(jobID), nil
}
