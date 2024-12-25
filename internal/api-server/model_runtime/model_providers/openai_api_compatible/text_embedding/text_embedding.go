package text_embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	biz_entity_model "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider/model_provider"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type IOpenApiCompactTextEmbeddingModel interface {
	Invoke(ctx context.Context) (*biz_entity.TextEmbeddingResult, error)
}

type openApiCompactLargeLanguageModel struct {
	credentials map[string]interface{}
	model       string
	texts       []string
	biz_entity_model.AIModelRuntime
}

func NewOpenApiCompactLargeLanguageModel(ctx context.Context, model string, credentials map[string]interface{}, texts []string) *openApiCompactLargeLanguageModel {
	return &openApiCompactLargeLanguageModel{
		model:       model,
		credentials: credentials,
		texts:       texts,
	}
}

func (o *openApiCompactLargeLanguageModel) Invoke(ctx context.Context) (*biz_entity.TextEmbeddingResult, error) {
	var (
		err            error
		batchEmbedding [][]float32
	)
	headers := map[string]string{
		"Content-Type":   "application/json",
		"Accept-Charset": "utf-8",
	}

	if apiKey, ok := o.credentials["api_key"]; ok {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	}

	endpointUrl, ok := o.credentials["endpoint_url"].(string)

	if !ok || endpointUrl == "" {
		return nil, errors.WithCode(code.ErrModelNotHaveEndPoint, "Model %s not have endpoint url", o.model)
	}

	if !strings.HasSuffix(endpointUrl, "/") {
		endpointUrl = fmt.Sprintf("%s/", endpointUrl)
	}

	requestData := map[string]interface{}{
		"model":           o.model,
		"encoding_format": "float",
	}

	endpointUrl, err = url.JoinPath(endpointUrl, "embeddings")

	if err != nil {
		return nil, errors.WithSCode(code.ErrRunTimeCaller, err.Error())
	}

	requestData["input"] = o.texts

	client := http.Client{
		Timeout: time.Duration(300) * time.Second,
	}

	log.Infof("Invoke llm request body %+v", requestData)
	requestBodyData, err := json.Marshal(requestData)

	if err != nil {
		return nil, errors.WithSCode(code.ErrEncodingJSON, err.Error())
	}

	req, err := http.NewRequest("POST", endpointUrl, bytes.NewReader(requestBodyData))

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

	decoder := json.NewDecoder(response.Body)
	var LLMResult biz_entity.TextEmbeddingLargeModelResult

	if err = decoder.Decode(&LLMResult); err != nil {
		return nil, errors.WithSCode(code.ErrDecodingJSON, err.Error())
	}

	log.Info("Text embedding response %+v", LLMResult)

	for _, embedding := range LLMResult.Data {
		batchEmbedding = append(batchEmbedding, embedding.Embedding)
	}

	defer response.Body.Close()

	embeddingUsage, err := o.calcResponseUsage(LLMResult.Usage.TotalTokens)

	if err != nil {
		return nil, err
	}

	return &biz_entity.TextEmbeddingResult{
		Usage:      embeddingUsage,
		Embeddings: batchEmbedding,
		Model:      o.model,
	}, nil
}

func (o *openApiCompactLargeLanguageModel) calcResponseUsage(tokens int) (*biz_entity.EmbeddingUsage, error) {

	priceInfo, err := o.GetPrice(o.model, o.credentials, biz_entity_model.INPUT, int64(tokens))

	if err != nil {
		return nil, err
	}

	return &biz_entity.EmbeddingUsage{
		Tokens:      tokens,
		TotalTokens: tokens,
		UnitPrice:   priceInfo.UnitPrice,
		PriceUnit:   priceInfo.Unit,
		TotalPrice:  priceInfo.TotalAmount,
		Currency:    priceInfo.Currency,
		Latency:     0.3,
	}, nil
}
