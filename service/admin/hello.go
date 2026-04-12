package admin

import (
	"context"

	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"go.uber.org/zap"
)

func (s *Service) Hello(ctx context.Context, adminUser *common.AdminUser, req *dto.HelloReq) (*dto.HelloResp, common.Errno) {
	msg, err := s.adminUser.Hello(ctx, &do.Hello{})
	if err != nil {
		logger.Error("Hello error", zap.Error(err), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return &dto.HelloResp{Hello: msg, World: "world"}, common.OK
}
