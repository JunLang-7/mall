package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/secure"
	"github.com/JunLang-7/mall/utils/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) ListCustomerCourse(ctx context.Context, userID int64, req *dto.CourseListReq) (*dto.CourseListResp, common.Errno) {
	q := query.Use(s.db).CourseGood
	tx := q.WithContext(ctx).Where(q.Status.Eq(consts.IsEnable))
	if req.ID > 0 {
		tx = tx.Where(q.ID.Eq(req.ID))
	}
	if req.NameKW != "" {
		tx = tx.Where(q.Name.Like(tools.GetAllLike(req.NameKW)))
	}
	list, total, err := tx.Order(q.Sort.Asc(), q.ID.Desc()).FindByPage(req.GetOffset(), req.Limit)
	if err != nil {
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
	q := query.Use(s.db).CourseGood
	course, err := q.WithContext(ctx).Where(q.ID.Eq(id), q.Status.Eq(consts.IsEnable)).First()
	if err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	catalogs, totalDuration, lessonCount, err := s.courseCatalogs(ctx, id, s.hasPurchased(ctx, userID, id))
	if err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	dtoCourse := s.courseDto(ctx, course, s.hasPurchased(ctx, userID, id))
	return &dto.CustomerCourseDetailDto{
		CourseDto:      *dtoCourse,
		TotalDuration:  totalDuration,
		LessonCount:    lessonCount,
		Catalogs:       catalogs,
		HasPurchased:   dtoCourse.HasPurchased,
	}, common.OK
}

func (s *Service) GetLessonInfo(ctx context.Context, userID int64, lessonID int64) (*dto.LessonDto, common.Errno) {
	var lesson model.Lesson
	if err := s.db.WithContext(ctx).First(&lesson, lessonID).Error; err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	var rel model.CourseLesson
	_ = s.db.WithContext(ctx).Where("lesson_id = ?", lessonID).First(&rel).Error
	allowed := rel.ID == 0 || rel.EnableTrial == consts.IsEnable || s.hasPurchased(ctx, userID, rel.CourseGoodsID)
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
	var rights []model.UserCourseGood
	tx := s.db.WithContext(ctx).Where("user_id = ? AND learn_expire_time > ?", userID, time.Now().UnixMilli())
	var total int64
	_ = tx.Model(&model.UserCourseGood{}).Count(&total).Error
	if err := tx.Offset(pager.GetOffset()).Limit(pager.Limit).Find(&rights).Error; err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	list := make([]*dto.PurchasedCourseDto, 0, len(rights))
	for _, right := range rights {
		var course model.CourseGood
		if err := s.db.WithContext(ctx).First(&course, right.GoodsID).Error; err != nil {
			continue
		}
		cd := s.courseDto(ctx, &course, true)
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
	var progress model.LessonLearnProgress
	err := s.db.WithContext(ctx).Where("user_id = ? AND course_id = ? AND lesson_id = ?", userID, req.CourseID, req.LessonID).First(&progress).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return &dto.LessonLearnInfoResp{
		CourseID: req.CourseID, LessonID: req.LessonID, PlayPosition: progress.PlayPosition,
		LearnStatus: progress.LearnStatus, LastReportTime: progress.UpdateAt.UnixMilli(),
		InLearning: progress.LearnStatus == consts.LearnStatusLearning,
	}, common.OK
}

func (s *Service) ReportLessonLearn(ctx context.Context, userID int64, req *dto.LessonLearnReportReq) common.Errno {
	now := time.Now()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		progress := model.LessonLearnProgress{
			CourseID: req.CourseID, LessonID: req.LessonID, UserID: userID,
			PlayPosition: req.PlayPosition, LearnStatus: consts.LearnStatusLearning,
			CreateAt: now, UpdateAt: now,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "course_id"}, {Name: "lesson_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"play_position": req.PlayPosition,
				"learn_status":  consts.LearnStatusLearning,
				"update_at":     now,
			}),
		}).Create(&progress).Error; err != nil {
			return err
		}
		return tx.Create(&model.LessonLearnRecord{
			CourseID: req.CourseID, LessonID: req.LessonID, UserID: userID,
			EntryTime: now, ExitTime: now, Duration: 0, LastType: req.Type,
		}).Error
	})
	if err != nil {
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) AddCartGoods(ctx context.Context, userID int64, req *dto.AddGoodsReq) (int64, common.Errno) {
	if s.hasPurchased(ctx, userID, req.GoodsID) {
		return 0, *common.ParamErr.WithMsg("course already purchased")
	}
	var course model.CourseGood
	if err := s.db.WithContext(ctx).Where("id = ? AND status = ?", req.GoodsID, consts.IsEnable).First(&course).Error; err != nil {
		return 0, *common.DataBaseErr.WithErr(err)
	}
	cart := model.UserCart{UserID: userID, GoodsID: req.GoodsID, Quantity: 1, AddAt: time.Now()}
	err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "goods_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"quantity": 1, "add_at": time.Now()}),
	}).Create(&cart).Error
	if err != nil {
		return 0, *common.DataBaseErr.WithErr(err)
	}
	return cart.ID, common.OK
}

func (s *Service) RemoveCartGoods(ctx context.Context, userID int64, req *dto.RemoveGoodsReq) common.Errno {
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", req.ID, userID).Delete(&model.UserCart{}).Error; err != nil {
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) ListCartGoods(ctx context.Context, userID int64, req *dto.ListCartGoodsReq) (*dto.ListGoodsResp, common.Errno) {
	var carts []model.UserCart
	tx := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("add_at desc")
	var total int64
	_ = tx.Model(&model.UserCart{}).Count(&total).Error
	if err := tx.Offset(req.GetOffset()).Limit(req.Limit).Find(&carts).Error; err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	list := make([]*dto.CartGoodsDto, 0, len(carts))
	for _, cart := range carts {
		var course model.CourseGood
		if err := s.db.WithContext(ctx).First(&course, cart.GoodsID).Error; err != nil {
			continue
		}
		cd := s.courseDto(ctx, &course, s.hasPurchased(ctx, userID, course.ID))
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
		if s.hasPurchased(ctx, userID, id) {
			return nil, *common.ParamErr.WithMsg("course already purchased")
		}
		var course model.CourseGood
		if err := s.db.WithContext(ctx).Where("id = ? AND status = ?", id, consts.IsEnable).First(&course).Error; err != nil {
			return nil, *common.DataBaseErr.WithErr(err)
		}
		snap := s.courseDto(ctx, &course, false)
		resp.TotalFee += course.CoursePrice
		resp.TotalPayFee += course.CoursePrice
		resp.CourseFees = append(resp.CourseFees, &dto.CourseFeeDto{CourseID: id, Price: course.CoursePrice, PayFee: course.CoursePrice, GoodsSnap: snap})
	}
	resp.ExpireTime = time.Now().Add(s.feeTTL()).UnixMilli()
	data, _ := json.Marshal(resp)
	if err := s.rds.Set(s.orderFeeKey(resp.FeeUUID), data, s.feeTTL()).Err(); err != nil {
		return nil, *common.RedisErr.WithErr(err)
	}
	return resp, common.OK
}

func (s *Service) PayNow(ctx context.Context, userID int64, req *dto.OrderPayNowReq) (*dto.OrderPayNowResp, common.Errno) {
	data, err := s.rds.Get(s.orderFeeKey(req.FeeUUID)).Bytes()
	if err != nil {
		return nil, *common.ParamErr.WithMsg("fee expired")
	}
	fee := &dto.OrderCalcFeeResp{}
	if err = json.Unmarshal(data, fee); err != nil {
		return nil, *common.ServerErr.WithErr(err)
	}
	orderID := s.snow.NextID()
	now := time.Now().UnixMilli()
	innerTradeNo := tools.UUIDHex()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order := &model.Order{
			ID: orderID, UserID: userID, Status: consts.OrderStatusPending, OrderSource: consts.OrderSourceCustomer,
			OrderAmount: fee.TotalFee, DiscountAmount: fee.TotalDiscountFee, PaymentAmount: fee.TotalPayFee,
			InnerTradeNo: innerTradeNo, OrderDesc: "课程订单", UserRemark: req.Remark, CreateAt: now, CreateBy: userID,
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for _, item := range fee.CourseFees {
			snap, _ := json.Marshal(item.GoodsSnap)
			if err := tx.Create(&model.OrderItem{OrderID: orderID, UserID: userID, GoodsID: item.CourseID, GoodsType: consts.GoodsTypeCourse, Quantity: 1, PaymentAmount: item.PayFee, DiscountAmount: item.DiscountFee, GoodsSnap: string(snap)}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return s.mockPayResp(orderID, innerTradeNo), common.OK
}

func (s *Service) PayLater(ctx context.Context, userID int64, req *dto.OrderPayLaterReq) (*dto.OrderPayNowResp, common.Errno) {
	var order model.Order
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ? AND status = ?", req.OrderID, userID, consts.OrderStatusPending).First(&order).Error; err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return s.mockPayResp(order.ID, order.InnerTradeNo), common.OK
}

func (s *Service) CancelOrder(ctx context.Context, userID int64, req *dto.CancelOrderReq) common.Errno {
	now := time.Now().UnixMilli()
	err := s.db.WithContext(ctx).Model(&model.Order{}).Where("id = ? AND user_id = ? AND status = ?", req.OrderID, userID, consts.OrderStatusPending).Updates(map[string]interface{}{
		"status": consts.OrderStatusCanceled, "cancel_at": now, "cancel_type": consts.CancelTypeUser, "cancel_by": userID, "cancel_reason": req.Reason,
	}).Error
	if err != nil {
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) ListOrders(ctx context.Context, userID int64, req *dto.OrderListReq) (*dto.UserOrderListResp, common.Errno) {
	var orders []model.Order
	tx := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("create_at desc")
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	var total int64
	_ = tx.Model(&model.Order{}).Count(&total).Error
	if err := tx.Offset(req.GetOffset()).Limit(req.Limit).Find(&orders).Error; err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	list := make([]*dto.OrderInfoResp, 0, len(orders))
	for _, order := range orders {
		list = append(list, s.orderInfoDto(ctx, &order))
	}
	return &dto.UserOrderListResp{Pager: req.Pager, Total: total, List: list}, common.OK
}

func (s *Service) GetOrderInfo(ctx context.Context, userID int64, orderID int64) (*dto.OrderInfoResp, common.Errno) {
	var order model.Order
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return s.orderInfoDto(ctx, &order), common.OK
}

func (s *Service) PaymentQuery(ctx context.Context, userID int64, orderID int64) (*dto.OrderInfoResp, common.Errno) {
	return s.GetOrderInfo(ctx, userID, orderID)
}

func (s *Service) DeliverOrder(ctx context.Context, orderID int64) common.Errno {
	var order model.Order
	if err := s.db.WithContext(ctx).Where("id = ?", orderID).First(&order).Error; err != nil {
		return *common.DataBaseErr.WithErr(err)
	}
	if order.Status == consts.OrderStatusDone {
		return common.OK
	}
	now := time.Now().UnixMilli()
	returnErr := common.OK
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var items []model.OrderItem
		if err := tx.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
			return err
		}
		for _, item := range items {
			if item.GoodsType != consts.GoodsTypeCourse {
				continue
			}
			if err := tx.Where("user_id = ? AND order_item_id = ?", order.UserID, item.ID).FirstOrCreate(&model.UserCourseGood{
				UserID: order.UserID, OrderID: order.ID, OrderItemID: item.ID, GoodsID: item.GoodsID, GoodsType: item.GoodsType,
				BuyTime: now, LearnExpireTime: addMonthByCode(now, 12), ServiceExpireTime: addMonthByCode(now, 12),
			}).Error; err != nil {
				return err
			}
			_ = tx.Where("user_id = ? AND goods_id = ?", order.UserID, item.GoodsID).Delete(&model.UserCart{}).Error
		}
		return tx.Model(&model.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{"status": consts.OrderStatusDone, "payment_at": now}).Error
	})
	if err != nil {
		returnErr = *common.DataBaseErr.WithErr(err)
	}
	return returnErr
}

func (s *Service) RefundOrder(ctx context.Context, adminID int64, req *dto.RefundOrderReq) common.Errno {
	var order model.Order
	if err := s.db.WithContext(ctx).First(&order, req.OrderID).Error; err != nil {
		return *common.DataBaseErr.WithErr(err)
	}
	now := time.Now().UnixMilli()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		refund := &model.OrderRefund{UserID: order.UserID, OrderID: order.ID, ItemIds: joinInt64(req.ItemIDs), ApplyAt: now, Reason: req.Reason, Status: consts.RefundStatusDone, Amount: req.Amount, InnerTradeNo: tools.UUIDHex(), RefundID: "mock_" + tools.UUIDHex(), ApplyUserID: adminID, DoneAt: now}
		if err := tx.Create(refund).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Order{}).Where("id = ?", order.ID).Updates(map[string]interface{}{"status": consts.OrderStatusRefunded, "refund_amount": req.Amount, "refund_at": now}).Error; err != nil {
			return err
		}
		return tx.Model(&model.UserCourseGood{}).Where("order_id = ?", order.ID).Updates(map[string]interface{}{"learn_expire_time": now, "service_expire_time": now}).Error
	})
	if err != nil {
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) CancelTimeoutOrders(ctx context.Context) common.Errno {
	timeout := time.Now().Add(-s.payTTL()).UnixMilli()
	now := time.Now().UnixMilli()
	err := s.db.WithContext(ctx).Model(&model.Order{}).Where("status = ? AND create_at < ?", consts.OrderStatusPending, timeout).Updates(map[string]interface{}{
		"status": consts.OrderStatusCanceled, "cancel_at": now, "cancel_type": consts.CancelTypeTimeout, "cancel_by": consts.SystemUserID, "cancel_reason": "timeout",
	}).Error
	if err != nil {
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

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
	var catalogs []model.CourseCatalog
	if err := s.db.WithContext(ctx).Where("course_id = ?", courseID).Order("sort asc").Find(&catalogs).Error; err != nil {
		return nil, 0, 0, err
	}
	var rels []model.CourseLesson
	if err := s.db.WithContext(ctx).Where("course_goods_id = ?", courseID).Order("sort asc").Find(&rels).Error; err != nil {
		return nil, 0, 0, err
	}
	byCatalog := map[int64][]model.CourseLesson{}
	for _, rel := range rels {
		byCatalog[rel.CatalogID] = append(byCatalog[rel.CatalogID], rel)
	}
	ret := make([]*dto.CatalogDto, 0, len(catalogs))
	var duration int64
	var count int32
	for _, catalog := range catalogs {
		lessons := make([]*dto.CatalogLessonDto, 0)
		for _, rel := range byCatalog[catalog.ID] {
			var lesson model.Lesson
			if err := s.db.WithContext(ctx).First(&lesson, rel.LessonID).Error; err != nil {
				continue
			}
			count++
			duration += int64(lesson.Duration)
			videoURL := ""
			if purchased || rel.EnableTrial == consts.IsEnable {
				videoURL = lesson.VideoKey
			}
			lessons = append(lessons, &dto.CatalogLessonDto{ID: rel.ID, LessonID: rel.LessonID, Name: rel.Name, LessonName: lesson.Name, Detail: lesson.Detail, VideoURL: videoURL, VideoFileName: lesson.VideoFileName, Duration: lesson.Duration, Status: lesson.Status, ShowTime: rel.ShowTime.UnixMilli(), EnableTrial: rel.EnableTrial == consts.IsEnable})
		}
		ret = append(ret, &dto.CatalogDto{ID: catalog.ID, ParentID: catalog.ParentID, Level: catalog.Level, Name: catalog.Name, CourseID: catalog.CourseID, Sort: catalog.Sort, Lessons: lessons, LessonCount: int32(len(lessons))})
	}
	return ret, duration, count, nil
}

func (s *Service) purchasedMap(ctx context.Context, userID int64) map[int64]bool {
	ret := map[int64]bool{}
	if userID <= 0 {
		return ret
	}
	var rights []model.UserCourseGood
	_ = s.db.WithContext(ctx).Where("user_id = ? AND learn_expire_time > ?", userID, time.Now().UnixMilli()).Find(&rights).Error
	for _, right := range rights {
		ret[right.GoodsID] = true
	}
	return ret
}

func (s *Service) hasPurchased(ctx context.Context, userID int64, courseID int64) bool {
	if userID <= 0 {
		return false
	}
	var count int64
	_ = s.db.WithContext(ctx).Model(&model.UserCourseGood{}).Where("user_id = ? AND goods_id = ? AND learn_expire_time > ?", userID, courseID, time.Now().UnixMilli()).Count(&count).Error
	return count > 0
}

func (s *Service) orderInfoDto(ctx context.Context, order *model.Order) *dto.OrderInfoResp {
	resp := &dto.OrderInfoResp{OrderDto: dto.OrderDto{
		ID: order.ID, UserID: order.UserID, Status: order.Status, OrderSource: order.OrderSource,
		OrderAmount: order.OrderAmount, DiscountAmount: order.DiscountAmount, PaymentAmount: order.PaymentAmount,
		TradeNo: order.TradeNo, InnerTradeNo: order.InnerTradeNo, OrderDesc: order.OrderDesc, PaymentAt: order.PaymentAt,
		UserRemark: order.UserRemark, ReceiverConfirmAt: order.ReceiverConfirmAt, ReceiverConfirmType: order.ReceiverConfirmType,
		RefundAmount: order.RefundAmount, RefundAt: order.RefundAt, CancelAt: order.CancelAt, CancelType: order.CancelType,
		CancelBy: order.CancelBy, CancelReason: order.CancelReason, CreateAt: order.CreateAt, CreateBy: order.CreateBy,
	}, Items: make([]*dto.OrderItemDto, 0), Refunds: make([]*dto.RefundDto, 0)}
	var items []model.OrderItem
	_ = s.db.WithContext(ctx).Where("order_id = ?", order.ID).Find(&items).Error
	for _, item := range items {
		var snap any
		_ = json.Unmarshal([]byte(item.GoodsSnap), &snap)
		resp.Items = append(resp.Items, &dto.OrderItemDto{ID: item.ID, OrderID: item.OrderID, UserID: item.UserID, GoodsID: item.GoodsID, GoodsType: item.GoodsType, Quantity: item.Quantity, PaymentAmount: item.PaymentAmount, DiscountAmount: item.DiscountAmount, GoodsSnap: snap})
	}
	var refunds []model.OrderRefund
	_ = s.db.WithContext(ctx).Where("order_id = ?", order.ID).Find(&refunds).Error
	for _, refund := range refunds {
		resp.Refunds = append(resp.Refunds, &dto.RefundDto{ID: refund.ID, Amount: refund.Amount, ItemIDs: splitInt64(refund.ItemIds), ApplyAt: refund.ApplyAt, Status: refund.Status, DoneAt: refund.DoneAt, Reason: refund.Reason, RefundID: refund.RefundID, ApplyUserID: refund.ApplyUserID})
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

func addMonthByCode(ms int64, months int) int64 {
	return time.UnixMilli(ms).AddDate(0, months, 0).UnixMilli()
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

func (s *Service) AdminListCustomerUsers(ctx context.Context, req *dto.AdminCustomerUserListReq) (*dto.AdminCustomerUserListResp, common.Errno) {
	tx := s.db.WithContext(ctx).Model(&model.User{})
	if req.UserID > 0 {
		tx = tx.Where("id = ?", req.UserID)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	if req.Mobile != "" {
		mobileHash := secure.MobileSHA256(req.Mobile, s.conf.Security.MobileSHA256Salt)
		var mu model.MobileUser
		if err := s.db.WithContext(ctx).Where("mobile_sha256 = ?", mobileHash).First(&mu).Error; err == nil {
			tx = tx.Where("id = ?", mu.UserID)
		} else {
			return &dto.AdminCustomerUserListResp{Pager: req.Pager, List: []*dto.CustomerUserInfoDto{}}, common.OK
		}
	}
	var total int64
	_ = tx.Count(&total).Error
	var users []model.User
	if err := tx.Offset(req.GetOffset()).Limit(req.Limit).Find(&users).Error; err != nil {
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
	if err := s.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", req.UserID).Update("status", req.Status).Error; err != nil {
		return *common.DataBaseErr.WithErr(err)
	}
	if req.Status != consts.IsEnable {
		_ = s.verify.CleanCustomerToken(ctx, req.UserID)
	}
	return common.OK
}

func (s *Service) AdminListOrders(ctx context.Context, req *dto.OrderListReq) (*dto.UserOrderListResp, common.Errno) {
	var orders []model.Order
	tx := s.db.WithContext(ctx).Model(&model.Order{}).Order("create_at desc")
	if req.OrderID > 0 {
		tx = tx.Where("id = ?", req.OrderID)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	var total int64
	_ = tx.Count(&total).Error
	if err := tx.Offset(req.GetOffset()).Limit(req.Limit).Find(&orders).Error; err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	list := make([]*dto.OrderInfoResp, 0, len(orders))
	for _, order := range orders {
		list = append(list, s.orderInfoDto(ctx, &order))
	}
	return &dto.UserOrderListResp{Pager: req.Pager, Total: total, List: list}, common.OK
}

func (s *Service) AdminOrderInfo(ctx context.Context, orderID int64) (*dto.OrderInfoResp, common.Errno) {
	var order model.Order
	if err := s.db.WithContext(ctx).First(&order, orderID).Error; err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return s.orderInfoDto(ctx, &order), common.OK
}

func (s *Service) AdminOrderStats(ctx context.Context, req *dto.AdminOrderStatsReq) (*dto.AdminOrderStatsResp, common.Errno) {
	resp := &dto.AdminOrderStatsResp{ByStatus: []dto.StatusStat{}, ByGoods: []dto.GoodsStat{}}
	var orders []model.Order
	tx := s.db.WithContext(ctx)
	if req.CreateStart > 0 && req.CreateEnd > 0 {
		tx = tx.Where("create_at BETWEEN ? AND ?", req.CreateStart, req.CreateEnd)
	}
	if err := tx.Find(&orders).Error; err != nil {
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
	var items []model.OrderItem
	_ = s.db.WithContext(ctx).Find(&items).Error
	goodsMap := map[int64]*dto.GoodsStat{}
	for _, item := range items {
		stat := goodsMap[item.GoodsID]
		if stat == nil {
			stat = &dto.GoodsStat{GoodsID: item.GoodsID}
			var course model.CourseGood
			if err := s.db.WithContext(ctx).First(&course, item.GoodsID).Error; err == nil {
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

func formatOrderID(v int64) string {
	return fmt.Sprintf("%d", v)
}
