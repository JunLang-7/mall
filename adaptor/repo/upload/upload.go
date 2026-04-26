package upload

import (
	"context"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/service/do"
	"gorm.io/gorm"
)

type IUploadFile interface {
	CreateUploadFile(ctx context.Context, fileList []do.AddUploadFile) error
	DeleteUploadFile(ctx context.Context, strings []string) error
}

type UploadFile struct {
	db *gorm.DB
}

func NewUploadFile(adaptor adaptor.IAdaptor) *UploadFile {
	return &UploadFile{
		db: adaptor.GetDB(),
	}
}

// CreateUploadFile 创建上传文件记录
func (u *UploadFile) CreateUploadFile(ctx context.Context, fileList []do.AddUploadFile) error {
	var addList []model.ResourceUploadFile
	for _, file := range fileList {
		addList = append(addList, model.ResourceUploadFile{
			Scene:          file.Scene,
			FileKey:        file.FileKey,
			FileName:       file.FileName,
			FileSize:       file.FileSize,
			FileType:       file.FileType,
			UploadClientIP: file.ClientIP,
			UserID:         file.UserID,
			UserType:       file.UserType,
		})
	}
	return u.db.WithContext(ctx).CreateInBatches(&addList, 100).Error
}

// DeleteUploadFile 删除上传文件记录
func (u *UploadFile) DeleteUploadFile(ctx context.Context, fileList []string) error {
	if len(fileList) == 0 {
		return nil
	}
	return u.db.WithContext(ctx).Where("file_key IN ?", fileList).Delete(&model.ResourceUploadFile{}).Error
}
