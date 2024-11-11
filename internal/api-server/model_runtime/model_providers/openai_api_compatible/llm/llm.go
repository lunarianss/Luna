package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lunarianss/Luna/internal/api-server/core/app/apps"
	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/entities/llm"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model_runtime"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
)

type OpenApiCompactLargeLanguageModel struct {
	*apps.AppRunner
	*model_runtime.StreamGenerateQueue
	FullAssistantContent string
	Usage                interface{}
	ChunkIndex           int
	Delimiter            string
	Model                string
	Stream               bool
	User                 string
	Stop                 []string
	Credentials          map[string]interface{}
	PromptMessages       []*message.PromptMessage
	ModelParameters      map[string]interface{}
}

func (m *OpenApiCompactLargeLanguageModel) Invoke(ctx context.Context, promptMessages []*message.PromptMessage, modelParameters map[string]interface{}, credentials map[string]interface{}) error {
	m.Credentials = credentials
	m.ModelParameters = modelParameters
	m.PromptMessages = promptMessages

	return m.generate(ctx)
}

func (m *OpenApiCompactLargeLanguageModel) generate(ctx context.Context) error {
	headers := map[string]string{
		"Content-Type":   "application/json",
		"Accept-Charset": "utf-8",
	}

	if extraHeaders, ok := m.Credentials["extra_headers"]; ok {
		if extraHeadersMap, ok := extraHeaders.(map[string]string); ok {
			for k, v := range extraHeadersMap {
				if _, ok := headers[k]; !ok {
					headers[k] = v
				}
			}
		}
	}

	if apiKey, ok := m.Credentials["api_key"]; ok {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	}

	endpointUrl, ok := m.Credentials["endpoint_url"]

	if !ok || endpointUrl == "" {
		return errors.WithCode(code.ErrModelNotHaveEndPoint, fmt.Sprintf("model %s not have endpoint url", m.Model))
	}

	endpointUrlStr, ok := endpointUrl.(string)

	if !ok {
		return errors.WithCode(code.ErrModelNotHaveEndPoint, fmt.Sprintf("model %s not have endpoint url", m.Model))
	}

	if !strings.HasSuffix(endpointUrlStr, "/") {
		endpointUrlStr = fmt.Sprintf("%s/", endpointUrlStr)
	}

	requestData := map[string]interface{}{
		"model":  m.Model,
		"stream": m.Stream,
	}

	for k, v := range m.ModelParameters {
		requestData[k] = v
	}
	messageItems := make([]map[string]interface{}, 0)

	//todo util now is only support simple chat
	completionType := m.Credentials["mode"]
	if completionType == string(base.CHAT) {
		endpointJoinUrl, err := url.JoinPath(endpointUrlStr, "chat/completions")

		if err != nil {
			return errors.WithCode(code.ErrRunTimeCaller, err.Error())
		}
		endpointUrlStr = endpointJoinUrl

		for _, promptMessage := range m.PromptMessages {
			messageItem, err := promptMessage.ConvertToRequestData()

			if err != nil {
				return err
			}
			messageItems = append(messageItems, messageItem)
		}
	}

	requestData["messages"] = messageItems

	if len(m.Stop) > 1 {
		requestData["stop"] = m.Stop
	}

	if m.User != "" {
		requestData["user"] = m.User
	}

	client := http.Client{
		Timeout: time.Duration(300) * time.Second,
	}

	requestBodyData, err := json.Marshal(requestData)

	if err != nil {
		return errors.WithCode(code.ErrEncodingJSON, err.Error())
	}

	req, err := http.NewRequest("POST", endpointUrlStr, bytes.NewReader(requestBodyData))

	if err != nil {
		return errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	if len(headers) > 0 {
		for headerKey, headerValue := range headers {
			req.Header.Set(headerKey, headerValue)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return errors.WithCode(code.ErrCallLargeLanguageModel, err.Error())
	}

	defer response.Body.Close()

	if m.Stream {
		m.handleStreamResponse(ctx, response)
	}
	return nil
}

func (m *OpenApiCompactLargeLanguageModel) sendStreamChunkToQueue(ctx context.Context, messageId string, assistantPromptMessage *message.AssistantPromptMessage) {
	streamResultChunk := &llm.LLMResultChunk{
		ID:            messageId,
		Model:         m.Model,
		PromptMessage: m.PromptMessages,
		Delta: &llm.LLMResultChunkDelta{
			Index:   m.ChunkIndex,
			Message: assistantPromptMessage,
		},
	}
	m.AppRunner.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, false)
}

func (m *OpenApiCompactLargeLanguageModel) sendErrorChunkToQueue(ctx context.Context, code error, messageID string) {

	finishErrReason := fmt.Sprintf("error ocurred when handle stream: %#+v", code)

	log.Error(finishErrReason)

	m.sendStreamFinalChunkToQueue(ctx, messageID, finishErrReason)
}

func (m *OpenApiCompactLargeLanguageModel) sendStreamFinalChunkToQueue(ctx context.Context, messageId string, finalReason string) {
	defer m.Close()

	streamResultChunk := &llm.LLMResultChunk{
		ID:            messageId,
		Model:         m.Model,
		PromptMessage: m.PromptMessages,
		Delta: &llm.LLMResultChunkDelta{
			Index:        m.ChunkIndex,
			Message:      message.NewEmptyAssistantPromptMessage(),
			FinishReason: finalReason,
		},
	}
	m.AppRunner.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, true)
}

func (m *OpenApiCompactLargeLanguageModel) handleStreamResponse(ctx context.Context, response *http.Response) {

	var (
		messageID    string
		finishReason string
	)

	delimiter, ok := m.Credentials["stream_mode_delimiter"]
	if !ok {
		delimiter = "\n\n"
	}

	m.Delimiter, ok = delimiter.(string)

	if !ok {
		m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrConvertDelimiterString, fmt.Sprintf("cannot convert delimiter %+v to string", delimiter)), messageID)
		return
	}

	scanner := bufio.NewScanner(response.Body)

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := strings.Index(string(data), m.Delimiter); i >= 0 {
			return i + len(m.Delimiter), data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})

	for scanner.Scan() {
		var (
			assistantPromptMessage *message.AssistantPromptMessage
		)
		chunk := strings.TrimSpace(scanner.Text())

		if chunk == "" || strings.HasPrefix(chunk, ":") {
			continue
		}

		chunk = strings.TrimPrefix(chunk, "data: ")
		chunk = strings.TrimSpace(chunk)

		if chunk == "[DONE]" {
			continue
		}

		var chunkJson map[string]interface{}

		err := json.Unmarshal([]byte(chunk), &chunkJson)

		if err != nil {
			m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrEncodingJSON, fmt.Sprintf("JSON data %+v could not be decoded, failed: %+v", chunk, err.Error())), messageID)
		}

		var chunkChoice = make(map[string]interface{})

		if chunkChoices, ok := chunkJson["choices"]; ok {
			if v, ok := chunkChoices.([]interface{}); ok {
				if vv, ok := v[0].(map[string]interface{}); ok {
					chunkChoice = vv
				}
			}
		}

		messageID, ok = chunkChoice["id"].(string)

		if !ok {
			messageID = ""
		}

		finishReason, ok = chunkChoice["finish_reason"].(string)

		if !ok {
			finishReason = "finish_reason doesn't convert to string"
		}

		m.ChunkIndex += 1

		if delta, ok := chunkChoice["delta"]; ok {
			if deltaMap, ok := delta.(map[string]interface{}); ok {
				deltaContent := deltaMap["content"]
				assistantPromptMessage = message.NewAssistantPromptMessage(message.ASSISTANT, deltaContent)
				if deltaContentStr, ok := deltaContent.(string); ok {
					m.FullAssistantContent += deltaContentStr
				}
			}
		} else {
			log.Warn("this chunk not property of delta and text")
			continue
		}

		m.sendStreamChunkToQueue(ctx, messageID, assistantPromptMessage)
	}

	if err := scanner.Err(); err != nil {
		m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrRunTimeCaller, err.Error()), messageID)
		return
	}

	log.Infof("full answer %s", m.FullAssistantContent)

	m.sendStreamFinalChunkToQueue(ctx, messageID, finishReason)
}
