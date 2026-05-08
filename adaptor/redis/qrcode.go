package redis

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/go-redis/redis"
)

type IQrCode interface {
	SetScene(ctx context.Context, sceneToken string, data []byte, ttl time.Duration) error
	GetScene(ctx context.Context, sceneToken string) ([]byte, error)
	DeleteScene(ctx context.Context, sceneToken string) error
}

type QrCode struct {
	rds *redis.Client
}

func NewQrCode(adaptor adaptor.IAdaptor) *QrCode {
	return &QrCode{rds: adaptor.GetRedis()}
}

func (q *QrCode) SetScene(ctx context.Context, sceneToken string, data []byte, ttl time.Duration) error {
	return q.rds.Set(sceneKey(sceneToken), data, ttl).Err()
}

func (q *QrCode) GetScene(ctx context.Context, sceneToken string) ([]byte, error) {
	return q.rds.Get(sceneKey(sceneToken)).Bytes()
}

func (q *QrCode) DeleteScene(ctx context.Context, sceneToken string) error {
	return q.rds.Del(sceneKey(sceneToken)).Err()
}

func sceneKey(sceneToken string) string {
	return "mall:wechat:qrcode:" + sceneToken
}
