package do

import (
	"time"

	"github.com/JunLang-7/mall/common"
)

// Course 入参
type CourseListReq struct {
	ID     int64
	NameKW string
	Status int32
	common.Pager
}

type CourseInfoReq struct {
	ID     int64
	Status int32
}

// Lesson 入参
type LessonLearnProgressReq struct {
	UserID   int64
	CourseID int64
	LessonID int64
}

type LessonLearnProgressUpdate struct {
	UserID       int64
	CourseID     int64
	LessonID     int64
	PlayPosition int64
	LearnStatus  int32
}

type LessonLearnRecordCreate struct {
	UserID    int64
	CourseID  int64
	LessonID  int64
	EntryTime time.Time
	ExitTime  time.Time
	Duration  int32
	LastType  int32
}

// 用户课程权限入参
type UserCourseGoodCreate struct {
	UserID            int64
	OrderID           int64
	OrderItemID       int64
	GoodsID           int64
	GoodsType         int32
	BuyTime           int64
	LearnExpireTime   int64
	ServiceExpireTime int64
}

type HasPurchasedReq struct {
	UserID   int64
	CourseID int64
}

// 购物车入参
type AddCartReq struct {
	UserID   int64
	GoodsID  int64
	Quantity int32
}

type RemoveCartReq struct {
	ID     int64
	UserID int64
}

type ListCartReq struct {
	UserID int64
	common.Pager
}

// 订单入参
type CreateOrderReq struct {
	ID             int64
	UserID         int64
	Status         int32
	OrderSource    int32
	OrderAmount    int64
	DiscountAmount int64
	PaymentAmount  int64
	InnerTradeNo   string
	OrderDesc      string
	UserRemark     string
	CreateAt       int64
	CreateBy       int64
}

type CreateOrderItemReq struct {
	OrderID        int64
	UserID         int64
	GoodsID        int64
	GoodsType      int32
	Quantity       int32
	PaymentAmount  int64
	DiscountAmount int64
	GoodsSnap      string
}

type GetOrderReq struct {
	OrderID int64
	UserID  int64
	Status  int32
}

type UpdateOrderStatusReq struct {
	OrderID      int64
	NewStatus    int32
	CancelAt     int64
	CancelType   int32
	CancelBy     int64
	CancelReason string
	PaymentAt    int64
	RefundAmount int64
	RefundAt     int64
}

type ListOrderReq struct {
	UserID int64
	Status int32
	common.Pager
}

// 退款入参
type CreateRefundReq struct {
	UserID       int64
	OrderID      int64
	ItemIds      string
	ApplyAt      int64
	Reason       string
	Status       int32
	Amount       int64
	InnerTradeNo string
	RefundID     string
	ApplyUserID  int64
	DoneAt       int64
}

// 用户管理入参
type CustomerUserListReq struct {
	UserID int64
	Status int32
	common.Pager
}

type CustomerUserStatusReq struct {
	UserID int64
	Status int32
}
