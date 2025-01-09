package domain_service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type ToolFileMessageTransformer struct {
	toolMessages []*biz_entity.ToolInvokeMessage
	bucket       string
	*AgentDomain
}

func NewToolFileMessageTransformer(agentDomain *AgentDomain, bucket string) *ToolFileMessageTransformer {
	return &ToolFileMessageTransformer{
		AgentDomain:  agentDomain,
		toolMessages: make([]*biz_entity.ToolInvokeMessage, 0),
		bucket:       bucket,
	}
}

func (ttr *ToolFileMessageTransformer) TransformToolInvokeMessages(ctx context.Context, messages []*biz_entity.ToolInvokeMessage, userID, tenantID, conversationID string) ([]*biz_entity.ToolInvokeMessage, error) {

	for _, toolMessage := range messages {
		if toolMessage.Type == biz_entity.BLOB {
			if err := ttr.handleBlobMessage(ctx, toolMessage, userID, tenantID, conversationID); err != nil {
				return nil, err
			}
		}
	}

	return ttr.toolMessages, nil
}

func (ttr *ToolFileMessageTransformer) handleBlobMessage(ctx context.Context, toolMessage *biz_entity.ToolInvokeMessage, userID, tenantID, conversationID string) error {
	var (
		mimeType      any
		mimeTypeStr   string
		isExist       bool
		isConvertAble bool
		messageMeta   map[string]any
	)

	mimeType, isExist = toolMessage.Meta["mime_type"]
	if !isExist {
		mimeTypeStr = "octet/stream"
	}

	mimeTypeStr, isConvertAble = mimeType.(string)

	if !isConvertAble {
		return errors.WithCode(code.ErrInvokeToolUnConvertAble, "mimetype isn't able to convert to string")
	}

	meta := toolMessage.Meta

	if meta == nil {
		messageMeta = make(map[string]any)
	} else {
		messageMeta = meta
	}

	toolResByte, ok := toolMessage.Message.([]byte)

	if !ok {
		return errors.WithCode(code.ErrInvokeToolUnConvertAble, "tool blob response is not convert to []byte, actual it's %T", toolMessage.Message)
	}

	toolFile, err := NewToolFileManager(ttr.AgentDomain, ttr.bucket).CreateFileByRaw(ctx, userID, tenantID, conversationID, toolResByte, mimeTypeStr)

	if err != nil {
		return err
	}

	extension := filepath.Ext(toolFile.FileKey)

	imgUrl := ttr.getToolFileUrl(toolFile.ID, extension)

	if strings.Contains(mimeTypeStr, "image") {
		ttr.toolMessages = append(ttr.toolMessages, &biz_entity.ToolInvokeMessage{
			Type:       biz_entity.IMAGE_LINK,
			Message:    imgUrl,
			SaveAs:     toolMessage.SaveAs,
			Meta:       messageMeta,
			ToolFileID: toolMessage.ToolFileID,
		})
	} else {
		ttr.toolMessages = append(ttr.toolMessages, &biz_entity.ToolInvokeMessage{
			Type:       biz_entity.LINK,
			Message:    imgUrl,
			SaveAs:     toolMessage.SaveAs,
			Meta:       messageMeta,
			ToolFileID: toolMessage.ToolFileID,
		})
	}
	return nil
}

func (ttr *ToolFileMessageTransformer) getToolFileUrl(toolFileID, extension string) string {
	return fmt.Sprintf("/files/tools/%s%s", toolFileID, extension)
}
