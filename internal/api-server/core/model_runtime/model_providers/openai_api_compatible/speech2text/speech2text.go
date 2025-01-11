package speech2text

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/lunarianss/Luna/infrastructure/errors"

	biz_entity_openai_standard_response "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/openai_standard_response"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type IOpenAudioApiCompactLargeLanguage interface {
	Invoke(ctx context.Context) (*biz_entity_openai_standard_response.Speech2TextResp, error)
}

type OpenAudioApiCompactLargeLanguage struct {
	biz_entity.IAIModelRuntime
	model            string
	credentials      map[string]interface{}
	audioFileContent []byte
	modelParameters  map[string]interface{}
	filename         string
}

func NewOpenAudioApiCompactLargeLanguage(audioFileContent []byte, modelParameters map[string]interface{}, credentials map[string]interface{}, model, filename string, modelRuntime biz_entity.IAIModelRuntime) *OpenAudioApiCompactLargeLanguage {
	return &OpenAudioApiCompactLargeLanguage{
		audioFileContent: audioFileContent,
		credentials:      credentials,
		modelParameters:  modelParameters,
		model:            model,
		IAIModelRuntime:  modelRuntime,
		filename:         filename,
	}
}

func (m *OpenAudioApiCompactLargeLanguage) Invoke(ctx context.Context) (*biz_entity_openai_standard_response.Speech2TextResp, error) {

	headers := map[string]string{
		"Accept-Charset": "utf-8",
	}

	if extraHeaders, ok := m.credentials["extra_headers"]; ok {
		if extraHeadersMap, ok := extraHeaders.(map[string]string); ok {
			for k, v := range extraHeadersMap {
				if _, ok := headers[k]; !ok {
					headers[k] = v
				}
			}
		}
	}

	if apiKey, ok := m.credentials["api_key"]; ok {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	}

	endpointUrl, ok := m.credentials["endpoint_url"]

	if !ok || endpointUrl == "" {

		return nil, errors.WithCode(code.ErrModelNotHaveEndPoint, "Model %s not have endpoint url", m.model)
	}

	endpointUrlStr, ok := endpointUrl.(string)

	if !ok {
		return nil, errors.WithCode(code.ErrModelNotHaveEndPoint, "Model %s not have endpoint url", m.model)
	}

	if !strings.HasSuffix(endpointUrlStr, "/") {
		endpointUrlStr = fmt.Sprintf("%s/", endpointUrlStr)
	}

	endpointJoinUrl, err := url.JoinPath(endpointUrlStr, "audio/transcriptions")

	if err != nil {
		return nil, err
	}

	var requestBody bytes.Buffer

	writer := multipart.NewWriter(&requestBody)

	if err := writer.WriteField("model", m.model); err != nil {
		return nil, err
	}

	filePart, err := writer.CreateFormFile("file", m.filename)

	if err != nil {
		return nil, err
	}

	_, err = filePart.Write(m.audioFileContent)

	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpointJoinUrl, &requestBody)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	transBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	speechText := &biz_entity_openai_standard_response.Speech2TextResp{}

	json.Unmarshal(transBytes, speechText)

	return speechText, err
}
