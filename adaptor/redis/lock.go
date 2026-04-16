package redis

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/go-redis/redis"
)

type ILocker interface {
	GetLock(ctx context.Context, key string) (bool, error)
	AwaitLock(ctx context.Context, key string, duration time.Duration) error
}

type Locker struct {
	redis *redis.Client
}

func NewLocker(adaptor adaptor.IAdaptor) ILocker {
	return &Locker{redis: adaptor.GetRedis()}
}

// GetLock 尝试获取锁，成功返回 true，失败返回 false
func (l *Locker) GetLock(ctx context.Context, key string) (bool, error) {
	return l.redis.SetNX(key, "1", time.Second*5).Result()
}

// AwaitLock 等待锁释放，直到获取到锁或者超时
func (l *Locker) AwaitLock(ctx context.Context, key string, duration time.Duration) error {
	deadline := time.Now().Add(duration)
	for {
		exist, err := l.redis.Exists(key).Result()
		if err != nil {
			return err
		}
		if exist == 0 {
			return nil
		}
		if time.Now().After(deadline) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}
