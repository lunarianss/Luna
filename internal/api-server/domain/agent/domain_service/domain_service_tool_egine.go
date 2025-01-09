package domain_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/core/tools/tool_registry"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	po_agent "github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/po_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type ToolEngine struct {
	tool              *biz_entity.ToolRuntimeConfiguration
	agentToolCallback any
	traceManager      any
	message           *po_entity.Message
	meta              *biz_entity.ToolEngineInvokeMeta
	bucket            string
	providerType      string
	*AgentDomain
}

func NewToolEngine(tool *biz_entity.ToolRuntimeConfiguration, message *po_entity.Message, providerType string, agentDomain *AgentDomain, bucket string) *ToolEngine {
	te := &ToolEngine{
		tool:         tool,
		message:      message,
		providerType: providerType,
		AgentDomain:  agentDomain,
		bucket:       bucket,
	}

	te.constructInvokeMeta()
	return te
}

func (te *ToolEngine) AgentInvoke(ctx context.Context, toolParameters string, userID, tenantID string, invokeFrom biz_entity.InvokeFrom) *biz_entity.ToolEngineInvokeMessage {
	response, err := te.invoke(ctx, toolParameters, userID)

	if err != nil {
		return te.handleInvokeError(err)
	}

	convertedMessages, err := NewToolFileMessageTransformer(te.AgentDomain, te.bucket).TransformToolInvokeMessages(ctx, response, userID, tenantID, te.message.ConversationID)

	if err != nil {
		return te.handleInvokeError(err)
	}

	binaryFiles := te.extractToolResponseBinary(convertedMessages)

	messageFiles, err := te.createMessageFiles(ctx, binaryFiles, te.message, invokeFrom, userID)

	if err != nil {
		return te.handleInvokeError(err)
	}

	plainText, err := te.convertToolResponseToString(convertedMessages)

	if err != nil {
		return te.handleInvokeError(err)
	}

	return &biz_entity.ToolEngineInvokeMessage{
		InvokeToolPrompt: plainText,
		MessageFiles:     messageFiles,
		ToolInvokeMeta:   te.meta,
	}
}

func (te *ToolEngine) createMessageFiles(ctx context.Context, toolMessages []*biz_entity.ToolInvokeMessageBinary, agentMessage *po_entity.Message, invokeFrom biz_entity.InvokeFrom, userID string) ([]*biz_entity.ToolEngineInvokeMessageFiles, error) {
	var (
		result        []*biz_entity.ToolEngineInvokeMessageFiles
		fileType      string
		createdByRole string
	)

	for _, toolMessage := range toolMessages {
		if invokeFrom == biz_entity.DebuggerInvoke || invokeFrom == biz_entity.ExploreInvoke {
			createdByRole = "account"
		} else {
			createdByRole = "end_user"
		}

		fileType = te.getFileType(toolMessage.MimeType)

		messageFile := &po_agent.MessageFile{
			MessageID:      agentMessage.ID,
			Type:           fileType,
			TransferMethod: "tool_file",
			BelongsTo:      "assistant",
			URL:            toolMessage.Url,
			UploadFileID:   toolMessage.ToolFileID,
			CreatedByRole:  createdByRole,
			CreatedBy:      userID,
		}

		messageFile, err := te.AgentDomain.CreateMessageFile(ctx, messageFile)
		if err != nil {
			return nil, err
		}

		result = append(result, &biz_entity.ToolEngineInvokeMessageFiles{
			MessageFile: messageFile.ID,
			SaveAs:      toolMessage.SaveAs,
		})
	}

	return result, nil
}

func (te *ToolEngine) invoke(ctx context.Context, toolParameters string, userID string) ([]*biz_entity.ToolInvokeMessage, error) {
	toolCaller := tool_registry.NewModelRegisterCaller(userID, te.tool)

	invokeMessages, err := toolCaller.Invoke(ctx, []byte(toolParameters))

	if err != nil {
		return nil, err
	}

	return invokeMessages, nil
}

func (te *ToolEngine) handleInvokeError(err error) *biz_entity.ToolEngineInvokeMessage {

	var invokeMessage = &biz_entity.ToolEngineInvokeMessage{}

	te.meta.Error = err.Error()
	invokeMessage.ToolInvokeMeta = te.meta
	invokeMessage.MessageFiles = make([]*biz_entity.ToolEngineInvokeMessageFiles, 0)

	if errors.IsCode(err, code.ErrInvokeTool) || errors.IsCode(err, code.ErrInvokeToolUnConvertAble) {
		invokeMessage.InvokeToolPrompt = fmt.Sprintf("tool invoke failed due to the follow json format error: %#+-v", err)
	} else if errors.IsCode(err, code.ErrToolParameter) {
		invokeMessage.InvokeToolPrompt = "tool parameters validation error: please check your tool parameters"
	} else if errors.IsCode(err, code.ErrNotFoundToolRegistry) {
		invokeMessage.InvokeToolPrompt = fmt.Sprintf("there is not a tool named %s", te.tool.Identity.Name)
	} else {
		invokeMessage.InvokeToolPrompt = fmt.Sprintf("tool failed due to the follow json format error: %#+-v", err)
	}

	return invokeMessage
}

func (te *ToolEngine) extractToolResponseBinary(toolMessages []*biz_entity.ToolInvokeMessage) []*biz_entity.ToolInvokeMessageBinary {
	var result []*biz_entity.ToolInvokeMessageBinary
	for _, toolMessage := range toolMessages {
		if toolMessage.Type == biz_entity.IMAGE || toolMessage.Type == biz_entity.IMAGE_LINK {
			result = append(result, &biz_entity.ToolInvokeMessageBinary{
				MimeType:   toolMessage.Meta["mime_type"].(string),
				Url:        toolMessage.Message.(string),
				SaveAs:     toolMessage.SaveAs,
				ToolFileID: toolMessage.ToolFileID,
			})
		}
	}

	return result
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
			result += fmt.Sprintf("result link: %s. please inform the user", toolMessage.Message)
		} else if toolMessage.Type == biz_entity.IMAGE || toolMessage.Type == biz_entity.IMAGE_LINK {
			result += "image has been created and sent to user already, you do not need to create it, just tell the user to check it now"
		} else if toolMessage.Type == biz_entity.JSON {
			message, ok := toolMessage.Message.([]byte)
			if !ok {
				return result, errors.WithCode(code.ErrInvokeToolUnConvertAble, "(json)invoke message %+v isn't convert to []byte", message)
			}
			result += fmt.Sprintf("tool response: %s", string(message))
		} else {
			result += fmt.Sprintf("tool response: %+v", toolMessage.Message)
		}
	}

	return result, nil
}

func (te *ToolEngine) getFileType(mimeType string) string {
	var fileType string
	if strings.Contains(mimeType, "image") {
		fileType = "image"
	} else if strings.Contains(mimeType, "video") {
		fileType = "video"
	} else if strings.Contains(mimeType, "audio") {
		fileType = "audio"
	} else if strings.Contains(mimeType, "text") || strings.Contains(mimeType, "pdf") {
		fileType = "document"
	} else {
		fileType = "custom"
	}
	return fileType
}
