package domain_service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/core/tools/tool_registry"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type ToolEngine struct {
	tool              *biz_entity.ToolRuntimeConfiguration
	agentToolCallback any
	traceManager      any
	message           *po_entity.Message
	meta              *biz_entity.ToolEngineInvokeMeta
	providerType      string
}

func NewToolEngine(tool *biz_entity.ToolRuntimeConfiguration, message *po_entity.Message, providerType string) *ToolEngine {
	te := &ToolEngine{
		tool:         tool,
		message:      message,
		providerType: providerType,
	}

	te.constructInvokeMeta()
	return te
}

func (te *ToolEngine) AgentInvoke(ctx context.Context, toolParameters map[string]any, userID, tenantID string, invokeFrom biz_entity.InvokeFrom) *biz_entity.ToolEngineInvokeMessage {
	response, err := te.invoke(ctx, toolParameters, userID)

	if err != nil {
		return te.handleInvokeError(err)
	}

	convertedMessages, err := NewToolFileMessageTransformer().TransformToolInvokeMessages(response, userID, tenantID, te.message.ConversationID)

	if err != nil {
		return te.handleInvokeError(err)
	}

	plainText, err := te.convertToolResponseToString(convertedMessages)

	if err != nil {
		return te.handleInvokeError(err)
	}

	return &biz_entity.ToolEngineInvokeMessage{
		InvokeToolPrompt: plainText,
		MessageFiles:     nil,
		ToolInvokeMeta:   te.meta,
	}
}

func (te *ToolEngine) invoke(ctx context.Context, toolParameters map[string]any, userID string) ([]*biz_entity.ToolInvokeMessage, error) {

	toolCaller := tool_registry.NewModelRegisterCaller(userID, te.tool)

	parameterByte, err := json.Marshal(toolParameters)

	if err != nil {
		return nil, errors.WithSCode(code.ErrInvokeTool, err.Error())
	}

	invokeMessages, err := toolCaller.Invoke(ctx, parameterByte)

	if err != nil {
		return nil, errors.WithSCode(code.ErrInvokeTool, err.Error())
	}

	return invokeMessages, nil
}

func (te *ToolEngine) handleInvokeError(err error) *biz_entity.ToolEngineInvokeMessage {

	var invokeMessage = &biz_entity.ToolEngineInvokeMessage{}

	te.meta.Error = err.Error()
	invokeMessage.ToolInvokeMeta = te.meta
	invokeMessage.MessageFiles = make([]*biz_entity.ToolEngineInvokeMessageFiles, 0)

	if errors.IsCode(err, code.ErrInvokeTool) || errors.IsCode(err, code.ErrInvokeToolUnConvertAble) {
		invokeMessage.InvokeToolPrompt = fmt.Sprintf("tool invoke error: %s", err.Error())
	}

	return invokeMessage
}

func (te *ToolEngine) constructInvokeMeta() {
	te.meta = &biz_entity.ToolEngineInvokeMeta{
		TimeCost: 0,
		ToolConfig: map[string]any{
			"tool_name":          te.tool.Identity.Name,
			"tool_provider":      te.tool.Identity.Provider,
			"tool_provider_type": te.providerType,
			"tool_parameters":    te.tool.RuntimeParameters,
			"tool_icon":          te.tool.Identity.Icon,
		},
	}
}

func (te *ToolEngine) convertToolResponseToString(toolMessages []*biz_entity.ToolInvokeMessage) (string, error) {

	var result string

	for _, toolMessage := range toolMessages {
		if toolMessage.Type == biz_entity.TEXT {
			message, ok := toolMessage.Message.(string)
			if !ok {
				return result, errors.WithCode(code.ErrInvokeToolUnConvertAble, "(text)invoke message %+v isn't convert to string", toolMessage.Message)
			}
			result += message
		} else if toolMessage.Type == biz_entity.LINK {
			result += fmt.Sprintf("result link: %s. please tell user to check it", toolMessage.Message)
		} else if toolMessage.Type == biz_entity.IMAGE || toolMessage.Type == biz_entity.IMAGE_LINK {
			result += "image has been created and sent to user already, you do not need to create it, just tell the user to check it now"
		} else if toolMessage.Type == biz_entity.JSON {
			message, ok := toolMessage.Message.([]byte)
			if !ok {
				return result, errors.WithCode(code.ErrInvokeToolUnConvertAble, "(json)invoke message %+v isn't convert to []byte", message)
			}
			result += fmt.Sprintf("tool json response: %s", string(message))
		} else {
			result += fmt.Sprintf("tool json response: %+v", toolMessage.Message)
		}
	}

	return result, nil
}
