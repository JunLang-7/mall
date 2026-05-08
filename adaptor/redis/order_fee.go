package redis

import (
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/go-redis/redis"
)

type IOrderFee interface {
	SetOrderFee(key string, data []byte, expire time.Duration) error
	GetOrderFee(key string) ([]byte, error)
}

type OrderFee struct {
	rds *redis.Client
}

func NewOrderFee(adaptor adaptor.IAdaptor) *OrderFee {
	return &OrderFee{
		rds: adaptor.GetRedis(),
	}
}

func (o *OrderFee) SetOrderFee(key string, data []byte, expire time.Duration) error {
	return o.rds.Set(key, data, expire).Err()
}

func (o *OrderFee) GetOrderFee(key string) ([]byte, error) {
	return o.rds.Get(key).Bytes()
}
