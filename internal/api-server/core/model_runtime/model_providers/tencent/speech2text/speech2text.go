package speech2text

import (
	"context"
	"strings"

	"github.com/lunarianss/Luna/infrastructure/errors"
	provider_register "github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	biz_entity_openai_standard_response "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/openai_standard_response"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/tencentcloud/tencentcloud-speech-sdk-go/asr"
	"github.com/tencentcloud/tencentcloud-speech-sdk-go/common"
)

type tencentAudioLargeLanguageModel struct {
}

func init() {
	NewTencentAudioLargeLanguageModel().Register()
}

func NewTencentAudioLargeLanguageModel() *tencentAudioLargeLanguageModel {
	return &tencentAudioLargeLanguageModel{}
}

var _ provider_register.IAudioModelRegistry = (*tencentAudioLargeLanguageModel)(nil)

func (m *tencentAudioLargeLanguageModel) Invoke(ctx context.Context, model string, credentials map[string]interface{}, modelParameters map[string]interface{}, user, filename string, fileContent []byte, modelRuntime biz_entity.IAIModelRuntime) (*biz_entity_openai_standard_response.Speech2TextResp, error) {

	credential := common.NewCredential(credentials["secret_id"].(string), credentials["secret_key"].(string))

	recognizer := asr.NewFlashRecognizer(credentials["app_id"].(string), credential)

	req := new(asr.FlashRecognitionRequest)
	req.EngineType = "16k_zh"
	req.VoiceFormat = "mp3"
	req.SpeakerDiarization = 0
	req.FilterDirty = 0
	req.FilterModal = 0
	req.FilterPunc = 0
	req.ConvertNumMode = 1
	req.FirstChannelOnly = 1
	req.WordInfo = 0

	resp, err := recognizer.Recognize(req, fileContent)

	if err != nil {
		return nil, errors.WithCode(code.ErrTencentARS, err.Error())
	}

	var tranStr []string

	for _, channelResult := range resp.FlashResult {
		tranStr = append(tranStr, channelResult.Text)
	}

	return &biz_entity_openai_standard_response.Speech2TextResp{
		Text: strings.Join(tranStr, ""),
	}, nil
}

func (m *tencentAudioLargeLanguageModel) Register() {
	provider_register.AudioModelRuntimeRegistry.RegisterLargeModelInstance(m)
}

func (m *tencentAudioLargeLanguageModel) RegisterName() string {
	return "tencent/speech2text"
}
