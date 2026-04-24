package storage

import (
	"context"
	"errors"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/upload"
	"github.com/JunLang-7/mall/adaptor/rpc"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"go.uber.org/zap"
)

type Service struct {
	cos  rpc.IStorage
	repo upload.IUploadFile
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		cos:  rpc.NewStorage(adaptor),
		repo: upload.NewUploadFile(adaptor),
	}
}

// GetTempSecret 获取对象存储临时密钥
func (s *Service) GetTempSecret(ctx context.Context, req *dto.GetTempSecretReq) (*dto.TempSecretResp, common.Errno) {
	secret, err := s.cos.GetTempSecret(ctx, &do.GetTempSecret{
		Scene:    req.Scene,
		FileName: req.FileName,
		FileSize: req.FileSize,
		FileType: req.FileType,
		ClientIP: req.ClientIP,
	})
	if err != nil {
		if errors.Is(err, rpc.ErrInvalidStorageParam) {
			return nil, *common.ParamErr.WithErr(err)
		}
		logger.Error("GetTempSecret GetTempSecret error", zap.Any("req", req), zap.Error(err))
		return nil, *common.ServerErr.WithErr(err)
	}
	return &dto.TempSecretResp{
		SecretID:      secret.SecretID,
		SecretKey:     secret.SecretKey,
		SecurityToken: secret.SecurityToken,
		ExpiredTime:   secret.ExpireTime,
		Bucket:        secret.Bucket,
		Region:        secret.Region,
		Key:           secret.Key,
		FileURL:       secret.FileURL,
	}, common.OK
}
