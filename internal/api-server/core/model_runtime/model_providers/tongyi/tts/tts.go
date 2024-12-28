package tts

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_providers/tongyi/base"
	provider_register "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
)

type TongyiTTSModel struct {
}

func init() {
	NewTongyiTTSModel().Register()
}

func NewTongyiTTSModel() *TongyiTTSModel {
	return &TongyiTTSModel{}
}

func (m *TongyiTTSModel) RegisterName() string {
	return "tongyi/tts"
}

func (m *TongyiTTSModel) Register() {
	provider_register.TTSModelRuntimeRegistry.RegisterLargeModelInstance(m)
}

func (m *TongyiTTSModel) Invoke(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, user, tenantID string, voice string, modelRuntime biz_entity.IAIModelRuntime, format string, texts []string) error {
	var err error

	ginContext := ctx.(*gin.Context)
	ginContext.Writer.Header().Set("Content-Type", "audio/mpeg")
	ginContext.Writer.WriteHeader(http.StatusOK)

	credentials = m.addCustomParameters(credentials)

	tongyiSDK := base.NewTongyiTTSSDK(credentials["api_key"].(string), credentials["endpoint_url"].(string), model, voice, format, nil)

	defer tongyiSDK.Close()

	go tongyiSDK.Generate(ctx, texts)

Finish:
	for {
		select {
		case audioMessage := <-tongyiSDK.GetAudioBinaryQueues():
			ginContext.Writer.Write(audioMessage)
			ginContext.Writer.Flush()
		case event := <-tongyiSDK.GetEventQueue():
			// log.Infof("receive event %s", event.Header.Event)
			if event.Header.Event == "task-finished" || event.Header.Event == "task-failed" {
				for (len(tongyiSDK.GetAudioBinaryQueues())) > 0 {
					audioMessage := <-tongyiSDK.GetAudioBinaryQueues()
					ginContext.Writer.Write(audioMessage)
					ginContext.Writer.Flush()
				}
				break Finish
			}
		case err = <-tongyiSDK.GetErrorQueues():
			break Finish
		}
	}
	tongyiSDK.CloseConnection()
	<-tongyiSDK.GetDone()
	return err
}

func (m *TongyiTTSModel) addCustomParameters(credentials map[string]interface{}) map[string]interface{} {
	credentials["endpoint_url"] = "wss://dashscope.aliyuncs.com/api-ws/v1/inference/"
	credentials["api_key"] = credentials["dashscope_api_key"]
	return credentials
}
