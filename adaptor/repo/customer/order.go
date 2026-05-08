package customer

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"gorm.io/gorm"
)

type IOrder interface {
	CreateOrderWithItems(ctx context.Context, req *do.CreateOrderReq, items []*do.CreateOrderItemReq) error
	GetOrderByID(ctx context.Context, req *do.GetOrderReq) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, req *do.UpdateOrderStatusReq) error
	ListOrders(ctx context.Context, req *do.ListOrderReq) ([]*model.Order, int64, error)
	GetOrdersByTimeRange(ctx context.Context, createStart, createEnd int64) ([]*model.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]*model.OrderItem, error)
	GetAllOrderItems(ctx context.Context) ([]*model.OrderItem, error)
	CreateRefund(ctx context.Context, req *do.CreateRefundReq) error
	ListRefundsByOrderID(ctx context.Context, orderID int64) ([]*model.OrderRefund, error)
	DeliverOrder(ctx context.Context, orderID int64) error
	RefundOrder(ctx context.Context, orderID int64, req *do.CreateRefundReq, now int64) error
	CancelTimeoutOrders(ctx context.Context, timeout int64, now int64) error
	AutoReceiveOrders(ctx context.Context, shippedBefore int64, now int64) error
}

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(adaptor adaptor.IAdaptor) *OrderRepo {
	return &OrderRepo{db: adaptor.GetDB()}
}

func (r *OrderRepo) CreateOrderWithItems(ctx context.Context, req *do.CreateOrderReq, items []*do.CreateOrderItemReq) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		qs := query.Use(tx)
		if err := qs.Order.WithContext(ctx).Create(&model.Order{
			ID: req.ID, UserID: req.UserID, Status: req.Status, OrderSource: req.OrderSource,
			OrderAmount: req.OrderAmount, DiscountAmount: req.DiscountAmount, PaymentAmount: req.PaymentAmount,
			InnerTradeNo: req.InnerTradeNo, OrderDesc: req.OrderDesc, UserRemark: req.UserRemark,
			CreateAt: req.CreateAt, CreateBy: req.CreateBy,
		}); err != nil {
			return err
		}
		for _, item := range items {
			if err := qs.OrderItem.WithContext(ctx).Create(&model.OrderItem{
				OrderID: req.ID, UserID: req.UserID, GoodsID: item.GoodsID, GoodsType: item.GoodsType,
				Quantity: item.Quantity, PaymentAmount: item.PaymentAmount, DiscountAmount: item.DiscountAmount,
				GoodsSnap: item.GoodsSnap,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *OrderRepo) GetOrderByID(ctx context.Context, req *do.GetOrderReq) (*model.Order, error) {
	qs := query.Use(r.db).Order
	tx := qs.WithContext(ctx)
	if req.OrderID > 0 {
		tx = tx.Where(qs.ID.Eq(req.OrderID))
	}
	if req.UserID > 0 {
		tx = tx.Where(qs.UserID.Eq(req.UserID))
	}
	if req.Status != 0 {
		tx = tx.Where(qs.Status.Eq(req.Status))
	}
	return tx.First()
}

func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, req *do.UpdateOrderStatusReq) error {
	qs := query.Use(r.db).Order
	updateMap := map[string]interface{}{}
	if req.NewStatus != 0 {
		updateMap[qs.Status.ColumnName().String()] = req.NewStatus
	}
	if req.CancelAt > 0 {
		updateMap[qs.CancelAt.ColumnName().String()] = req.CancelAt
		updateMap[qs.CancelType.ColumnName().String()] = req.CancelType
		updateMap[qs.CancelBy.ColumnName().String()] = req.CancelBy
		updateMap[qs.CancelReason.ColumnName().String()] = req.CancelReason
	}
	if req.PaymentAt > 0 {
		updateMap[qs.PaymentAt.ColumnName().String()] = req.PaymentAt
	}
	if req.RefundAmount > 0 {
		updateMap[qs.RefundAmount.ColumnName().String()] = req.RefundAmount
		updateMap[qs.RefundAt.ColumnName().String()] = req.RefundAt
	}
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.OrderID)).Updates(updateMap)
	return err
}

func (r *OrderRepo) ListOrders(ctx context.Context, req *do.ListOrderReq) ([]*model.Order, int64, error) {
	qs := query.Use(r.db).Order
	tx := qs.WithContext(ctx).Order(qs.CreateAt.Desc())
	if req.UserID > 0 {
		tx = tx.Where(qs.UserID.Eq(req.UserID))
	}
	if req.Status != 0 {
		tx = tx.Where(qs.Status.Eq(req.Status))
	}
	return tx.FindByPage(req.GetOffset(), req.Limit)
}

func (r *OrderRepo) GetOrdersByTimeRange(ctx context.Context, createStart, createEnd int64) ([]*model.Order, error) {
	qs := query.Use(r.db).Order
	tx := qs.WithContext(ctx)
	if createStart > 0 && createEnd > 0 {
		tx = tx.Where(qs.CreateAt.Between(createStart, createEnd))
	}
	return tx.Find()
}

func (r *OrderRepo) GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]*model.OrderItem, error) {
	qs := query.Use(r.db).OrderItem
	return qs.WithContext(ctx).Where(qs.OrderID.Eq(orderID)).Find()
}

func (r *OrderRepo) GetAllOrderItems(ctx context.Context) ([]*model.OrderItem, error) {
	qs := query.Use(r.db).OrderItem
	return qs.WithContext(ctx).Find()
}

func (r *OrderRepo) CreateRefund(ctx context.Context, req *do.CreateRefundReq) error {
	qs := query.Use(r.db).OrderRefund
	return qs.WithContext(ctx).Create(&model.OrderRefund{
		UserID:       req.UserID,
		OrderID:      req.OrderID,
		ItemIds:      req.ItemIds,
		ApplyAt:      req.ApplyAt,
		Reason:       req.Reason,
		Status:       req.Status,
		Amount:       req.Amount,
		InnerTradeNo: req.InnerTradeNo,
		RefundID:     req.RefundID,
		ApplyUserID:  req.ApplyUserID,
		DoneAt:       req.DoneAt,
	})
}

func (r *OrderRepo) ListRefundsByOrderID(ctx context.Context, orderID int64) ([]*model.OrderRefund, error) {
	qs := query.Use(r.db).OrderRefund
	return qs.WithContext(ctx).Where(qs.OrderID.Eq(orderID)).Find()
}

func (r *OrderRepo) DeliverOrder(ctx context.Context, orderID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		qs := query.Use(tx)
		order, err := qs.Order.WithContext(ctx).Where(qs.Order.ID.Eq(orderID)).First()
		if err != nil {
			return err
		}
		if order.Status == consts.OrderStatusDone {
			return nil
		}
		items, err := qs.OrderItem.WithContext(ctx).Where(qs.OrderItem.OrderID.Eq(order.ID)).Find()
		if err != nil {
			return err
		}
		now := time.Now().UnixMilli()
		for _, item := range items {
			if item.GoodsType != consts.GoodsTypeCourse {
				continue
			}
			if err := tx.WithContext(ctx).Where("user_id = ? AND order_item_id = ?", order.UserID, item.ID).
				FirstOrCreate(&model.UserCourseGood{
					UserID:            order.UserID,
					OrderID:           order.ID,
					OrderItemID:       item.ID,
					GoodsID:           item.GoodsID,
					GoodsType:         item.GoodsType,
					BuyTime:           now,
					LearnExpireTime:   addMonthByCode(now, 12),
					ServiceExpireTime: addMonthByCode(now, 12),
				}).Error; err != nil {
				return err
			}
			_ = tx.WithContext(ctx).Where("user_id = ? AND goods_id = ?", order.UserID, item.GoodsID).Delete(&model.UserCart{}).Error
		}
		_, err = qs.Order.WithContext(ctx).Where(qs.Order.ID.Eq(order.ID)).Updates(map[string]interface{}{
			qs.Order.Status.ColumnName().String():    consts.OrderStatusDone,
			qs.Order.PaymentAt.ColumnName().String(): now,
		})
		return err
	})
}

func (r *OrderRepo) RefundOrder(ctx context.Context, orderID int64, req *do.CreateRefundReq, now int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		qs := query.Use(tx)
		if err := qs.OrderRefund.WithContext(ctx).Create(&model.OrderRefund{
			UserID:       req.UserID,
			OrderID:      orderID,
			ItemIds:      req.ItemIds,
			ApplyAt:      now,
			Reason:       req.Reason,
			Status:       consts.RefundStatusDone,
			Amount:       req.Amount,
			InnerTradeNo: req.InnerTradeNo,
			RefundID:     req.RefundID,
			ApplyUserID:  req.ApplyUserID,
			DoneAt:       now,
		}); err != nil {
			return err
		}
		_, err := qs.Order.WithContext(ctx).Where(qs.Order.ID.Eq(orderID)).Updates(map[string]interface{}{
			qs.Order.Status.ColumnName().String():       consts.OrderStatusRefunded,
			qs.Order.RefundAmount.ColumnName().String(): req.Amount,
			qs.Order.RefundAt.ColumnName().String():     now,
		})
		if err != nil {
			return err
		}
		_, err = qs.UserCourseGood.WithContext(ctx).Where(qs.UserCourseGood.OrderID.Eq(orderID)).Updates(map[string]interface{}{
			qs.UserCourseGood.LearnExpireTime.ColumnName().String():   now,
			qs.UserCourseGood.ServiceExpireTime.ColumnName().String(): now,
		})
		return err
	})
}

func (r *OrderRepo) CancelTimeoutOrders(ctx context.Context, timeout int64, now int64) error {
	qs := query.Use(r.db).Order
	_, err := qs.WithContext(ctx).Where(
		qs.Status.Eq(consts.OrderStatusPending),
		qs.CreateAt.Lt(timeout),
	).Updates(map[string]interface{}{
		qs.Status.ColumnName().String():       consts.OrderStatusCanceled,
		qs.CancelAt.ColumnName().String():     now,
		qs.CancelType.ColumnName().String():   consts.CancelTypeTimeout,
		qs.CancelBy.ColumnName().String():     consts.SystemUserID,
		qs.CancelReason.ColumnName().String(): "timeout",
	})
	return err
}

func (r *OrderRepo) AutoReceiveOrders(ctx context.Context, shippedBefore int64, now int64) error {
	qs := query.Use(r.db).Order
	_, err := qs.WithContext(ctx).Where(
		qs.Status.Eq(consts.OrderStatusShipped),
		qs.PaymentAt.Lt(shippedBefore),
	).Updates(map[string]interface{}{
		qs.Status.ColumnName().String():               consts.OrderStatusSigned,
		qs.ReceiverConfirmAt.ColumnName().String():     now,
		qs.ReceiverConfirmType.ColumnName().String():   consts.ReceiveConfirmAuto,
	})
	return err
}

func addMonthByCode(ms int64, months int) int64 {
	return time.UnixMilli(ms).AddDate(0, months, 0).UnixMilli()
}
