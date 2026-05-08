package user

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/JunLang-7/mall/utils/tools"
	goredis "github.com/go-redis/redis"
	"go.uber.org/zap"
)

const (
	qrStatePending   = "pending"
	qrStateScanned   = "scanned"
	qrStateConfirmed = "confirmed"
	qrStateExpired   = "expired"
)

type qrScene struct {
	State   string `json:"state"`
	Purpose string `json:"purpose"`
	Code    string `json:"code,omitempty"`
}

func (s *Service) GetWechatQrCode(ctx context.Context, purpose string) (*dto.WechatQrCodeResp, common.Errno) {
	sceneToken := tools.UUIDHex()
	scene := &qrScene{State: qrStatePending, Purpose: purpose}
	data, _ := json.Marshal(scene)
	if err := s.qrcode.SetScene(ctx, sceneToken, data, consts.ExpireTicketTime); err != nil {
		logger.Error("GetWechatQrCode SetScene error", zap.Error(err))
		return nil, *common.RedisErr.WithErr(err)
	}
	return &dto.WechatQrCodeResp{
		ExpireIn:   int64(consts.ExpireTicketTime.Seconds()),
		SceneToken: sceneToken,
	}, common.OK
}

func (s *Service) GetWechatQrCodeStatus(ctx context.Context, sceneToken string) (*dto.WechatQrCodeStatusResp, common.Errno) {
	data, err := s.qrcode.GetScene(ctx, sceneToken)
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return &dto.WechatQrCodeStatusResp{State: qrStateExpired, Message: "expired"}, common.OK
		}
		logger.Error("GetWechatQrCodeStatus GetScene error", zap.Error(err))
		return nil, *common.RedisErr.WithErr(err)
	}
	scene := &qrScene{}
	if err = json.Unmarshal(data, scene); err != nil {
		return nil, *common.ServerErr.WithErr(err)
	}
	switch scene.State {
	case qrStateConfirmed:
		return &dto.WechatQrCodeStatusResp{State: qrStateConfirmed, Purpose: scene.Purpose, Message: "login success"}, common.OK
	case qrStateScanned:
		return &dto.WechatQrCodeStatusResp{State: qrStateScanned, Purpose: scene.Purpose, Message: "scanned"}, common.OK
	default:
		return &dto.WechatQrCodeStatusResp{State: qrStatePending, Purpose: scene.Purpose, Message: "waiting"}, common.OK
	}
}

func (s *Service) ConfirmWechatScan(ctx context.Context, sceneToken, code string) (*dto.WechatQrCodeStatusResp, common.Errno) {
	data, err := s.qrcode.GetScene(ctx, sceneToken)
	if err != nil {
		logger.Error("ConfirmWechatScan GetScene error", zap.Error(err))
		return nil, *common.ParamErr.WithMsg("invalid scene_token")
	}
	scene := &qrScene{}
	_ = json.Unmarshal(data, scene)

	// 模拟微信用户 openid/unionid
	mockUnionID := "mock_union_" + tools.UUIDHex()[:8]
	user, errno := s.findOrCreateWechatUser(ctx, mockUnionID)
	if !errno.IsOK() {
		return nil, errno
	}

	// 生成登录 token
	resp, errno := s.handleCustomerLogin(ctx, user)
	if !errno.IsOK() {
		return nil, errno
	}

	// 更新 Redis 场景状态
	scene.State = qrStateConfirmed
	scene.Code = code
	newData, _ := json.Marshal(scene)
	_ = s.qrcode.SetScene(ctx, sceneToken, newData, consts.ExpireTicketTime)

	return &dto.WechatQrCodeStatusResp{
		State:   qrStateConfirmed,
		Purpose: scene.Purpose,
		Message: "login success",
		Token:   resp.Token,
		UserInfo: resp.UserInfo,
	}, common.OK
}

func (s *Service) findOrCreateWechatUser(ctx context.Context, unionID string) (*model.User, common.Errno) {
	wxUser, err := s.userRepo.GetWechatUserByUnionID(ctx, unionID)
	if err == nil && wxUser != nil {
		user, err := s.userRepo.GetUserByID(ctx, wxUser.UserID)
		if err == nil && user.Status == consts.IsEnable {
			return user, common.OK
		}
	}
	now := time.Now()
	nickname := "微信用户" + tools.UUIDHex()[:4]
	user := &model.User{
		NickName:    nickname,
		Sex:         0,
		Status:      consts.IsEnable,
		CreateAt:    now,
		UpdateAt:    now,
		LastLoginAt: now,
	}
	wu := &model.WechatUser{
		UnionID:  unionID,
		NickName: nickname,
		CreateAt: now,
		UpdateAt: now,
	}
	if err := s.userRepo.CreateUserWithWechatUser(ctx, user, wu); err != nil {
		logger.Error("findOrCreateWechatUser error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return user, common.OK
}
