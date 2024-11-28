package model_registry

import (
	"context"
	"fmt"

	biz_entity_model "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_configuration"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	"github.com/lunarianss/Luna/pkg/log"
)

type ModelInstance struct {
	ProviderModelBundle *biz_entity.ProviderModelBundleRuntime
	Model               string
	Provider            string
	Credentials         map[string]interface{}
	ModelTypeInstance   *biz_entity_model.AIModelRuntime
}

func (ac *ModelInstance) InvokeLLM(ctx context.Context, promptMessage []*message.PromptMessage, queueManager *model_runtime.StreamGenerateQueue, modelParameters map[string]interface{}, tools interface{}, stop []string, stream bool, user string, callbacks interface{}) {

	modelKeyMapInvoke := fmt.Sprintf("%s/%s", ac.Provider, ac.ProviderModelBundle.ModelTypeInstance.ModelType)

	log.Infof("invoke %s", modelKeyMapInvoke)

	AIModelIns, err := ModelRuntimeRegistry.Acquire(modelKeyMapInvoke)

	if err != nil {
		queueManager.PushErr(err)
	}

	AIModelIns.Invoke(ctx, queueManager, ac.Model, ac.Credentials, modelParameters, stop, stream, user, promptMessage)
}
