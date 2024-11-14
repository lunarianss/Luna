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

func (m *OpenApiCompactLargeLanguageModel) Invoke(ctx context.Context, promptMessages []*message.PromptMessage, modelParameters map[string]interface{}, credentials map[string]interface{}, queueManager *model_runtime.StreamGenerateQueue) {
	m.Credentials = credentials
	m.AppRunner = &apps.AppRunner{}
	m.ModelParameters = modelParameters
	m.PromptMessages = promptMessages
	m.StreamGenerateQueue = queueManager

	m.generate(ctx)
}

func (m *OpenApiCompactLargeLanguageModel) generate(ctx context.Context) {
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
		m.PushErr(errors.WithCode(code.ErrModelNotHaveEndPoint, fmt.Sprintf("Model %s not have endpoint url", m.Model)))
		return
	}

	endpointUrlStr, ok := endpointUrl.(string)

	if !ok {
		m.PushErr(errors.WithCode(code.ErrModelNotHaveEndPoint, fmt.Sprintf("Model %s not have endpoint url", m.Model)))
		return
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

	completionType := m.Credentials["mode"]
	if completionType == string(base.CHAT) {
		endpointJoinUrl, err := url.JoinPath(endpointUrlStr, "chat/completions")

		if err != nil {
			m.PushErr(errors.WithCode(code.ErrRunTimeCaller, err.Error()))
			return
		}
		endpointUrlStr = endpointJoinUrl

		for _, promptMessage := range m.PromptMessages {
			messageItem, err := promptMessage.ConvertToRequestData()

			if err != nil {
				m.PushErr(err)
				return
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

	log.Infof("Invoke llm request body %+v", requestData)
	requestBodyData, err := json.Marshal(requestData)

	if err != nil {
		m.PushErr(errors.WithCode(code.ErrEncodingJSON, err.Error()))
		return
	}

	req, err := http.NewRequest("POST", endpointUrlStr, bytes.NewReader(requestBodyData))

	if err != nil {
		m.PushErr(errors.WithCode(code.ErrRunTimeCaller, err.Error()))
		return
	}

	if len(headers) > 0 {
		for headerKey, headerValue := range headers {
			req.Header.Set(headerKey, headerValue)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		m.PushErr(errors.WithCode(code.ErrCallLargeLanguageModel, err.Error()))
		return
	}

	defer response.Body.Close()

	if m.Stream {
		m.handleStreamResponse(ctx, response)
	}
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
	m.AppRunner.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, false, nil)
}

func (m *OpenApiCompactLargeLanguageModel) sendErrorChunkToQueue(ctx context.Context, code error) {
	defer m.Close()
	err := errors.WithMessage(code, fmt.Sprintf("Error ocurred when handle stream: %#+v", code))
	m.AppRunner.HandleInvokeResultStream(ctx, nil, m.StreamGenerateQueue, false, err)
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
	m.AppRunner.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, true, nil)
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
		m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrConvertDelimiterString, fmt.Sprintf("Can't convert delimiter %+v to string", delimiter)))
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
			m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrEncodingJSON, fmt.Sprintf("JSON data %+v could not be decoded, failed: %+v", chunk, err.Error())))
			return
		}

		// groq 返回 error
		if apiError, ok := chunkJson["error"]; ok {
			if apiMapErr, ok := apiError.(map[string]interface{}); ok {
				if ok {
					apiByteErr, err := json.Marshal(apiMapErr)

					if err != nil {
						m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrEncodingJSON, err.Error()))
						return
					}

					m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrCallLargeLanguageModel, string(apiByteErr)))
					return
				}
			}
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
			finishReason = "Finish_reason doesn't convert to string"
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
			log.Warn("This chunk not property of delta and text")
			continue
		}

		m.sendStreamChunkToQueue(ctx, messageID, assistantPromptMessage)
	}

	if err := scanner.Err(); err != nil {
		m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrRunTimeCaller, err.Error()))
		return
	}

	log.Infof("Full Answer From AI: %s", m.FullAssistantContent)

	m.sendStreamFinalChunkToQueue(ctx, messageID, finishReason)
}
