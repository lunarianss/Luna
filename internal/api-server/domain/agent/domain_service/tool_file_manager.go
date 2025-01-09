package domain_service

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/po_entity"
	"github.com/lunarianss/Luna/internal/infrastructure/util"
)

type ToolFileManager struct {
	*AgentDomain
}

func NewToolFileManager(agentDomain *AgentDomain) *ToolFileManager {
	return &ToolFileManager{}
}

func (tf *ToolFileManager) SignFile(toolFileID string, extension string, secretKey string, baseUrl string) (string, error) {
	filePreviewURL := fmt.Sprintf("%s/files/tools/%s%s", baseUrl, toolFileID, extension)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	nonce, err := util.GenerateNonce(16)

	if err != nil {
		return "", nil
	}

	dataToSign := fmt.Sprintf("file-preview|%s|%s|%s", toolFileID, timestamp, nonce)

	sign := util.GenerateHMACSignature(dataToSign, secretKey)

	encodedSign := base64.URLEncoding.EncodeToString(sign)

	return fmt.Sprintf("%s?timestamp=%s&nonce=%s&sign=%s", filePreviewURL, timestamp, nonce, encodedSign), nil
}

func (tf *ToolFileManager) VerifyFile(fileID string, timestamp int64, nonce string, sign string, secretKey string, timeout int64) bool {
	dataToSign := fmt.Sprintf("file-preview|%s|%d|%s", fileID, timestamp, nonce)

	recalculatedSign := util.GenerateHMACSignature(dataToSign, secretKey)

	encodedRecalculatedSign := base64.URLEncoding.EncodeToString(recalculatedSign)

	if sign != encodedRecalculatedSign {
		return false
	}

	currentTime := time.Now().Unix()

	return currentTime-timestamp > timeout
}

func (tf *ToolFileManager) CreateFileByRaw(ctx context.Context, userID, tenantID, conversationID string, fileBinary []byte, mimeType string) (*po_entity.ToolFile, error) {

	var extension string

	extensions, err := mime.ExtensionsByType(mimeType)

	if err != nil || len(extensions) == 0 {
		extension = ".bin"
	}

	extension = extensions[0]

	filename := fmt.Sprintf("%s%s", uuid.NewString(), extension)
	filepath := fmt.Sprintf("tools/%s/%s", tenantID, filename)

	// todo storage save
	toolFile := &po_entity.ToolFile{
		UserID:         userID,
		TenantID:       tenantID,
		ConversationID: conversationID,
		FileKey:        filepath,
		MimeType:       mimeType,
		Size:           len(fileBinary),
	}

	toolFile, err = tf.AgentRepo.CreateToolFile(ctx, toolFile)

	if err != nil {
		return nil, err
	}

	return toolFile, nil
}
