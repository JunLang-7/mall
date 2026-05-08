package customer

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/service/do"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUser interface {
	GetUserPurchasedCourses(ctx context.Context, userID int64) ([]*model.UserCourseGood, error)
	HasPurchasedCourse(ctx context.Context, req *do.HasPurchasedReq) (bool, error)
	CreateUserCourseGood(ctx context.Context, req *do.UserCourseGoodCreate) error
	AddToCart(ctx context.Context, req *do.AddCartReq) (int64, error)
	RemoveFromCart(ctx context.Context, req *do.RemoveCartReq) error
	ListCart(ctx context.Context, req *do.ListCartReq) ([]*model.UserCart, int64, error)
	DeleteCartByUserAndGoodsID(ctx context.Context, userID, goodsID int64) error
	GetUserInfo(ctx context.Context, userID int64) (*model.User, error)
	ListUsers(ctx context.Context, req *do.CustomerUserListReq) ([]*model.User, int64, error)
	UpdateUserStatus(ctx context.Context, req *do.CustomerUserStatusReq) error
	GetMobileUserByHash(ctx context.Context, mobileHash string) (*model.MobileUser, error)
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	UpdateUserPassword(ctx context.Context, userID int64, password string) error
	UpdateUserLastLoginAt(ctx context.Context, userID int64, lastLoginAt time.Time) error
	CreateUserWithMobileUser(ctx context.Context, user *model.User, mobileUser *model.MobileUser) error
	GetUserByMobileHash(ctx context.Context, mobileHash string) (*model.User, error)
	GetMobileUserByUserID(ctx context.Context, userID int64) (*model.MobileUser, error)
	GetWechatUserByUserID(ctx context.Context, userID int64) (*model.WechatUser, error)
	GetWechatUserByUnionID(ctx context.Context, unionID string) (*model.WechatUser, error)
	GetAppUsersByUserID(ctx context.Context, userID int64) ([]*model.AppUser, error)
	GetAppUserByOpenID(ctx context.Context, openID string, appCode int32) (*model.AppUser, error)
	CreateUserWithWechatUser(ctx context.Context, user *model.User, wechatUser *model.WechatUser) error
	CreateUserWithAppUser(ctx context.Context, user *model.User, appUser *model.AppUser) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(adaptor adaptor.IAdaptor) *UserRepo {
	return &UserRepo{db: adaptor.GetDB()}
}

func (r *UserRepo) GetUserPurchasedCourses(ctx context.Context, userID int64) ([]*model.UserCourseGood, error) {
	qs := query.Use(r.db).UserCourseGood
	return qs.WithContext(ctx).Where(qs.UserID.Eq(userID), qs.LearnExpireTime.Gt(time.Now().UnixMilli())).Find()
}

func (r *UserRepo) HasPurchasedCourse(ctx context.Context, req *do.HasPurchasedReq) (bool, error) {
	qs := query.Use(r.db).UserCourseGood
	count, err := qs.WithContext(ctx).Where(
		qs.UserID.Eq(req.UserID),
		qs.GoodsID.Eq(req.CourseID),
		qs.LearnExpireTime.Gt(time.Now().UnixMilli()),
	).Count()
	return count > 0, err
}

func (r *UserRepo) CreateUserCourseGood(ctx context.Context, req *do.UserCourseGoodCreate) error {
	qs := query.Use(r.db).UserCourseGood
	return qs.WithContext(ctx).Create(&model.UserCourseGood{
		UserID: req.UserID, OrderID: req.OrderID, OrderItemID: req.OrderItemID,
		GoodsID: req.GoodsID, GoodsType: req.GoodsType, BuyTime: req.BuyTime,
		LearnExpireTime: req.LearnExpireTime, ServiceExpireTime: req.ServiceExpireTime,
	})
}

func (r *UserRepo) AddToCart(ctx context.Context, req *do.AddCartReq) (int64, error) {
	now := time.Now()
	qs := query.Use(r.db).UserCart
	cart := model.UserCart{UserID: req.UserID, GoodsID: req.GoodsID, Quantity: req.Quantity, AddAt: now}
	err := qs.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: qs.UserID.ColumnName().String()},
			{Name: qs.GoodsID.ColumnName().String()},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			qs.Quantity.ColumnName().String(): req.Quantity,
			qs.AddAt.ColumnName().String():    now,
		}),
	}).Create(&cart)
	return cart.ID, err
}

func (r *UserRepo) RemoveFromCart(ctx context.Context, req *do.RemoveCartReq) error {
	qs := query.Use(r.db).UserCart
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID), qs.UserID.Eq(req.UserID)).Delete()
	return err
}

func (r *UserRepo) ListCart(ctx context.Context, req *do.ListCartReq) ([]*model.UserCart, int64, error) {
	qs := query.Use(r.db).UserCart
	tx := qs.WithContext(ctx).Where(qs.UserID.Eq(req.UserID)).Order(qs.AddAt.Desc())
	return tx.FindByPage(req.GetOffset(), req.Limit)
}

func (r *UserRepo) DeleteCartByUserAndGoodsID(ctx context.Context, userID, goodsID int64) error {
	qs := query.Use(r.db).UserCart
	_, err := qs.WithContext(ctx).Where(qs.UserID.Eq(userID), qs.GoodsID.Eq(goodsID)).Delete()
	return err
}

func (r *UserRepo) GetUserInfo(ctx context.Context, userID int64) (*model.User, error) {
	qs := query.Use(r.db).User
	return qs.WithContext(ctx).Where(qs.ID.Eq(userID)).First()
}

func (r *UserRepo) ListUsers(ctx context.Context, req *do.CustomerUserListReq) ([]*model.User, int64, error) {
	qs := query.Use(r.db).User
	tx := qs.WithContext(ctx)
	if req.UserID > 0 {
		tx = tx.Where(qs.ID.Eq(req.UserID))
	}
	if req.Status != 0 {
		tx = tx.Where(qs.Status.Eq(req.Status))
	}
	return tx.FindByPage(req.GetOffset(), req.Limit)
}

func (r *UserRepo) UpdateUserStatus(ctx context.Context, req *do.CustomerUserStatusReq) error {
	qs := query.Use(r.db).User
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.UserID)).Update(qs.Status, req.Status)
	return err
}

func (r *UserRepo) GetMobileUserByHash(ctx context.Context, mobileHash string) (*model.MobileUser, error) {
	qs := query.Use(r.db).MobileUser
	return qs.WithContext(ctx).Where(qs.MobileSha256.Eq(mobileHash)).First()
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	qs := query.Use(r.db).User
	return qs.WithContext(ctx).Where(qs.ID.Eq(userID)).First()
}

func (r *UserRepo) UpdateUserPassword(ctx context.Context, userID int64, password string) error {
	qs := query.Use(r.db).User
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(userID)).Update(qs.Password, password)
	return err
}

func (r *UserRepo) UpdateUserLastLoginAt(ctx context.Context, userID int64, lastLoginAt time.Time) error {
	qs := query.Use(r.db).User
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(userID)).Update(qs.LastLoginAt, lastLoginAt)
	return err
}

func (r *UserRepo) CreateUserWithMobileUser(ctx context.Context, user *model.User, mobileUser *model.MobileUser) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		mobileUser.UserID = user.ID
		return tx.Create(mobileUser).Error
	})
}

func (r *UserRepo) GetUserByMobileHash(ctx context.Context, mobileHash string) (*model.User, error) {
	mqs := query.Use(r.db).MobileUser
	uqs := query.Use(r.db).User
	mobileUser, err := mqs.WithContext(ctx).Where(mqs.MobileSha256.Eq(mobileHash)).First()
	if err != nil {
		return nil, err
	}
	return uqs.WithContext(ctx).Where(uqs.ID.Eq(mobileUser.UserID)).First()
}

func (r *UserRepo) GetMobileUserByUserID(ctx context.Context, userID int64) (*model.MobileUser, error) {
	qs := query.Use(r.db).MobileUser
	return qs.WithContext(ctx).Where(qs.UserID.Eq(userID)).First()
}

func (r *UserRepo) GetWechatUserByUserID(ctx context.Context, userID int64) (*model.WechatUser, error) {
	qs := query.Use(r.db).WechatUser
	return qs.WithContext(ctx).Where(qs.UserID.Eq(userID)).First()
}

func (r *UserRepo) GetAppUsersByUserID(ctx context.Context, userID int64) ([]*model.AppUser, error) {
	qs := query.Use(r.db).AppUser
	return qs.WithContext(ctx).Where(qs.UserID.Eq(userID)).Find()
}

func (r *UserRepo) GetAppUserByOpenID(ctx context.Context, openID string, appCode int32) (*model.AppUser, error) {
	qs := query.Use(r.db).AppUser
	return qs.WithContext(ctx).Where(qs.OpenID.Eq(openID), qs.AppCode.Eq(appCode)).First()
}

func (r *UserRepo) CreateUserWithAppUser(ctx context.Context, user *model.User, appUser *model.AppUser) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		appUser.UserID = user.ID
		return tx.Create(appUser).Error
	})
}

func (r *UserRepo) GetWechatUserByUnionID(ctx context.Context, unionID string) (*model.WechatUser, error) {
	qs := query.Use(r.db).WechatUser
	return qs.WithContext(ctx).Where(qs.UnionID.Eq(unionID)).First()
}

func (r *UserRepo) CreateUserWithWechatUser(ctx context.Context, user *model.User, wechatUser *model.WechatUser) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		wechatUser.UserID = user.ID
		return tx.Create(wechatUser).Error
	})
}
