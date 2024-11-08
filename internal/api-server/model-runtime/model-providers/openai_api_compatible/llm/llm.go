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

	"github.com/lunarianss/Luna/internal/api-server/entities/base"
	"github.com/lunarianss/Luna/internal/api-server/entities/llm"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
	"github.com/lunarianss/Luna/pkg/log"
)

type OpenApiCompactLargeLanguageModel struct {
}

func (oacll *OpenApiCompactLargeLanguageModel) Invoke(ctx context.Context, model string, credentials map[string]interface{}, promptMessages []*message.PromptMessage, modelParameters map[string]interface{}, stop []string, stream bool, user string) error {
	return oacll.generate(ctx, model, credentials, modelParameters, stop, stream, user)
}

func (m *OpenApiCompactLargeLanguageModel) generate(ctx context.Context, model string, credentials map[string]interface{}, promptMessages []*message.PromptMessage, modelParameters map[string]interface{}, stop []string, stream bool, user string) error {

	headers := map[string]string{
		"Content-Type":   "application/json",
		"Accept-Charset": "utf-8",
	}

	if extraHeaders, ok := credentials["extra_headers"]; ok {
		if extraHeadersMap, ok := extraHeaders.(map[string]string); ok {
			for k, v := range extraHeadersMap {
				if _, ok := headers[k]; !ok {
					headers[k] = v
				}
			}
		}
	}

	if apiKey, ok := credentials["api_key"]; ok {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	}

	endpointUrl, ok := credentials["endpoint_url"]

	if !ok || endpointUrl == "" {
		return errors.WithCode(code.ErrModelNotHaveEndPoint, fmt.Sprintf("model %s not have endpoint url", model))
	}

	endpointUrlStr, ok := endpointUrl.(string)

	if !ok {
		return errors.WithCode(code.ErrModelNotHaveEndPoint, fmt.Sprintf("model %s not have endpoint url", model))
	}

	if !strings.HasSuffix(endpointUrlStr, "/") {
		endpointUrlStr = fmt.Sprintf("%s/", endpointUrlStr)
	}

	requestData := map[string]interface{}{
		"model":  model,
		"stream": stream,
	}

	for k, v := range modelParameters {
		requestData[k] = v
	}
	messageItems := make([]map[string]interface{}, 0)

	//todo util now is only support simple chat
	completionType := credentials["model"]
	if completionType == base.CHAT {
		endpointJoinUrl, err := url.JoinPath(endpointUrlStr, "chat/completions")

		if err != nil {
			return errors.WithCode(code.ErrRunTimeCaller, err.Error())
		}
		endpointUrlStr = endpointJoinUrl

		for _, promptMessage := range promptMessages {
			messageItem, err := promptMessage.ConvertToRequestData()

			if err != nil {
				return err
			}
			messageItems = append(messageItems, messageItem)
		}
	}

	requestData["message"] = messageItems

	if stop != nil && len(stop) > 1 {
		requestData["stop"] = stop
	}

	if user != "" {
		requestData["user"] = user
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

	defer response.Body.Close()

	if err != nil {
		return errors.WithCode(code.ErrCallLargeLanguageModel, err.Error())
	}

	if stream {
		m.handleStreamResponse(model, credentials, promptMessages, response)
	}
	return nil
}

func (m *OpenApiCompactLargeLanguageModel) handleStreamResponse(model string, credentials map[string]interface{}, promptMessages []*message.PromptMessage, response *http.Response) error {

	var full_assistant_content string
	var assistant_prompt_message *message.AssistantPromptMessage

	delimiter, ok := credentials["stream_mode_delimiter"]
	if !ok {
		delimiter = "\n\n"
	}

	delimiterStr, ok := delimiter.(string)

	if !ok {
		return errors.WithCode(code.ErrConvertDelimiterString, fmt.Sprintf("cannot convert delimiter %+v to string", delimiter))
	}

	scanner := bufio.NewScanner(response.Body)

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := strings.Index(string(data), delimiterStr); i >= 0 {
			return i + len(delimiterStr), data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})

	chunkIndex := 0

	for scanner.Scan() {
		chunk := strings.TrimSpace(scanner.Text())

		if chunk == "" || strings.HasPrefix(chunk, ":") {
			continue
		}

		chunk = strings.TrimPrefix(chunk, "data: ")
		chunk = strings.TrimSpace(chunk)
		var chunkJson map[string]interface{}

		err := json.Unmarshal([]byte(chunk), chunkJson)

		if err != nil {
			// 这里给结束消息，而不是只记得报错
			return errors.WithCode(code.ErrDecodingJSON, err.Error())
		}

		var chunkChoice = make(map[string]interface{})

		if chunkChoices, ok := chunkJson["choices"]; ok {
			if v, ok := chunkChoices.([]map[string]interface{}); ok {
				chunkChoice = v[0]
			}
		}

		messageId, ok := chunkChoice["id"].(string)

		if !ok {
			messageId = ""
		}

		// finishReason := chunkChoice["finish_reason"]

		chunkIndex += 1

		if delta, ok := chunkChoice["delta"]; ok {
			if deltaMap, ok := delta.(map[string]interface{}); ok {

				deltaContent := deltaMap["content"]

				assistant_prompt_message = &message.AssistantPromptMessage{
					PromptMessage: &message.PromptMessage{
						Role:    message.ASSISTANT,
						Content: deltaContent,
					},
				}

				if deltaContentStr, ok := deltaContent.(string); ok {
					full_assistant_content += deltaContentStr
				}
			}
		} else {
			log.Warn("this chunk not property of delta and text")
			continue
		}

		streamResultChunk := &llm.LLMResultChunk{
			ID:            messageId,
			Model:         model,
			PromptMessage: promptMessages,
			Delta: &llm.LLMResultChunkDelta{
				Index:   chunkIndex,
				Message: assistant_prompt_message,
			},
		}

		// 将 chunk 处理

	}

	if err := scanner.Err(); err != nil {
		return errors.WithCode(code.ErrRunTimeCaller, err.Error())
	}

	return nil
}
