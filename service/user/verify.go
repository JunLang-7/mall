package user

import (
	"context"
	"encoding/json"

	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/JunLang-7/mall/utils/tools"
	"github.com/wenlng/go-captcha/v2/slide"
	"go.uber.org/zap"
)

// GetSlideCaptcha 获取滑动验证码
func (s *Service) GetSlideCaptcha(ctx context.Context) (*dto.GetVerifyCaptchaResp, common.Errno) {
	// 生成滑动验证码
	captData, err := s.captcha.Generate()
	if err != nil {
		logger.Error("GetSlideCaptcha Generate error", zap.Error(err))
		return nil, *common.ServerErr.WithErr(err)
	}
	// 获取验证码数据
	dotData := captData.GetData()
	if dotData == nil {
		logger.Error("GetSlideCaptcha captcha data is nil")
		return nil, *common.ServerErr.WithMsg("GetData is nil")
	}
	// 将验证码数据转换为 JSON 字符串
	dots, err := json.Marshal(dotData)
	if err != nil {
		logger.Error("GetSlideCaptcha json.Marshal error", zap.Error(err))
		return nil, *common.ServerErr.WithErr(err)
	}

	// 将验证码图片转换为 Base64 编码字符串
	mBs64Data, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		logger.Error("GetSlideCaptcha GetMasterImage error", zap.Error(err))
		return nil, *common.ServerErr.WithErr(err)
	}
	// 将验证码缺口图片转换为 Base64 编码字符串
	tBs64Data, err := captData.GetTileImage().ToBase64()
	if err != nil {
		logger.Error("GetSlideCaptcha GetTileImage error", zap.Error(err))
		return nil, *common.ServerErr.WithErr(err)
	}

	// 生成一个唯一的验证码 key，并将验证码数据存储到 Redis 中，设置过期时间
	key := tools.UUIDHex()
	err = s.verify.SetCaptchaKey(ctx, key, string(dots), consts.ExpireLoginTime)
	if err != nil {
		logger.Error("GetSlideCaptcha SetCaptchaKey error", zap.Error(err))
		return nil, *common.RedisErr.WithErr(err)
	}

	return &dto.GetVerifyCaptchaResp{
		Key:            key,
		ImageBs64:      mBs64Data,
		TitleImageBs64: tBs64Data,
		TitleHeight:    dotData.Height,
		TitleWidth:     dotData.Width,
		TitleX:         dotData.X,
		TitleY:         dotData.Y,
		Expire:         int64(consts.ExpireLoginTime.Seconds()),
	}, common.OK
}

// CheckSlideCaptcha 校验滑动验证码
func (s *Service) CheckSlideCaptcha(ctx context.Context, req *dto.CheckCaptchaReq) (*dto.CheckCaptchaResp, common.Errno) {
	captData, err := s.verify.GetCaptchaKey(ctx, req.Key)
	if err != nil {
		logger.Error("CheckSlideCaptcha GetCaptchaKey error", zap.Error(err))
		return nil, *common.ServerErr.WithErr(err)
	}
	if captData == "" {
		return nil, *common.ParamErr.WithMsg("captcha is expired")
	}
	var dot slide.Block
	if err = json.Unmarshal([]byte(captData), &dot); err != nil {
		logger.Error("CheckSlideCaptcha json.Unmarshal error", zap.Error(err))
		return nil, common.InvalidCaptchaErr
	}
	ok := slide.CheckPoint(int64(req.SlideX), int64(req.SlideY), int64(dot.X), int64(dot.Y), 5)
	if !ok {
		logger.Error("CheckSlideCaptcha slide.CheckPoint error", zap.Error(err))
		return nil, common.InvalidCaptchaErr
	}
	ticket := tools.UUIDHex()
	err = s.verify.SetCaptchaTicket(ctx, ticket, req.Key, consts.ExpireTicketTime)
	if err != nil {
		logger.Error("CheckSlideCaptcha SetCaptchaTicket error", zap.Error(err))
		return nil, *common.RedisErr.WithErr(err)
	}
	return &dto.CheckCaptchaResp{
		Ticket: ticket,
		Expire: int64(consts.ExpireTicketTime.Seconds()),
	}, common.OK
}
