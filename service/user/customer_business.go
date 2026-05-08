package user

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/JunLang-7/mall/utils/secure"
	"github.com/JunLang-7/mall/utils/tools"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *Service) ListCustomerCourse(ctx context.Context, userID int64, req *dto.CourseListReq) (*dto.CourseListResp, common.Errno) {
	list, total, err := s.course.ListCourse(ctx, &do.CourseListReq{
		ID: req.ID, NameKW: req.NameKW, Pager: req.Pager,
	})
	if err != nil {
		logger.Error("ListCustomerCourse error", zap.Error(err), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	purchased := s.purchasedMap(ctx, userID)
	courses := make([]*dto.CourseDto, 0, len(list))
	for _, course := range list {
		courses = append(courses, s.courseDto(ctx, course, purchased[course.ID]))
	}
	return &dto.CourseListResp{Pager: req.Pager, List: courses, Total: total}, common.OK
}

func (s *Service) GetCustomerCourseDetail(ctx context.Context, userID int64, id int64) (*dto.CustomerCourseDetailDto, common.Errno) {
	course, err := s.course.GetCourseInfo(ctx, &do.CourseInfoReq{ID: id, Status: consts.IsEnable})
	if err != nil {
		logger.Error("GetCustomerCourseDetail error", zap.Error(err), zap.Int64("id", id))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	purchased, err := s.userRepo.HasPurchasedCourse(ctx, &do.HasPurchasedReq{UserID: userID, CourseID: id})
	if err != nil {
		logger.Error("GetCustomerCourseDetail HasPurchasedCourse error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	catalogs, totalDuration, lessonCount, err := s.courseCatalogs(ctx, id, purchased)
	if err != nil {
		logger.Error("GetCustomerCourseDetail courseCatalogs error", zap.Error(err), zap.Int64("id", id))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	dtoCourse := s.courseDto(ctx, course, purchased)
	return &dto.CustomerCourseDetailDto{
		CourseDto:     *dtoCourse,
		TotalDuration: totalDuration,
		LessonCount:   lessonCount,
		Catalogs:      catalogs,
		HasPurchased:  dtoCourse.HasPurchased,
	}, common.OK
}

func (s *Service) GetLessonInfo(ctx context.Context, userID int64, lessonID int64) (*dto.LessonDto, common.Errno) {
	lesson, err := s.lesson.GetLessonInfo(ctx, lessonID)
	if err != nil {
		logger.Error("GetLessonInfo error", zap.Error(err), zap.Int64("lessonID", lessonID))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	rel, _ := s.lesson.GetCourseLesson(ctx, lessonID)
	allowed := rel == nil || rel.ID == 0 || rel.EnableTrial == consts.IsEnable
	if !allowed && userID > 0 {
		purchased, _ := s.userRepo.HasPurchasedCourse(ctx, &do.HasPurchasedReq{UserID: userID, CourseID: rel.CourseGoodsID})
		allowed = purchased
	}
	videoURL := ""
	if allowed {
		videoURL = lesson.VideoKey
	}
	return &dto.LessonDto{
		ID:            lesson.ID,
		Name:          lesson.Name,
		Detail:        lesson.Detail,
		CategoryID:    lesson.CategoryID,
		VideoKey:      lesson.VideoKey,
		VideoURL:      videoURL,
		VideoFileName: lesson.VideoFileName,
		Duration:      lesson.Duration,
		Status:        lesson.Status,
		CreateBy:      lesson.CreateBy,
		UpdateBy:      lesson.UpdateBy,
		CreateAt:      lesson.CreateAt.UnixMilli(),
		UpdateAt:      lesson.UpdateAt.UnixMilli(),
	}, common.OK
}

func (s *Service) ListPurchasedCourse(ctx context.Context, userID int64, pager common.Pager) (*dto.PurchasedCourseListResp, common.Errno) {
	rights, err := s.userRepo.GetUserPurchasedCourses(ctx, userID)
	if err != nil {
		logger.Error("ListPurchasedCourse error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	total := int64(len(rights))
	list := make([]*dto.PurchasedCourseDto, 0)
	start := pager.GetOffset()
	for i, right := range rights {
		if int(i) < start || len(list) >= pager.Limit {
			continue
		}
		course, err := s.course.GetCourseByID(ctx, right.GoodsID)
		if err != nil {
			continue
		}
		cd := s.courseDto(ctx, course, true)
		list = append(list, &dto.PurchasedCourseDto{
			ID: course.ID, Name: course.Name, ServiceExpireTime: right.ServiceExpireTime,
			LearnExpireTime: right.LearnExpireTime, Features: cd.Features, UpdateStatus: course.UpdateStatus,
			HasPurchased: true, CoverKey: course.CoverKey, CoverURL: cd.CoverURL,
			DetailCoverKey: course.DetailCoverKey, DetailCoverURL: cd.DetailCOverURL, Detail: course.Detail,
		})
	}
	return &dto.PurchasedCourseListResp{Pager: pager, List: list, Total: total}, common.OK
}

func (s *Service) GetLessonLearnInfo(ctx context.Context, userID int64, req *dto.LessonLearnInfoReq) (*dto.LessonLearnInfoResp, common.Errno) {
	progress, err := s.lesson.GetLessonLearnProgress(ctx, &do.LessonLearnProgressReq{
		UserID:   userID,
		CourseID: req.CourseID,
		LessonID: req.LessonID,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("GetLessonLearnInfo error", zap.Error(err), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	if progress == nil {
		return &dto.LessonLearnInfoResp{CourseID: req.CourseID, LessonID: req.LessonID}, common.OK
	}
	return &dto.LessonLearnInfoResp{
		CourseID:       req.CourseID,
		LessonID:       req.LessonID,
		PlayPosition:   progress.PlayPosition,
		LearnStatus:    progress.LearnStatus,
		EntryTime:      progress.CreateAt.UnixMilli(),
		LastReportTime: progress.UpdateAt.UnixMilli(),
		InLearning:     progress.LearnStatus == consts.LearnStatusLearning,
	}, common.OK
}

func (s *Service) ReportLessonLearn(ctx context.Context, userID int64, req *dto.LessonLearnReportReq) common.Errno {
	now := time.Now()
	if err := s.lesson.UpsertLessonLearnProgress(ctx, &do.LessonLearnProgressUpdate{
		UserID:       userID,
		CourseID:     req.CourseID,
		LessonID:     req.LessonID,
		PlayPosition: req.PlayPosition,
		LearnStatus:  consts.LearnStatusLearning,
	}); err != nil {
		logger.Error("ReportLessonLearn UpsertLessonLearnProgress error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	if err := s.lesson.CreateLessonLearnRecord(ctx, &do.LessonLearnRecordCreate{
		UserID: userID, CourseID: req.CourseID, LessonID: req.LessonID,
		EntryTime: now, ExitTime: now, Duration: 0, LastType: req.Type,
	}); err != nil {
		logger.Error("ReportLessonLearn CreateLessonLearnRecord error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) AddCartGoods(ctx context.Context, userID int64, req *dto.AddGoodsReq) (int64, common.Errno) {
	purchased, err := s.userRepo.HasPurchasedCourse(ctx, &do.HasPurchasedReq{UserID: userID, CourseID: req.GoodsID})
	if err != nil {
		logger.Error("AddCartGoods HasPurchasedCourse error", zap.Error(err))
		return 0, *common.DataBaseErr.WithErr(err)
	}
	if purchased {
		return 0, *common.ParamErr.WithMsg("course already purchased")
	}
	if _, err := s.course.GetCourseByID(ctx, req.GoodsID); err != nil {
		logger.Error("AddCartGoods GetCourseByID error", zap.Error(err))
		return 0, *common.DataBaseErr.WithErr(err)
	}
	id, err := s.userRepo.AddToCart(ctx, &do.AddCartReq{UserID: userID, GoodsID: req.GoodsID, Quantity: 1})
	if err != nil {
		logger.Error("AddCartGoods AddToCart error", zap.Error(err))
		return 0, *common.DataBaseErr.WithErr(err)
	}
	return id, common.OK
}

func (s *Service) RemoveCartGoods(ctx context.Context, userID int64, req *dto.RemoveGoodsReq) common.Errno {
	if err := s.userRepo.RemoveFromCart(ctx, &do.RemoveCartReq{ID: req.ID, UserID: userID}); err != nil {
		logger.Error("RemoveCartGoods error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) ListCartGoods(ctx context.Context, userID int64, req *dto.ListCartGoodsReq) (*dto.ListGoodsResp, common.Errno) {
	carts, total, err := s.userRepo.ListCart(ctx, &do.ListCartReq{UserID: userID, Pager: req.Pager})
	if err != nil {
		logger.Error("ListCartGoods error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	list := make([]*dto.CartGoodsDto, 0, len(carts))
	for _, cart := range carts {
		course, err := s.course.GetCourseByID(ctx, cart.GoodsID)
		if err != nil {
			continue
		}
		purchased, _ := s.userRepo.HasPurchasedCourse(ctx, &do.HasPurchasedReq{UserID: userID, CourseID: course.ID})
		cd := s.courseDto(ctx, course, purchased)
		list = append(list, &dto.CartGoodsDto{CourseDto: *cd, CartID: cart.ID, GoodsID: cart.GoodsID, Quantity: cart.Quantity})
	}
	return &dto.ListGoodsResp{Pager: req.Pager, Total: total, List: list}, common.OK
}

func (s *Service) CalcOrderFee(ctx context.Context, userID int64, req *dto.OrderCalcFeeReq) (*dto.OrderCalcFeeResp, common.Errno) {
	if len(req.CourseIDs) == 0 {
		return nil, *common.ParamErr.WithMsg("empty course_ids")
	}
	resp := &dto.OrderCalcFeeResp{FeeUUID: tools.UUIDHex(), CourseFees: make([]*dto.CourseFeeDto, 0, len(req.CourseIDs))}
	for _, id := range req.CourseIDs {
		purchased, _ := s.userRepo.HasPurchasedCourse(ctx, &do.HasPurchasedReq{UserID: userID, CourseID: id})
		if purchased {
			return nil, *common.ParamErr.WithMsg("course already purchased")
		}
		course, err := s.course.GetCourseInfo(ctx, &do.CourseInfoReq{ID: id, Status: consts.IsEnable})
		if err != nil {
			logger.Error("CalcOrderFee GetCourseInfo error", zap.Error(err), zap.Int64("id", id))
			return nil, *common.DataBaseErr.WithErr(err)
		}
		snap := s.courseDto(ctx, course, false)
		resp.TotalFee += course.CoursePrice
		resp.TotalPayFee += course.CoursePrice
		resp.CourseFees = append(resp.CourseFees, &dto.CourseFeeDto{CourseID: id, Price: course.CoursePrice, PayFee: course.CoursePrice, GoodsSnap: snap})
	}
	resp.ExpireTime = time.Now().Add(s.feeTTL()).UnixMilli()
	data, _ := json.Marshal(resp)
	if err := s.orderFee.SetOrderFee(s.orderFeeKey(resp.FeeUUID), data, s.feeTTL()); err != nil {
		logger.Error("CalcOrderFee SetOrderFee error", zap.Error(err))
		return nil, *common.RedisErr.WithErr(err)
	}
	return resp, common.OK
}

func (s *Service) PayNow(ctx context.Context, userID int64, req *dto.OrderPayNowReq) (*dto.OrderPayNowResp, common.Errno) {
	data, err := s.orderFee.GetOrderFee(s.orderFeeKey(req.FeeUUID))
	if err != nil {
		return nil, *common.ParamErr.WithMsg("fee expired")
	}
	fee := &dto.OrderCalcFeeResp{}
	if err = json.Unmarshal(data, fee); err != nil {
		logger.Error("PayNow unmarshal fee error", zap.Error(err))
		return nil, *common.ServerErr.WithErr(err)
	}
	orderID := s.snow.NextID()
	now := time.Now().UnixMilli()
	innerTradeNo := tools.UUIDHex()
	items := make([]*do.CreateOrderItemReq, 0, len(fee.CourseFees))
	for _, item := range fee.CourseFees {
		snap, _ := json.Marshal(item.GoodsSnap)
		items = append(items, &do.CreateOrderItemReq{
			GoodsID: item.CourseID, GoodsType: consts.GoodsTypeCourse, Quantity: 1,
			PaymentAmount: item.PayFee, DiscountAmount: item.DiscountFee, GoodsSnap: string(snap),
		})
	}
	err = s.order.CreateOrderWithItems(ctx, &do.CreateOrderReq{
		ID: orderID, UserID: userID, Status: consts.OrderStatusPending,
		OrderSource: consts.OrderSourceCustomer, OrderAmount: fee.TotalFee,
		DiscountAmount: fee.TotalDiscountFee, PaymentAmount: fee.TotalPayFee,
		InnerTradeNo: innerTradeNo, OrderDesc: "课程订单", UserRemark: req.Remark,
		CreateAt: now, CreateBy: userID,
	}, items)
	if err != nil {
		logger.Error("PayNow CreateOrderWithItems error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return s.mockPayResp(orderID, innerTradeNo), common.OK
}

func (s *Service) PayLater(ctx context.Context, userID int64, req *dto.OrderPayLaterReq) (*dto.OrderPayNowResp, common.Errno) {
	order, err := s.order.GetOrderByID(ctx, &do.GetOrderReq{
		OrderID: req.OrderID, UserID: userID, Status: consts.OrderStatusPending,
	})
	if err != nil {
		logger.Error("PayLater GetOrderByID error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return s.mockPayResp(order.ID, order.InnerTradeNo), common.OK
}

func (s *Service) CancelOrder(ctx context.Context, userID int64, req *dto.CancelOrderReq) common.Errno {
	now := time.Now().UnixMilli()
	if err := s.order.UpdateOrderStatus(ctx, &do.UpdateOrderStatusReq{
		OrderID: req.OrderID, NewStatus: consts.OrderStatusCanceled,
		CancelAt: now, CancelType: consts.CancelTypeUser,
		CancelBy: userID, CancelReason: req.Reason,
	}); err != nil {
		logger.Error("CancelOrder error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) ListOrders(ctx context.Context, userID int64, req *dto.OrderListReq) (*dto.UserOrderListResp, common.Errno) {
	orders, total, err := s.order.ListOrders(ctx, &do.ListOrderReq{UserID: userID, Status: req.Status, Pager: req.Pager})
	if err != nil {
		logger.Error("ListOrders error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	list := make([]*dto.OrderInfoResp, 0, len(orders))
	for _, order := range orders {
		list = append(list, s.orderInfoDto(ctx, order))
	}
	return &dto.UserOrderListResp{Pager: req.Pager, Total: total, List: list}, common.OK
}

func (s *Service) GetOrderInfo(ctx context.Context, userID int64, orderID int64) (*dto.OrderInfoResp, common.Errno) {
	order, err := s.order.GetOrderByID(ctx, &do.GetOrderReq{OrderID: orderID, UserID: userID})
	if err != nil {
		logger.Error("GetOrderInfo error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return s.orderInfoDto(ctx, order), common.OK
}

func (s *Service) DeliverOrder(ctx context.Context, orderID int64) common.Errno {
	if err := s.order.DeliverOrder(ctx, orderID); err != nil {
		logger.Error("DeliverOrder error", zap.Error(err), zap.Int64("orderID", orderID))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) RefundOrder(ctx context.Context, adminID int64, req *dto.RefundOrderReq) common.Errno {
	now := time.Now().UnixMilli()
	order, err := s.order.GetOrderByID(ctx, &do.GetOrderReq{OrderID: req.OrderID})
	if err != nil {
		logger.Error("RefundOrder GetOrderByID error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	if err := s.order.RefundOrder(ctx, req.OrderID, &do.CreateRefundReq{
		UserID: order.UserID, OrderID: req.OrderID, ItemIds: joinInt64(req.ItemIDs),
		ApplyAt: now, Reason: req.Reason, Status: consts.RefundStatusDone, Amount: req.Amount,
		InnerTradeNo: tools.UUIDHex(), RefundID: "mock_" + tools.UUIDHex(), ApplyUserID: adminID, DoneAt: now,
	}, now); err != nil {
		logger.Error("RefundOrder error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) CancelTimeoutOrders(ctx context.Context) common.Errno {
	timeout := time.Now().Add(-s.payTTL()).UnixMilli()
	now := time.Now().UnixMilli()
	if err := s.order.CancelTimeoutOrders(ctx, timeout, now); err != nil {
		logger.Error("CancelTimeoutOrders error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) ListContinueLearn(ctx context.Context, userID int64, pager common.Pager) (*dto.ContinueLearnResp, common.Errno) {
	progresses, total, err := s.lesson.GetUserLearnProgresses(ctx, userID, pager.GetOffset(), pager.Limit)
	if err != nil {
		logger.Error("ListContinueLearn error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	list := make([]*dto.ContinueLearnCourseDto, 0, len(progresses))
	for _, progress := range progresses {
		course, err := s.course.GetCourseByID(ctx, progress.CourseID)
		if err != nil {
			continue
		}
		lesson, err := s.lesson.GetLessonInfo(ctx, progress.LessonID)
		if err != nil {
			continue
		}
		coverURL := ""
		if course.CoverKey != "" {
			urlMap, _ := s.storage.GetPreviewUrl(ctx, &do.GetPreviewUrl{Keys: []string{course.CoverKey}, Expire: 6})
			coverURL = urlMap[course.CoverKey]
		}
		list = append(list, &dto.ContinueLearnCourseDto{
			CourseID:       course.ID,
			CourseName:     course.Name,
			CourseCoverKey: course.CoverKey,
			CourseCoverURL: coverURL,
			LessonID:       lesson.ID,
			LessonName:     lesson.Name,
			PlayPosition:   progress.PlayPosition,
			LearnStatus:    progress.LearnStatus,
		})
	}
	return &dto.ContinueLearnResp{Pager: pager, List: list, Total: total}, common.OK
}

func (s *Service) AdminListCustomerUsers(ctx context.Context, req *dto.AdminCustomerUserListReq) (*dto.AdminCustomerUserListResp, common.Errno) {
	doReq := &do.CustomerUserListReq{UserID: req.UserID, Status: req.Status, Pager: req.Pager}
	if req.Mobile != "" {
		mobileHash := secure.MobileSHA256(req.Mobile, s.conf.Security.MobileSHA256Salt)
		mu, err := s.userRepo.GetMobileUserByHash(ctx, mobileHash)
		if err != nil {
			return &dto.AdminCustomerUserListResp{Pager: req.Pager, List: []*dto.CustomerUserInfoDto{}}, common.OK
		}
		doReq.UserID = mu.UserID
	}
	users, total, err := s.userRepo.ListUsers(ctx, doReq)
	if err != nil {
		logger.Error("AdminListCustomerUsers error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	list := make([]*dto.CustomerUserInfoDto, 0, len(users))
	for _, user := range users {
		if info, err := s.buildCustomerUserInfo(ctx, user.ID); err == nil {
			list = append(list, info)
		}
	}
	return &dto.AdminCustomerUserListResp{Pager: req.Pager, Total: total, List: list}, common.OK
}

func (s *Service) AdminUpdateCustomerStatus(ctx context.Context, req *dto.AdminCustomerUserStatusReq) common.Errno {
	if err := s.userRepo.UpdateUserStatus(ctx, &do.CustomerUserStatusReq{UserID: req.UserID, Status: req.Status}); err != nil {
		logger.Error("AdminUpdateCustomerStatus error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	if req.Status != consts.IsEnable {
		_ = s.verify.CleanCustomerToken(ctx, req.UserID)
	}
	return common.OK
}

func (s *Service) AdminListOrders(ctx context.Context, req *dto.OrderListReq) (*dto.UserOrderListResp, common.Errno) {
	orders, total, err := s.order.ListOrders(ctx, &do.ListOrderReq{Status: req.Status, Pager: req.Pager})
	if err != nil {
		logger.Error("AdminListOrders error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	list := make([]*dto.OrderInfoResp, 0, len(orders))
	for _, order := range orders {
		list = append(list, s.orderInfoDto(ctx, order))
	}
	return &dto.UserOrderListResp{Pager: req.Pager, Total: total, List: list}, common.OK
}

func (s *Service) AdminOrderInfo(ctx context.Context, orderID int64) (*dto.OrderInfoResp, common.Errno) {
	order, err := s.order.GetOrderByID(ctx, &do.GetOrderReq{OrderID: orderID})
	if err != nil {
		logger.Error("AdminOrderInfo error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return s.orderInfoDto(ctx, order), common.OK
}

func (s *Service) AdminOrderStats(ctx context.Context, req *dto.AdminOrderStatsReq) (*dto.AdminOrderStatsResp, common.Errno) {
	resp := &dto.AdminOrderStatsResp{ByStatus: []dto.StatusStat{}, ByGoods: []dto.GoodsStat{}}
	orders, err := s.order.GetOrdersByTimeRange(ctx, req.CreateStart, req.CreateEnd)
	if err != nil {
		logger.Error("AdminOrderStats GetOrdersByTimeRange error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	statusMap := map[int32]*dto.StatusStat{}
	for _, order := range orders {
		stat := statusMap[order.Status]
		if stat == nil {
			stat = &dto.StatusStat{Status: order.Status}
			statusMap[order.Status] = stat
		}
		stat.Count++
		stat.Amount += order.PaymentAmount
		resp.TotalPay += order.PaymentAmount
	}
	for _, stat := range statusMap {
		resp.ByStatus = append(resp.ByStatus, *stat)
	}
	items, _ := s.order.GetAllOrderItems(ctx)
	goodsMap := map[int64]*dto.GoodsStat{}
	for _, item := range items {
		stat := goodsMap[item.GoodsID]
		if stat == nil {
			stat = &dto.GoodsStat{GoodsID: item.GoodsID}
			if course, err := s.course.GetCourseByID(ctx, item.GoodsID); err == nil {
				stat.Name = course.Name
			}
			goodsMap[item.GoodsID] = stat
		}
		stat.Count += int64(item.Quantity)
		stat.Amount += item.PaymentAmount
	}
	for _, stat := range goodsMap {
		resp.ByGoods = append(resp.ByGoods, *stat)
	}
	return resp, common.OK
}

func (s *Service) WechatNotifySuccess(ctx context.Context, orderID int64) common.Errno {
	if orderID <= 0 {
		return common.OK
	}
	return s.DeliverOrder(ctx, orderID)
}

func (s *Service) WechatNotifyResponse() map[string]string {
	return map[string]string{"code": "SUCCESS", "message": "成功"}
}

// Helper functions

func (s *Service) courseDto(ctx context.Context, course *model.CourseGood, hasPurchased bool) *dto.CourseDto {
	features := make([]string, 0)
	_ = json.Unmarshal([]byte(course.Features), &features)
	urlMap := map[string]string{}
	if course.CoverKey != "" || course.DetailCoverKey != "" {
		urlMap, _ = s.storage.GetPreviewUrl(ctx, &do.GetPreviewUrl{Keys: []string{course.CoverKey, course.DetailCoverKey}, Expire: 6})
	}
	return &dto.CourseDto{
		ID: course.ID, Name: course.Name, CoursePrice: course.CoursePrice, ServiceTime: course.ServiceTime,
		LearnTime: course.LearnTime, Status: course.Status, Sort: course.Sort, Features: features,
		UpdateStatus: course.UpdateStatus, CoverKey: course.CoverKey, CoverURL: urlMap[course.CoverKey],
		DetailCoverKey: course.DetailCoverKey, DetailCOverURL: urlMap[course.DetailCoverKey], Detail: course.Detail,
		CreateBy: course.CreateBy, UpdateBy: course.UpdateBy, CreateAt: course.CreateAt.UnixMilli(), UpdateAt: course.UpdateAt.UnixMilli(),
		HasPurchased: hasPurchased,
	}
}

func (s *Service) courseCatalogs(ctx context.Context, courseID int64, purchased bool) ([]*dto.CatalogDto, int64, int32, error) {
	catalogs, err := s.course.GetCatalogs(ctx, courseID)
	if err != nil {
		return nil, 0, 0, err
	}
	rels, err := s.course.GetCourseLessons(ctx, courseID)
	if err != nil {
		return nil, 0, 0, err
	}
	byCatalog := map[int64][]*model.CourseLesson{}
	for _, rel := range rels {
		byCatalog[rel.CatalogID] = append(byCatalog[rel.CatalogID], rel)
	}
	ret := make([]*dto.CatalogDto, 0, len(catalogs))
	var duration int64
	var count int32
	for _, catalog := range catalogs {
		lessons := make([]*dto.CatalogLessonDto, 0)
		for _, rel := range byCatalog[catalog.ID] {
			lesson, err := s.lesson.GetLessonInfo(ctx, rel.LessonID)
			if err != nil {
				continue
			}
			count++
			duration += int64(lesson.Duration)
			videoURL := ""
			if purchased || rel.EnableTrial == consts.IsEnable {
				videoURL = lesson.VideoKey
			}
			lessons = append(lessons, &dto.CatalogLessonDto{
				ID: rel.ID, LessonID: rel.LessonID, Name: rel.Name, LessonName: lesson.Name,
				Detail: lesson.Detail, VideoURL: videoURL, VideoFileName: lesson.VideoFileName,
				Duration: lesson.Duration, Status: lesson.Status, ShowTime: rel.ShowTime.UnixMilli(),
				EnableTrial: rel.EnableTrial == consts.IsEnable,
			})
		}
		ret = append(ret, &dto.CatalogDto{
			ID: catalog.ID, ParentID: catalog.ParentID, Level: catalog.Level, Name: catalog.Name,
			CourseID: catalog.CourseID, Sort: catalog.Sort, Lessons: lessons, LessonCount: int32(len(lessons)),
		})
	}
	return ret, duration, count, nil
}

func (s *Service) purchasedMap(ctx context.Context, userID int64) map[int64]bool {
	ret := map[int64]bool{}
	if userID <= 0 {
		return ret
	}
	rights, _ := s.userRepo.GetUserPurchasedCourses(ctx, userID)
	for _, right := range rights {
		ret[right.GoodsID] = true
	}
	return ret
}

func (s *Service) orderInfoDto(ctx context.Context, order *model.Order) *dto.OrderInfoResp {
	od := dto.OrderDto{
		ID: order.ID, UserID: order.UserID, Status: order.Status, OrderSource: order.OrderSource,
		OrderAmount: order.OrderAmount, DiscountAmount: order.DiscountAmount, PaymentAmount: order.PaymentAmount,
		TradeNo: order.TradeNo, InnerTradeNo: order.InnerTradeNo, OrderDesc: order.OrderDesc, PaymentAt: order.PaymentAt,
		UserRemark: order.UserRemark, RefundAmount: order.RefundAmount,
		CreateAt: order.CreateAt, CreateBy: order.CreateBy,
	}
	if order.ReceiverConfirmAt != nil {
		od.ReceiverConfirmAt = order.ReceiverConfirmAt
	}
	if order.ReceiverConfirmType != nil {
		od.ReceiverConfirmType = order.ReceiverConfirmType
	}
	if order.RefundAt != nil {
		od.RefundAt = order.RefundAt
	}
	if order.CancelAt != nil {
		od.CancelAt = order.CancelAt
	}
	if order.CancelType != nil {
		od.CancelType = order.CancelType
	}
	if order.CancelBy != nil {
		od.CancelBy = order.CancelBy
	}
	if order.CancelReason != nil {
		od.CancelReason = order.CancelReason
	}
	resp := &dto.OrderInfoResp{OrderDto: od, Items: make([]*dto.OrderItemDto, 0), Refunds: make([]*dto.RefundDto, 0)}
	items, _ := s.order.GetOrderItemsByOrderID(ctx, order.ID)
	for _, item := range items {
		var snap any
		_ = json.Unmarshal([]byte(item.GoodsSnap), &snap)
		resp.Items = append(resp.Items, &dto.OrderItemDto{
			ID: item.ID, OrderID: item.OrderID, UserID: item.UserID, GoodsID: item.GoodsID,
			GoodsType: item.GoodsType, Quantity: item.Quantity, PaymentAmount: item.PaymentAmount,
			DiscountAmount: item.DiscountAmount, GoodsSnap: snap,
		})
	}
	refunds, _ := s.order.ListRefundsByOrderID(ctx, order.ID)
	for _, refund := range refunds {
		resp.Refunds = append(resp.Refunds, &dto.RefundDto{
			ID: refund.ID, Amount: refund.Amount, ItemIDs: splitInt64(refund.ItemIds),
			ApplyAt: refund.ApplyAt, Status: refund.Status, DoneAt: refund.DoneAt,
			Reason: refund.Reason, RefundID: refund.RefundID, ApplyUserID: refund.ApplyUserID,
		})
	}
	if info, err := s.buildCustomerUserInfo(ctx, order.UserID); err == nil && info.User != nil {
		resp.UserName = info.User.NickName
		if info.MobileUser != nil {
			resp.UserMobile = info.MobileUser.Mobile
		}
	}
	return resp
}

func (s *Service) mockPayResp(orderID int64, innerTradeNo string) *dto.OrderPayNowResp {
	return &dto.OrderPayNowResp{
		OrderID: orderID, AppID: s.conf.Wechat.AppID, TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr: tools.UUIDHex(), Package: "prepay_id=" + innerTradeNo, SignType: "RSA", PaySign: "mock",
		CodeURL: "weixin://wxpay/bizpayurl?pr=" + innerTradeNo, TradeType: "NATIVE",
	}
}

func (s *Service) orderFeeKey(uuid string) string { return "mall:order:fee:" + uuid }

func (s *Service) feeTTL() time.Duration {
	if s.conf.Order.FeeLockTimeoutMinutes > 0 {
		return time.Duration(s.conf.Order.FeeLockTimeoutMinutes) * time.Minute
	}
	return consts.ExpireOrderFeeTime
}

func (s *Service) payTTL() time.Duration {
	if s.conf.Order.PayTimeoutMinutes > 0 {
		return time.Duration(s.conf.Order.PayTimeoutMinutes) * time.Minute
	}
	return consts.ExpireOrderPayTime
}

func joinInt64(ids []int64) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.FormatInt(id, 10))
	}
	return strings.Join(parts, ",")
}

func splitInt64(raw string) []int64 {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	ret := make([]int64, 0, len(parts))
	for _, part := range parts {
		id, _ := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if id > 0 {
			ret = append(ret, id)
		}
	}
	return ret
}
