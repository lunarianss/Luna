// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dto

import (
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
)

type WebAppParameterResponse struct {
	OpeningStatement              map[string]interface{}              `json:"opening_statement" gorm:"column:opening_statement;serializer:json"`
	SuggestedQuestions            []string                            `json:"suggested_questions" gorm:"column:suggested_questions;serializer:json"`
	SuggestedQuestionsAfterAnswer po_entity.AppModelConfigEnable      `json:"suggested_questions_after_answer" gorm:"column:suggested_questions_after_answer;serializer:json"`
	TextToSpeech                  po_entity.AppModelConfigEnable      `json:"text_to_speech" gorm:"column:text_to_speech;serializer:json"`
	SpeechToText                  po_entity.AppModelConfigEnable      `json:"speech_to_text" gorm:"column:speech_to_text;serializer:json"`
	RetrieverResource             po_entity.AppModelConfigEnable      `json:"retriever_resource" gorm:"column:retriever_resource;serializer:json"`
	MoreLikeThis                  po_entity.AppModelConfigEnable      `json:"more_like_this" gorm:"column:more_like_this;serializer:json"`
	UserInputForm                 []map[string]map[string]interface{} `json:"user_input_form" gorm:"column:user_input_form;serializer:json"`
	SensitiveWordAvoidance        map[string]interface{}              `json:"sensitive_word_avoidance" gorm:"column:sensitive_word_avoidance;serializer:json"`
	FileUpload                    map[string]map[string]interface{}   `json:"file_upload" gorm:"column:file_upload;serializer:json"`
}

func AppConfigRecordToParameter(appConfig *po_entity.AppModelConfig) *WebAppParameterResponse {

	parameterDetail := &WebAppParameterResponse{
		OpeningStatement:              appConfig.OpeningStatement,
		SuggestedQuestions:            appConfig.SuggestedQuestions,
		SuggestedQuestionsAfterAnswer: appConfig.SuggestedQuestionsAfterAnswer,
		TextToSpeech:                  appConfig.TextToSpeech,
		SpeechToText:                  appConfig.SpeechToText,
		RetrieverResource:             appConfig.RetrieverResource,
		MoreLikeThis:                  appConfig.MoreLikeThis,
		UserInputForm:                 appConfig.UserInputForm,
		SensitiveWordAvoidance:        appConfig.SensitiveWordAvoidance,
		FileUpload:                    appConfig.FileUpload,
	}

	if !parameterDetail.RetrieverResource.Enable {
		parameterDetail.RetrieverResource.Enable = true
	}

	if parameterDetail.SensitiveWordAvoidance == nil {
		parameterDetail.SensitiveWordAvoidance = map[string]any{
			"enabled": false,
			"type":    "",
			"configs": []any{},
		}
	}

	if parameterDetail.FileUpload == nil {
		parameterDetail.FileUpload = map[string]map[string]interface{}{
			"image": {
				"enabled":          false,
				"number_limits":    3,
				"detail":           "high",
				"transfer_methods": []string{"remote_url", "local_file"},
			},
		}
	}

	if parameterDetail.UserInputForm == nil {
		parameterDetail.UserInputForm = []map[string]map[string]interface{}{}
	}

	return parameterDetail
}
