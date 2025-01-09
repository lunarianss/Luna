package service

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/api-server/config"
	"github.com/lunarianss/Luna/internal/api-server/core/storage"
	storage_interface "github.com/lunarianss/Luna/internal/api-server/core/storage/interface"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
	agentDomain "github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
	dto "github.com/lunarianss/Luna/internal/api-server/dto/file"
	"github.com/lunarianss/Luna/internal/infrastructure/code"
)

type FileService struct {
	agentDomain *agentDomain.AgentDomain
	config      *config.Config
}

func NewFileService(agentDomain *agentDomain.AgentDomain, config *config.Config) *FileService {
	return &FileService{
		config:      config,
		agentDomain: agentDomain,
	}
}

func (ts *FileService) PreviewFile(ctx context.Context, fileID string, args *dto.PreviewFileQuery) error {

	if isPass := domain_service.NewToolFileManager(nil, "").VerifyFile(fileID, args.Timestamp, args.Nonce, args.Sign, ts.config.SystemOptions.SecretKey, ts.config.SystemOptions.FileTimeout); !isPass {
		return errors.WithSCode(code.ErrForbidden, "verify file sign failed")
	}

	storage, err := storage.NewStorage(ctx, ts.config.MinioOptions.Bucket, storage_interface.MINIO)

	if err != nil {
		return err
	}

	toolFile, err := ts.agentDomain.AgentRepo.GetToolFileByID(ctx, fileID)

	if err != nil {
		return err
	}

	streams, err := storage.LoadStream(ctx, toolFile.FileKey)

	if err != nil {
		return err
	}

	g, ok := ctx.(*gin.Context)

	if !ok {
		return errors.WithSCode(code.ErrRunTimeCaller, "ctx must be *gin.Context")
	}

	if toolFile.Size > 0 {
		g.Header("Content-Length", strconv.Itoa(toolFile.Size))

	}
	if args.AsAttachment {
		g.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", toolFile.Name))
	}

	g.Stream(func(w io.Writer) bool {
		for stream := range streams {
			_, err := w.Write(stream)
			if err != nil {
				return false
			}
		}
		return false
	})

	return nil
}
