package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/core/tools/tool_registry"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

func init() {
	tool_registry.ToolRuntimeRegistry.RegisterAgentToolInstance(NewStableDiffusionTool())
}

type StableDiffusionTool struct {
	parameter        *StableDiffusionToolParameter
	payloadMap       map[string]string
	modelEndpointMap map[string]string
	credentials      string
	*biz_entity.ToolRuntimeConfiguration
}

var _ tool_registry.IToolCallRegistry = (*StableDiffusionTool)(nil)

func NewStableDiffusionTool() *StableDiffusionTool {
	return &StableDiffusionTool{
		parameter: &StableDiffusionToolParameter{
			AspectRatio:  "16:9",
			Mode:         "text-to-image",
			OutputFormat: "png",
			Model:        "core",
		},
		payloadMap: make(map[string]string, 6),
		modelEndpointMap: map[string]string{
			"sd3":       "https://api.stability.ai/v2beta/stable-image/generate/sd3",
			"sd3-turbo": "https://api.stability.ai/v2beta/stable-image/generate/sd3",
			"core":      "https://api.stability.ai/v2beta/stable-image/generate/core",
		},
	}
}

type StableDiffusionToolParameter struct {
	Prompt         string `json:"prompt"`
	AspectRatio    string `json:"aspect_ratio"`
	Mode           string `json:"mode"`
	Seed           int    `json:"seed"`
	OutputFormat   string `json:"output_format"`
	NegativePrompt string `json:"negative_prompt"`
	Model          string `json:"model"`
}

func (st *StableDiffusionTool) Register() string {
	return "stability/stability_text2image"
}

func (st *StableDiffusionTool) Invoke(ctx context.Context, userID string, toolParameters []byte, toolRuntime *biz_entity.ToolRuntimeConfiguration) ([]*biz_entity.ToolInvokeMessage, error) {

	log.Infof("----- toolParameter %s -----\n\n", string(toolParameters))

	if err := st.parseParameter(toolParameters, toolRuntime); err != nil {
		return nil, err
	}

	response, err := st.post(st.modelEndpointMap[st.parameter.Model])

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, errors.WithCode(code.ErrInvokeTool, "failed to invoke stable diffusion tool error: %s", err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.WithCode(code.ErrInvokeTool, "failed to invoke stable diffusion tool error: %s", body)
	}

	return []*biz_entity.ToolInvokeMessage{st.CreateBlobMessage(body, map[string]any{
		"mime_type": "image/png",
	}, "image")}, nil
}

func (st *StableDiffusionTool) parseParameter(parameter []byte, toolRuntime *biz_entity.ToolRuntimeConfiguration) error {

	if err := json.Unmarshal(parameter, st.parameter); err != nil {
		return errors.WithCode(code.ErrDecodingJSON, "decoding json bytes: %s, error %+v", string(parameter), err.Error())
	}

	model := st.parameter.Model

	if model == "sd3" || model == "sd3-turbo" {
		st.payloadMap["model"] = model
	}

	if model != "sd3-turbo" {
		st.payloadMap["negative_prompt"] = st.parameter.NegativePrompt
	}

	st.payloadMap["prompt"] = st.parameter.Prompt
	st.payloadMap["aspect_ratio"] = st.parameter.AspectRatio
	st.payloadMap["mode"] = st.parameter.Mode
	st.payloadMap["seed"] = strconv.Itoa(st.parameter.Seed)
	st.payloadMap["output_format"] = "png"

	credentials, ok := toolRuntime.Credentials["api_key"].(string)

	if !ok {
		errors.WithSCode(code.ErrRunTimeCaller, "stable diffusion tool api key is not converted to string")
	}

	st.credentials = credentials

	return nil
}

func (st *StableDiffusionTool) post(url string) (*http.Response, error) {

	var b bytes.Buffer

	writer := multipart.NewWriter(&b)

	for key, value := range st.payloadMap {
		if err := writer.WriteField(key, value); err != nil {
			return nil, errors.WithSCode(code.ErrRunTimeCaller, err.Error())
		}
	}

	if err := writer.Close(); err != nil {
		return nil, errors.WithSCode(code.ErrRunTimeCaller, err.Error())
	}

	req, err := http.NewRequest("POST", url, &b)

	if err != nil {
		return nil, errors.WithSCode(code.ErrInvokeTool, err.Error())
	}

	for k, v := range st.generateHeaders(writer) {
		req.Header.Add(k, v)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	return client.Do(req)
}

func (st *StableDiffusionTool) generateHeaders(writer *multipart.Writer) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + st.credentials,
		"Accept":        "image/*",
		"Content-Type":  writer.FormDataContentType(),
	}
}
