package rpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/config"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/utils/tools"
	"github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"gorm.io/gorm"
)

var ErrInvalidStorageParam = errors.New("invalid storage param")

type IStorage interface {
	GetTempSecret(ctx context.Context, req *do.GetTempSecret) (*do.TempSecret, error)
	GetPreviewUrl(ctx context.Context, req *do.GetPreviewUrl) (map[string]string, error)
	DeleteFile(ctx context.Context, req *do.DeleteFile) error
}

type Storage struct {
	db   *gorm.DB
	conf config.Config
}

func NewStorage(adaptor adaptor.IAdaptor) *Storage {
	return &Storage{
		db:   adaptor.GetDB(),
		conf: *adaptor.GetConf(),
	}
}

func (s *Storage) GetTempSecret(ctx context.Context, req *do.GetTempSecret) (*do.TempSecret, error) {
	client, err := s.getPreviewClient(ctx)
	if err != nil {
		return nil, err
	}
	// 根据业务场景和文件类型等参数判断上传路径，生成文件key
	path, ok := s.conf.Storage.Buckets.Paths[req.Scene]
	if !ok || strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("%w: scene %s", ErrInvalidStorageParam, req.Scene)
	}
	idx := strings.LastIndex(req.FileName, ".")
	if idx < 0 || idx == len(req.FileName)-1 {
		return nil, fmt.Errorf("%w: file name %s", ErrInvalidStorageParam, req.FileName)
	}
	postFix := req.FileName[idx+1:]
	objectPath := strings.Trim(path, "/")
	fileKey := objectPath + "/" + tools.UUIDHex() + "." + postFix

	// 获取临时密钥
	stClient := sts.NewClient(s.conf.Storage.SecretID, s.conf.Storage.SecretKey, nil)
	res, err := stClient.GetCredential(s.getCredentialOptions(objectPath))
	if err != nil {
		return nil, err
	}

	// 生成预签名URL
	preUrl, err := client.Object.GetPresignedURL(ctx, http.MethodGet, fileKey, s.conf.Storage.SecretID, s.conf.Storage.SecretKey, time.Hour, nil)
	if err != nil {
		return nil, err
	}

	return &do.TempSecret{
		SecretID:      res.Credentials.TmpSecretID,
		SecretKey:     res.Credentials.TmpSecretKey,
		SecurityToken: res.Credentials.SessionToken,
		Bucket:        s.conf.Storage.Buckets.BucketName,
		Region:        s.conf.Storage.Buckets.Region,
		Key:           fileKey,
		FileURL:       preUrl.String(),
		ExpireTime:    int64(res.ExpiredTime),
		StartTime:     int64(res.StartTime),
	}, nil
}

// getCredentialOptions 根据上传路径生成临时密钥的权限策略
func (s *Storage) getCredentialOptions(path string) *sts.CredentialOptions {
	return &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          s.conf.Storage.Buckets.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://cloud.tencent.com/document/product/436/31923
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						// 这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						// 存储桶的命名格式为 BucketName-APPID，此处填写的 bucket 必须为此格式
						fmt.Sprintf("qcs::cos:%s:uid/%s:%s/%s/*",
							s.conf.Storage.Buckets.Region,
							s.conf.Storage.AppID,
							s.conf.Storage.Buckets.BucketName,
							path,
						),
					},
					// 开始构建生效条件 condition
					// 关于 condition 的详细设置规则和COS支持的condition类型可以参考https://cloud.tencent.com/document/product/436/71306
					Condition: map[string]map[string]interface{}{},
				},
			},
		},
	}
}

// getPreviewClient 获取预览文件的COS客户端，预览文件需要使用CDN域名生成URL
func (s *Storage) getPreviewClient(_ context.Context) (*cos.Client, error) {
	u, err := url.Parse(s.conf.Storage.Buckets.CdnDomain)
	if err != nil {
		return nil, err
	}
	baseUrl := &cos.BaseURL{BucketURL: u}
	previewClient := cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  s.conf.Storage.SecretID,
			SecretKey: s.conf.Storage.SecretKey,
		},
	})
	return previewClient, nil
}

func (s *Storage) GetPreviewUrl(ctx context.Context, req *do.GetPreviewUrl) (map[string]string, error) {
	client, err := s.getPreviewClient(ctx)
	if err != nil {
		return nil, err
	}
	fileKeyMap := make(map[string]string)
	for _, fileKey := range req.Keys {
		pre, err := client.Object.GetPresignedURL(ctx, http.MethodGet, fileKey, s.conf.Storage.SecretID, s.conf.Storage.SecretKey, time.Hour, nil)
		if err != nil {
			return nil, err
		}
		fileKeyMap[fileKey] = pre.String()
	}
	return fileKeyMap, nil
}

func (s *Storage) DeleteFile(ctx context.Context, req *do.DeleteFile) error {
	//TODO implement me
	panic("implement me")
}
