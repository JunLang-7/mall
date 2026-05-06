package dto

import "github.com/JunLang-7/mall/common"

type CustomerUserDto struct {
	UserID      int64  `json:"user_id"`
	NickName    string `json:"nick_name"`
	CreateAt    int64  `json:"create_at"`
	IconURL     string `json:"icon_url"`
	Sex         int32  `json:"sex"`
	Status      int32  `json:"status"`
	LastLoginAt int64  `json:"last_login_at"`
	UpdateAt    int64  `json:"update_at"`
	WechatBind  bool   `json:"wechat_bind"`
	HasPassword bool   `json:"has_password"`
}

type CustomerMobileUserDto struct {
	Mobile string `json:"mobile"`
	UserID int64  `json:"user_id"`
}

type CustomerWechatUserDto struct {
	UserID  int64  `json:"user_id"`
	UnionID string `json:"union_id"`
}

type CustomerAppUserDto struct {
	OpenID  string `json:"open_id"`
	UserID  int64  `json:"user_id"`
	AppCode int32  `json:"app_code"`
}

type CustomerUserInfoDto struct {
	User       *CustomerUserDto         `json:"user"`
	MobileUser *CustomerMobileUserDto   `json:"mobile_user,omitempty"`
	WechatUser *CustomerWechatUserDto   `json:"wechat_user,omitempty"`
	AppUsers   []*CustomerAppUserDto    `json:"app_users"`
}

type CustomerLoginResp struct {
	Token    string               `json:"token"`
	UserInfo *CustomerUserInfoDto `json:"user_info"`
}

type ChangePasswordReq struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
	VerifyCode      string `json:"verify_code"`
}

type ChangePasswordSmsCodeReq struct {
	Ticket string `json:"ticket"`
}

type ChangePasswordResp struct {
	ReloginRequired bool `json:"relogin_required"`
}

type WechatQrCodeReq struct {
	Purpose string `json:"purpose"`
}

type WechatQrCodeResp struct {
	ExpireIn   int64  `json:"expire_in"`
	SceneToken string `json:"scene_token"`
	QrcodeURL  string `json:"qrcode_url"`
}

type WechatQrCodeStatusResp struct {
	State     string               `json:"state"`
	Purpose   string               `json:"purpose"`
	Message   string               `json:"message"`
	Token     string               `json:"token,omitempty"`
	UserInfo  *CustomerUserInfoDto `json:"user_info,omitempty"`
	BindExist bool                 `json:"bind_exist"`
}

type WechatScanConfirmReq struct {
	SceneToken string `json:"scene_token"`
	Code       string `json:"code"`
}

type CustomerCourseDetailDto struct {
	CourseDto
	TotalDuration int64         `json:"total_duration"`
	LessonCount   int32         `json:"lesson_count"`
	Catalogs      []*CatalogDto `json:"catalogs"`
	HasPurchased bool          `json:"has_purchased"`
}

type PurchasedCourseDto struct {
	ID                int64    `json:"id"`
	Name              string   `json:"name"`
	ServiceExpireTime int64    `json:"service_expire_time"`
	LearnExpireTime   int64    `json:"learn_expire_time"`
	Features          []string `json:"features"`
	UpdateStatus      int32    `json:"update_status"`
	HasPurchased      bool     `json:"has_purchased"`
	CoverKey          string   `json:"cover_key"`
	CoverURL          string   `json:"cover_url"`
	DetailCoverKey    string   `json:"detail_cover_key"`
	DetailCoverURL    string   `json:"detail_cover_url"`
	Detail            string   `json:"detail"`
}

type PurchasedCourseListResp struct {
	common.Pager
	List  []*PurchasedCourseDto `json:"list"`
	Total int64                 `json:"total"`
}

type LessonLearnInfoReq struct {
	CourseID int64 `form:"course_id"`
	LessonID int64 `form:"lesson_id"`
}

type LessonLearnInfoResp struct {
	CourseID       int64 `json:"course_id"`
	LessonID       int64 `json:"lesson_id"`
	PlayPosition   int64 `json:"play_position"`
	LearnStatus    int32 `json:"learn_status"`
	LastType       int32 `json:"last_type"`
	EntryTime      int64 `json:"entry_time"`
	LastReportTime int64 `json:"last_report_time"`
	InLearning     bool  `json:"in_learning"`
}

type LessonLearnReportReq struct {
	CourseID      int64 `json:"course_id"`
	LessonID      int64 `json:"lesson_id"`
	Type          int32 `json:"type"`
	PlayPosition  int64 `json:"play_position"`
}

type ContinueLearnCourseDto struct {
	CourseID       int64  `json:"course_id"`
	CourseName     string `json:"course_name"`
	CourseCoverKey string `json:"course_cover_key"`
	CourseCoverURL string `json:"course_cover_url"`
	LessonID       int64  `json:"lesson_id"`
	LessonName     string `json:"lesson_name"`
	LessonIndex    int32  `json:"lesson_index"`
	LessonCount    int64  `json:"lesson_count"`
	PlayPosition   int64  `json:"play_position"`
	LearnStatus    int32  `json:"learn_status"`
	LastLearnTime  int64  `json:"last_learn_time"`
}

type ContinueLearnResp struct {
	common.Pager
	List  []*ContinueLearnCourseDto `json:"list"`
	Total int64                     `json:"total"`
}

type AddGoodsReq struct {
	GoodsID int64 `json:"goods_id"`
}

type UpdateCartGoodsReq struct {
	ID       int64 `json:"id"`
	Quantity int32 `json:"quantity"`
}

type RemoveGoodsReq struct {
	ID int64 `json:"id"`
}

type ListCartGoodsReq struct {
	common.Pager
	GoodsNameKW string `form:"goods_name_kw"`
}

type CartGoodsDto struct {
	CourseDto
	CartID   int64 `json:"cart_id"`
	GoodsID  int64 `json:"goods_id"`
	Quantity int32 `json:"quantity"`
}

type ListGoodsResp struct {
	common.Pager
	Total int64           `json:"total"`
	List  []*CartGoodsDto `json:"list"`
}

type OrderCalcFeeReq struct {
	CourseIDs []int64 `json:"course_ids"`
	Platform  string  `json:"platform"`
}

type CourseFeeDto struct {
	CourseID    int64 `json:"course_id"`
	Price       int64 `json:"price"`
	DiscountFee int64 `json:"discount_fee"`
	PayFee      int64 `json:"pay_fee"`
	GoodsSnap   any   `json:"goods_snap"`
}

type OrderCalcFeeResp struct {
	FeeUUID          string          `json:"fee_uuid"`
	TotalFee         int64           `json:"total_fee"`
	TotalDiscountFee int64           `json:"total_discount_fee"`
	TotalPayFee      int64           `json:"total_pay_fee"`
	ExpireTime       int64           `json:"expire_time"`
	CourseFees       []*CourseFeeDto `json:"course_fees"`
}

type OrderPayNowReq struct {
	FeeUUID  string `json:"fee_uuid"`
	Remark   string `json:"remark"`
	Platform string `json:"platform"`
}

type OrderPayLaterReq struct {
	OrderID  int64  `json:"order_id"`
	Platform string `json:"platform"`
}

type OrderPayNowResp struct {
	OrderID   int64  `json:"order_id"`
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
	CodeURL   string `json:"code_url"`
	TradeType string `json:"trade_type"`
}

type CancelOrderReq struct {
	OrderID int64  `json:"order_id"`
	Reason  string `json:"reason"`
}

type OrderListReq struct {
	common.Pager
	OrderID      int64  `form:"order_id"`
	Status       int32  `form:"status"`
	StatusList   string `form:"status_list"`
	GoodsNameKW  string `form:"goods_name_kw"`
	CreateStart  int64  `form:"create_start"`
	CreateEnd    int64  `form:"create_end"`
	PaymentStart int64  `form:"payment_start"`
	PaymentEnd   int64  `form:"payment_end"`
	RefundStart  int64  `form:"refund_start"`
	RefundEnd    int64  `form:"refund_end"`
}

type OrderInfoReq struct {
	OrderID int64 `form:"order_id"`
}

type OrderDto struct {
	ID                  int64   `json:"id"`
	UserID              int64   `json:"user_id"`
	Status              int32   `json:"status"`
	OrderSource         int32   `json:"order_source"`
	OrderAmount         int64   `json:"order_amount"`
	DiscountAmount      int64   `json:"discount_amount"`
	PaymentAmount       int64   `json:"payment_amount"`
	TradeNo             string  `json:"trade_no"`
	InnerTradeNo        string  `json:"inner_trade_no"`
	OrderDesc           string  `json:"order_desc"`
	PaymentAt           int64   `json:"payment_at"`
	UserRemark          string  `json:"user_remark"`
	ReceiverConfirmAt   *int64  `json:"receiver_confirm_at"`
	ReceiverConfirmType *int32  `json:"receiver_confirm_type"`
	RefundAmount        int64   `json:"refund_amount"`
	RefundAt            *int64  `json:"refund_at"`
	CancelAt            *int64  `json:"cancel_at"`
	CancelType          *int32  `json:"cancel_type"`
	CancelBy            *int64  `json:"cancel_by"`
	CancelReason        *string `json:"cancel_reason"`
	CreateAt            int64   `json:"create_at"`
	CreateBy            int64   `json:"create_by"`
	CancelName          string  `json:"cancel_name"`
	CreateName          string  `json:"create_name"`
	UserName            string  `json:"user_name"`
	UserMobile          string  `json:"user_mobile"`
}

type OrderItemDto struct {
	ID             int64  `json:"id"`
	OrderID        int64  `json:"order_id"`
	UserID         int64  `json:"user_id"`
	GoodsID        int64  `json:"goods_id"`
	GoodsType      int32  `json:"goods_type"`
	Quantity       int32  `json:"quantity"`
	PaymentAmount  int64  `json:"payment_amount"`
	DiscountAmount int64  `json:"discount_amount"`
	GoodsSnap      any    `json:"goods_snap"`
	RefundStatus   int32  `json:"refund_status"`
}

type RefundDto struct {
	ID          int64           `json:"id"`
	Amount      int64           `json:"amount"`
	ItemIDs     []int64         `json:"item_ids"`
	ApplyAt     int64           `json:"apply_at"`
	Status      int32           `json:"status"`
	DoneAt      int64           `json:"done_at"`
	Reason      string          `json:"reason"`
	RefundID    string          `json:"refund_id"`
	ApplyUserID int64           `json:"apply_user_id"`
	RefundName  string          `json:"refund_name"`
	Items       []*OrderItemDto `json:"items"`
}

type OrderInfoResp struct {
	OrderDto
	Items   []*OrderItemDto `json:"items"`
	Refunds []*RefundDto    `json:"refunds"`
}

type UserOrderListResp struct {
	common.Pager
	Total int64            `json:"total"`
	List  []*OrderInfoResp `json:"list"`
}

type PaymentQueryReq struct {
	OrderID int64 `form:"order_id"`
}

type RefundOrderReq struct {
	OrderID int64   `json:"order_id"`
	ItemIDs []int64 `json:"item_ids"`
	Amount  int64   `json:"amount"`
	Reason  string  `json:"reason"`
}

type AdminCustomerUserListReq struct {
	common.Pager
	UserID int64  `form:"user_id"`
	Mobile string `form:"mobile"`
	Status int32  `form:"status"`
}

type AdminCustomerUserListResp struct {
	common.Pager
	Total int64                  `json:"total"`
	List  []*CustomerUserInfoDto `json:"list"`
}

type AdminCustomerUserStatusReq struct {
	UserID int64 `json:"user_id"`
	Status int32 `json:"status"`
}

type AdminOrderStatsReq struct {
	CreateStart int64 `form:"create_start"`
	CreateEnd   int64 `form:"create_end"`
}

type AdminOrderStatsResp struct {
	ByStatus []StatusStat `json:"by_status"`
	ByGoods  []GoodsStat  `json:"by_goods"`
	TotalPay int64        `json:"total_pay"`
}

type StatusStat struct {
	Status int32 `json:"status"`
	Count  int64 `json:"count"`
	Amount int64 `json:"amount"`
}

type GoodsStat struct {
	GoodsID int64  `json:"goods_id"`
	Name    string `json:"name"`
	Count   int64  `json:"count"`
	Amount  int64  `json:"amount"`
}
