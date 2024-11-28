package app

import (
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	biz_entity_provider_config "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
)

// Define types for your enums
type InvokeFrom string

const (
	ServiceAPI InvokeFrom = "service-api"
	WebApp     InvokeFrom = "web-app"
	Explore    InvokeFrom = "explore"
	Debugger   InvokeFrom = "debugger"
)

// ModelConfigWithCredentialsEntity struct
type ModelConfigWithCredentialsEntity struct {
	Provider            string                                                 `json:"provider"`
	Model               string                                                 `json:"model"`
	ModelSchema         *biz_entity.AIModelStaticConfiguration                 `json:"model_schema"`
	Mode                string                                                 `json:"mode"`
	ProviderModelBundle *biz_entity_provider_config.ProviderModelBundleRuntime `json:"provider_model_bundle"`
	Credentials         interface{}                                            `json:"credentials"`
	Parameters          map[string]interface{}                                 `json:"parameters"`
	Stop                []string                                               `json:"stop"`
}

// AppGenerateEntity struct
type AppGenerateEntity struct {
	TaskID     string                 `json:"task_id"`
	AppConfig  *app_config.AppConfig  `json:"app_config"`
	Inputs     map[string]interface{} `json:"inputs"`
	UserID     string                 `json:"user_id"`
	Stream     bool                   `json:"stream"`
	InvokeFrom InvokeFrom             `json:"invoke_from"`
	CallDepth  int                    `json:"call_depth"`
	Extras     map[string]interface{} `json:"extras"`
}

// EasyUIBasedAppGenerateEntity struct
type EasyUIBasedAppGenerateEntity struct {
	*AppGenerateEntity
	AppConfig *app_config.EasyUIBasedAppConfig  `json:"app_config"`
	ModelConf *ModelConfigWithCredentialsEntity `json:"model_conf"`
	Query     string                            `json:"query"`
}

// ConversationAppGenerateEntity struct
type ConversationAppGenerateEntity struct {
	*AppGenerateEntity
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
}

// ChatAppGenerateEntity struct
type ChatAppGenerateEntity struct {
	*EasyUIBasedAppGenerateEntity
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
}
