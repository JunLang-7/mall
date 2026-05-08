package user

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/utils/logger"
	"go.uber.org/zap"
)

func (s *Service) StartBackgroundTasks(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if errno := s.CancelTimeoutOrders(ctx); !errno.IsOK() {
					logger.Error("CancelTimeoutOrders error", zap.Any("errno", errno))
				}
				if errno := s.AutoReceiveOrders(ctx); !errno.IsOK() {
					logger.Error("AutoReceiveOrders error", zap.Any("errno", errno))
				}
			}
		}
	}()
}

func (s *Service) AutoReceiveOrders(ctx context.Context) common.Errno {
	autoDays := s.conf.Order.AutoReceiveDays
	if autoDays <= 0 {
		return common.OK
	}
	shippedBefore := time.Now().Add(-time.Duration(autoDays) * 24 * time.Hour).UnixMilli()
	now := time.Now().UnixMilli()
	if err := s.order.AutoReceiveOrders(ctx, shippedBefore, now); err != nil {
		logger.Error("AutoReceiveOrders error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}
