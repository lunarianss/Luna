package domain_service

import (
	"strings"

	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type ToolFileMessageTransformer struct {
	toolMessages []*biz_entity.ToolInvokeMessage
}

func NewToolFileMessageTransformer() *ToolFileMessageTransformer {
	return &ToolFileMessageTransformer{
		toolMessages: make([]*biz_entity.ToolInvokeMessage, 0),
	}
}

func (ttr *ToolFileMessageTransformer) TransformToolInvokeMessages(messages []*biz_entity.ToolInvokeMessage, userID, tenantID, conversationID string) ([]*biz_entity.ToolInvokeMessage, error) {

	for _, toolMessage := range messages {
		if toolMessage.Type == biz_entity.BLOB {
			if err := ttr.handleBlobMessage(toolMessage); err != nil {
				return nil, err
			}
		}
	}

	return ttr.toolMessages, nil
}

func (ttr *ToolFileMessageTransformer) handleBlobMessage(toolMessage *biz_entity.ToolInvokeMessage) error {
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

	if strings.Contains(mimeTypeStr, "image") {
		ttr.toolMessages = append(ttr.toolMessages, &biz_entity.ToolInvokeMessage{
			Type:    biz_entity.IMAGE_LINK,
			Message: "http://pic.com/file",
			SaveAs:  toolMessage.SaveAs,
			Meta:    messageMeta,
		})
	} else {
		ttr.toolMessages = append(ttr.toolMessages, &biz_entity.ToolInvokeMessage{
			Type:    biz_entity.LINK,
			Message: "http://pic.com/file",
			SaveAs:  toolMessage.SaveAs,
			Meta:    messageMeta,
		})
	}

	return nil
}
