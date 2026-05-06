package user

import (
	"context"
	"time"

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
				errno := s.CancelTimeoutOrders(ctx)
				if !errno.IsOK() {
					logger.Error("CancelTimeoutOrders error", zap.Any("errno", errno))
				}
			}
		}
	}()
}
