// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/domain/app/entity/po_entity"
	biz_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	po_entity_chat "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
	"github.com/shopspring/decimal"
)

type IOpenApiCompactLargeLanguage interface {
	Invoke(ctx context.Context, queue *biz_entity_chat.StreamGenerateQueue)
	InvokeNonStream(ctx context.Context) (*biz_entity_chat.LLMResult, error)
}

type openApiCompactLargeLanguageModel struct {
	*biz_entity_chat.StreamGenerateQueue
	biz_entity.IAIModelRuntime
	FullAssistantContent string
	Usage                interface{}
	ChunkIndex           int
	Delimiter            string
	Model                string
	User                 string
	Stop                 []string
	Credentials          map[string]interface{}
	PromptMessages       []po_entity_chat.IPromptMessage
	ModelParameters      map[string]interface{}
	finalToolCalls       []*biz_entity_chat.ToolCallStream
	agent                bool
	tools                []*biz_entity_chat.PromptMessageTool
}

func NewOpenApiCompactLargeLanguageModel(promptMessages []po_entity_chat.IPromptMessage, modelParameters map[string]interface{}, credentials map[string]interface{}, model string, modelRuntime biz_entity.IAIModelRuntime, tools []*biz_entity_chat.PromptMessageTool) *openApiCompactLargeLanguageModel {
	return &openApiCompactLargeLanguageModel{
		PromptMessages:  promptMessages,
		Credentials:     credentials,
		ModelParameters: modelParameters,
		Model:           model,
		IAIModelRuntime: modelRuntime,
		tools:           tools,
	}
}

func (m *openApiCompactLargeLanguageModel) InvokeNonStream(ctx context.Context) (*biz_entity_chat.LLMResult, error) {
	if len(m.tools) > 0 {
		m.agent = true
	}

	return m.generateNonStream(ctx)
}

func (m *openApiCompactLargeLanguageModel) Invoke(ctx context.Context, queue *biz_entity_chat.StreamGenerateQueue) {
	if len(m.tools) > 0 {
		m.agent = true
	}
	m.StreamGenerateQueue = queue
	m.generate(ctx)
}

func (m *openApiCompactLargeLanguageModel) generateNonStream(ctx context.Context) (*biz_entity_chat.LLMResult, error) {
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

		return nil, errors.WithCode(code.ErrModelNotHaveEndPoint, "Model %s not have endpoint url", m.Model)
	}

	endpointUrlStr, ok := endpointUrl.(string)

	if !ok {

		return nil, errors.WithCode(code.ErrModelNotHaveEndPoint, "Model %s not have endpoint url", m.Model)
	}

	if !strings.HasSuffix(endpointUrlStr, "/") {
		endpointUrlStr = fmt.Sprintf("%s/", endpointUrlStr)
	}

	requestData := map[string]interface{}{
		"model":  m.Model,
		"stream": false,
	}

	for k, v := range m.ModelParameters {
		requestData[k] = v
	}
	messageItems := make([]map[string]interface{}, 0)

	completionType := m.Credentials["mode"]
	if completionType == string(po_entity.CHAT) {
		endpointJoinUrl, err := url.JoinPath(endpointUrlStr, "chat/completions")

		if err != nil {
			return nil, errors.WithSCode(code.ErrRunTimeCaller, err.Error())
		}
		endpointUrlStr = endpointJoinUrl

		for _, promptMessage := range m.PromptMessages {
			messageItem, err := promptMessage.ConvertToRequestData()

			if err != nil {
				return nil, err
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
		return nil, errors.WithSCode(code.ErrEncodingJSON, err.Error())
	}

	req, err := http.NewRequest("POST", endpointUrlStr, bytes.NewReader(requestBodyData))

	if err != nil {
		return nil, errors.WithSCode(code.ErrRunTimeCaller, err.Error())
	}

	if len(headers) > 0 {
		for headerKey, headerValue := range headers {
			req.Header.Set(headerKey, headerValue)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.WithSCode(code.ErrCallLargeLanguageModel, err.Error())
	}

	defer response.Body.Close()
	return m.handleNoStreamResponse(ctx, response)
}

func (m *openApiCompactLargeLanguageModel) handleNoStreamResponse(_ context.Context, response *http.Response) (*biz_entity_chat.LLMResult, error) {
	defer response.Body.Close()

	completion_type := m.Credentials["mode"].(string)
	var responseJSON biz_entity_chat.OpenaiResponse

	if err := json.NewDecoder(response.Body).Decode(&responseJSON); err != nil {
		return nil, errors.WithSCode(code.ErrDecodingJSON, err.Error())
	}

	choices := responseJSON.Choices

	if len(choices) == 0 {
		return nil, errors.WithCode(code.ErrCallLargeLanguageModel, "")
	}

	output := choices[0]

	messageId := responseJSON.ID

	responseContent := ""

	reason := output.FinishReason

	if completion_type == "chat" {
		responseContent = output.Message.Content
	}

	assistantMessage := biz_entity_chat.NewAssistantPromptMessage(responseContent)

	var (
		err               error
		ok                bool
		completeTokensInt int64
		promptTokensInt   int64
	)

	promptTokens := responseJSON.Usage.PromptTokens

	promptTokensInt, ok = handleTokenCount(promptTokens)

	if !ok {
		// 使用 openai tokenizlier
		promptTokensInt = 10
	}

	completeTokens := responseJSON.Usage.CompletionTokens

	completeTokensInt, ok = handleTokenCount(completeTokens)

	if !ok {
		// 使用 openai tokenizlier
		completeTokensInt = 20
	}

	promptPriceInfo, err := m.GetPrice(m.Model, m.Credentials, biz_entity.INPUT, promptTokensInt)

	if err != nil {
		return nil, err
	}

	completePriceInfo, err := m.GetPrice(m.Model, m.Credentials, biz_entity.OUTPUT, completeTokensInt)

	promptTotal := decimal.NewFromFloat(promptPriceInfo.TotalAmount)
	completeTotal := decimal.NewFromFloat(completePriceInfo.TotalAmount)
	totalAmount := promptTotal.Add(completeTotal).InexactFloat64()

	if err != nil {
		return nil, err
	}

	llmUsage := &biz_entity_chat.LLMUsage{
		PromptTokens:        promptTokensInt,
		PromptUnitPrice:     promptPriceInfo.UnitPrice,
		PromptPriceUnit:     promptPriceInfo.Unit,
		PromptPrice:         promptPriceInfo.TotalAmount,
		CompletionTokens:    completeTokensInt,
		CompletionUnitPrice: completePriceInfo.UnitPrice,
		CompletionPriceUnit: completePriceInfo.Unit,
		CompletionPrice:     completePriceInfo.TotalAmount,
		Currency:            promptPriceInfo.Currency,
		Latency:             1.0,
		TotalTokens:         promptTokensInt + completeTokensInt,
		TotalPrice:          totalAmount,
	}

	return &biz_entity_chat.LLMResult{
		ID:                messageId,
		Model:             m.Model,
		Message:           assistantMessage,
		Usage:             llmUsage,
		SystemFingerprint: responseJSON.SystemFingerprint,
		Reason:            reason,
	}, nil
}

func (m *openApiCompactLargeLanguageModel) generate(ctx context.Context) {
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
		m.PushErr(errors.WithCode(code.ErrModelNotHaveEndPoint, "Model %s not have endpoint url", m.Model))
		return
	}

	endpointUrlStr, ok := endpointUrl.(string)

	if !ok {
		m.PushErr(errors.WithCode(code.ErrModelNotHaveEndPoint, "Model %s not have endpoint url", m.Model))
		return
	}

	if !strings.HasSuffix(endpointUrlStr, "/") {
		endpointUrlStr = fmt.Sprintf("%s/", endpointUrlStr)
	}

	requestData := map[string]interface{}{
		"model":  m.Model,
		"stream": true,
	}

	for k, v := range m.ModelParameters {
		requestData[k] = v
	}
	messageItems := make([]map[string]interface{}, 0)

	completionType := m.Credentials["mode"]
	if completionType == string(po_entity.CHAT) {
		endpointJoinUrl, err := url.JoinPath(endpointUrlStr, "chat/completions")

		if err != nil {
			m.PushErr(errors.WithSCode(code.ErrRunTimeCaller, err.Error()))
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

	if len(m.tools) > 0 {
		requestData["tool_choice"] = "auto"

		messageFunctions := make([]*biz_entity_chat.PromptMessageFunction, 0)

		for _, tool := range m.tools {
			messageFunctions = append(messageFunctions, biz_entity_chat.NewFunctionTools(tool))
		}

		requestData["tools"] = messageFunctions
	}

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

	util.LogCompleteInfo(requestData)
	requestBodyData, err := json.Marshal(requestData)

	if err != nil {
		m.PushErr(errors.WithSCode(code.ErrEncodingJSON, err.Error()))
		return
	}

	req, err := http.NewRequest("POST", endpointUrlStr, bytes.NewReader(requestBodyData))

	if err != nil {
		m.PushErr(errors.WithSCode(code.ErrRunTimeCaller, err.Error()))
		return
	}

	if len(headers) > 0 {
		for headerKey, headerValue := range headers {
			req.Header.Set(headerKey, headerValue)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		m.PushErr(errors.WithSCode(code.ErrCallLargeLanguageModel, err.Error()))
		return
	}

	defer response.Body.Close()

	m.handleStreamResponse(ctx, response)
}

func (m *openApiCompactLargeLanguageModel) sendStreamChunkToQueue(ctx context.Context, messageId string, assistantPromptMessage *biz_entity_chat.AssistantPromptMessage) {
	streamResultChunk := &biz_entity_chat.LLMResultChunk{
		ID:            messageId,
		Model:         m.Model,
		PromptMessage: m.PromptMessages,
		Delta: &biz_entity_chat.LLMResultChunkDelta{
			Index:   m.ChunkIndex,
			Message: assistantPromptMessage,
		},
	}
	m.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, false, nil, false)
}

func (m *openApiCompactLargeLanguageModel) sendAgentStreamChunkToQueue(ctx context.Context, messageId string, assistantPromptMessage *biz_entity_chat.AssistantPromptMessage) {
	streamResultChunk := &biz_entity_chat.LLMResultChunk{
		ID:            messageId,
		Model:         m.Model,
		PromptMessage: m.PromptMessages,
		Delta: &biz_entity_chat.LLMResultChunkDelta{
			Index:   m.ChunkIndex,
			Message: assistantPromptMessage,
		},
	}
	m.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, false, nil, true)
}

func (m *openApiCompactLargeLanguageModel) sendAgentEndChunkToQueue(ctx context.Context, messageId string, assistantPromptMessage *biz_entity_chat.AssistantPromptMessage, finishReason string) {
	streamResultChunk := &biz_entity_chat.LLMResultChunk{
		ID:            messageId,
		Model:         m.Model,
		PromptMessage: m.PromptMessages,
		Delta: &biz_entity_chat.LLMResultChunkDelta{
			Index:        m.ChunkIndex,
			Message:      assistantPromptMessage,
			FinishReason: finishReason,
		},
	}
	m.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, false, nil, true)
}

func (m *openApiCompactLargeLanguageModel) sendErrorChunkToQueue(ctx context.Context, code error) {
	defer m.Close()
	err := errors.WithMessage(code, fmt.Sprintf("Error ocurred when handle stream: %#+v", code))
	m.HandleInvokeResultStream(ctx, nil, m.StreamGenerateQueue, false, err, false)
}

func handleTokenCount(count any) (int64, bool) {
	var countInt int64
	ok := true

	switch v := count.(type) {
	case float64:
		countInt = int64(v)
	case int:
		countInt = int64(v)
	case string:
		countFloat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			ok = false
			log.Infof("token count %v can't convert to float", count)
		}
		countInt = int64(countFloat)
	default:
		ok = false
	}

	if countInt == 0 {
		return 0, false
	} else {
		return countInt, ok
	}
}

func (m *openApiCompactLargeLanguageModel) convertToToolCall(toolStreamCalls []*biz_entity_chat.ToolCallStream) ([]*biz_entity_chat.ToolCall, error) {
	var (
		toolCalls []*biz_entity_chat.ToolCall
	)

	for _, toolCallStream := range toolStreamCalls {
		var args map[string]any

		if err := json.Unmarshal([]byte(toolCallStream.Function.Arguments), &args); err != nil {
			return nil, errors.WithSCode(code.ErrDecodingJSON, err.Error())
		}

		toolCalls = append(toolCalls, &biz_entity_chat.ToolCall{
			ID:   toolCallStream.ID,
			Type: toolCallStream.Type,
			Function: &biz_entity_chat.ToolCallFunction{
				Name:      toolCallStream.Function.Name,
				Arguments: args,
			},
		})

	}

	return toolCalls, nil

}
func (m *openApiCompactLargeLanguageModel) sendStreamFinalChunkToQueue(ctx context.Context, messageId string, finalReason string, fullAssistant string, usage map[string]interface{}) {
	defer m.Close()

	var (
		err               error
		ok                bool
		completeTokensInt int64
		promptTokensInt   int64
	)

	promptTokens, found := usage["prompt_tokens"]

	if found {
		promptTokensInt, ok = handleTokenCount(promptTokens)
	}

	if !ok {
		// 使用 openai tokenizlier
		promptTokensInt = 10
	}

	completeTokens, found := usage["completion_tokens"]

	if found {
		completeTokensInt, ok = handleTokenCount(completeTokens)
	}

	if !ok {
		// 使用 openai tokenizlier
		completeTokensInt = 20
	}

	promptPriceInfo, err := m.GetPrice(m.Model, m.Credentials, biz_entity.INPUT, promptTokensInt)

	if err != nil {
		m.HandleInvokeResultStream(ctx, nil, m.StreamGenerateQueue, false, err, false)
	}

	completePriceInfo, err := m.GetPrice(m.Model, m.Credentials, biz_entity.OUTPUT, completeTokensInt)

	promptTotal := decimal.NewFromFloat(promptPriceInfo.TotalAmount)
	completeTotal := decimal.NewFromFloat(completePriceInfo.TotalAmount)
	totalAmount := promptTotal.Add(completeTotal).InexactFloat64()

	if err != nil {
		m.HandleInvokeResultStream(ctx, nil, m.StreamGenerateQueue, false, err, false)
	}

	llmUsage := &biz_entity_chat.LLMUsage{
		PromptTokens:        promptTokensInt,
		PromptUnitPrice:     promptPriceInfo.UnitPrice,
		PromptPriceUnit:     promptPriceInfo.Unit,
		PromptPrice:         promptPriceInfo.TotalAmount,
		CompletionTokens:    completeTokensInt,
		CompletionUnitPrice: completePriceInfo.UnitPrice,
		CompletionPriceUnit: completePriceInfo.Unit,
		CompletionPrice:     completePriceInfo.TotalAmount,
		Currency:            promptPriceInfo.Currency,
		Latency:             1.0,
		TotalTokens:         promptTokensInt + completeTokensInt,
		TotalPrice:          totalAmount,
	}

	streamResultChunk := &biz_entity_chat.LLMResultChunk{
		ID:            messageId,
		Model:         m.Model,
		PromptMessage: m.PromptMessages,
		Delta: &biz_entity_chat.LLMResultChunkDelta{
			Index: m.ChunkIndex,
			Message: &biz_entity_chat.AssistantPromptMessage{
				PromptMessage: &po_entity_chat.PromptMessage{
					Content: fullAssistant,
				},
			},
			FinishReason: finalReason,
			Usage:        llmUsage,
		},
	}
	m.HandleInvokeResultStream(ctx, streamResultChunk, m.StreamGenerateQueue, true, nil, false)
}

func (m *openApiCompactLargeLanguageModel) handleStreamResponse(ctx context.Context, response *http.Response) {

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
		m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrConvertDelimiterString, "Can't convert delimiter %+v to string", delimiter))
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

	var usage = make(map[string]interface{})

	for scanner.Scan() {
		var (
			assistantPromptMessage *biz_entity_chat.AssistantPromptMessage
			deltaToolCall          []*biz_entity_chat.ToolCallStream
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
			m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrEncodingJSON, "JSON data %+v could not be decoded, failed: %+v", chunk, err.Error()))
			return
		}

		// groq 返回 error
		if apiError, ok := chunkJson["error"]; ok {
			if apiMapErr, ok := apiError.(map[string]interface{}); ok {
				if ok {
					apiByteErr, err := json.Marshal(apiMapErr)

					if err != nil {
						m.sendErrorChunkToQueue(ctx, errors.WithSCode(code.ErrEncodingJSON, err.Error()))
						return
					}

					m.sendErrorChunkToQueue(ctx, errors.WithSCode(code.ErrCallLargeLanguageModel, string(apiByteErr)))
					return
				}
			}
		}

		var chunkChoice = make(map[string]interface{})

		if usageChunk, ok := chunkJson["x_groq"]; ok {
			if v, ok := usageChunk.(map[string]interface{}); ok {
				usageAny, ok := v["usage"]
				if ok {
					usageMap, ok := usageAny.(map[string]interface{})
					if ok {
						usage = usageMap
					}
				}
			}
		}

		if usageChunk, ok := chunkJson["usage"]; ok {
			if usageChunk != nil {
				if v, ok := usageChunk.(map[string]interface{}); ok {
					if ok {
						usage = v
					}
				}
			}
		}

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
				deltaContent, existContent := deltaMap["content"]

				if tool_calls, ok := deltaMap["tool_calls"]; ok {

					toolCallSlice, ok := tool_calls.([]any)

					if !ok {
						m.sendErrorChunkToQueue(ctx, errors.WithCode(code.ErrInvokeToolUnConvertAble, "too_calls returned by llm isn't converted to []any"))
						return
					}

					toolCallByte, err := json.Marshal(toolCallSlice)

					if err != nil {
						m.sendErrorChunkToQueue(ctx, errors.WithSCode(code.ErrEncodingJSON, err.Error()))
						return
					}

					if err := json.Unmarshal(toolCallByte, &deltaToolCall); err != nil {
						m.sendErrorChunkToQueue(ctx, errors.WithSCode(code.ErrDecodingJSON, err.Error()))
						return
					}
					m.increaseToolCall(deltaToolCall)
				}

				if !existContent || deltaContent == "" {
					continue
				}

				if deltaContentStr, ok := deltaContent.(string); ok {
					m.FullAssistantContent += deltaContentStr
					assistantPromptMessage = biz_entity_chat.NewAssistantPromptMessage(deltaContentStr)
				}
			}
		} else {
			log.Warn("This chunk not property of delta and text")
			continue
		}

		if assistantPromptMessage != nil {
			if m.agent {
				m.sendAgentStreamChunkToQueue(ctx, messageID, assistantPromptMessage)
			} else {
				m.sendStreamChunkToQueue(ctx, messageID, assistantPromptMessage)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		m.sendErrorChunkToQueue(ctx, errors.WithSCode(code.ErrRunTimeCaller, err.Error()))
		return
	}

	if !m.agent {
		m.sendStreamFinalChunkToQueue(ctx, messageID, finishReason, m.FullAssistantContent, usage)
	} else {
		if len(m.finalToolCalls) > 0 {
			toolCalls, err := m.convertToToolCall(m.finalToolCalls)
			if err != nil {
				m.sendErrorChunkToQueue(ctx, errors.WithSCode(code.ErrRunTimeCaller, err.Error()))
				return
			}
			assistantPromptMessage := biz_entity_chat.NewAssistantPromptMessage(m.FullAssistantContent)
			assistantPromptMessage.ToolCalls = toolCalls

			m.sendAgentStreamChunkToQueue(ctx, messageID, assistantPromptMessage)
		} else {
			assistantPromptMessage := biz_entity_chat.NewAssistantPromptMessage(m.FullAssistantContent)
			m.sendAgentEndChunkToQueue(ctx, messageID, assistantPromptMessage, biz_entity_chat.AGENT_END)

		}
	}

}

func (m *openApiCompactLargeLanguageModel) getToolCall(tooCallID string) *biz_entity_chat.ToolCallStream {
	if tooCallID == "" {
		return m.finalToolCalls[len(m.finalToolCalls)-1]
	}
	for _, toolCall := range m.finalToolCalls {
		if toolCall.ID == tooCallID {
			return toolCall
		}
	}

	newToolCall := &biz_entity_chat.ToolCallStream{
		ID:   tooCallID,
		Type: "function",
		Function: &biz_entity_chat.ToolCallStreamFunction{
			Name:      "",
			Arguments: "",
		},
	}

	m.finalToolCalls = append(m.finalToolCalls, newToolCall)
	return newToolCall
}

func (m *openApiCompactLargeLanguageModel) increaseToolCall(newToolCalls []*biz_entity_chat.ToolCallStream) {
	for _, toolCall := range newToolCalls {
		toolCallFind := m.getToolCall(toolCall.ID)

		if toolCall.ID != "" {
			toolCallFind.ID = toolCall.ID
		}

		if toolCall.Type != "" {
			toolCallFind.Type = toolCall.Type
		}

		if toolCall.Function.Name != "" {
			toolCallFind.Function.Name = toolCall.Function.Name
		}

		if toolCall.Function.Arguments != "" {
			toolCallFind.Function.Arguments += toolCall.Function.Arguments
		}

	}

}

func (runner *openApiCompactLargeLanguageModel) HandleInvokeResultStream(ctx context.Context, invokeResult *biz_entity_chat.LLMResultChunk, streamGenerator *biz_entity_chat.StreamGenerateQueue, end bool, err error, isAgent bool) {

	if err != nil && invokeResult == nil {
		streamGenerator.Final(&biz_entity_chat.QueueErrorEvent{
			AppQueueEvent: biz_entity_chat.NewAppQueueEvent(biz_entity_chat.Error),
			Err:           err,
		})
		return
	}

	if end {
		llmResult := &biz_entity_chat.LLMResult{
			Model:         invokeResult.Model,
			PromptMessage: invokeResult.PromptMessage,
			Reason:        invokeResult.Delta.FinishReason,
			Message: &biz_entity_chat.AssistantPromptMessage{
				PromptMessage: &po_entity_chat.PromptMessage{
					Content: invokeResult.Delta.Message.Content,
				},
			},
			Usage: invokeResult.Delta.Usage,
		}

		event := biz_entity_chat.NewAppQueueEvent(biz_entity_chat.MessageEnd)
		streamGenerator.Final(&biz_entity_chat.QueueMessageEndEvent{
			AppQueueEvent: event,
			LLMResult:     llmResult,
		})
		return
	}

	if isAgent {
		event := biz_entity_chat.NewAppQueueEvent(biz_entity_chat.AgentMessage)
		streamGenerator.Push(&biz_entity_chat.QueueAgentMessageEvent{
			AppQueueEvent: event,
			Chunk:         invokeResult})
	} else {
		event := biz_entity_chat.NewAppQueueEvent(biz_entity_chat.LLMChunk)
		streamGenerator.Push(&biz_entity_chat.QueueLLMChunkEvent{
			AppQueueEvent: event,
			Chunk:         invokeResult})
	}
}
