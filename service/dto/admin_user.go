package dto

type AdminUserDto struct {
	UserID     int64  `json:"user_id"`
	Name       string `json:"name"`
	NickName   string `json:"nick_name"`
	Sex        int32  `json:"sex"`
	Status     int32  `json:"status"`
	Mobile     string `json:"mobile"`
	LarkOpenID string `json:"lark_open_id"`
	UpdateAt   int64  `json:"update_at"`
	CreateAt   int64  `json:"create_at"`
}

type GetUserInfoReq struct {
	ID int64 `form:"id" json:"id"`
}

type CreateUserReq struct {
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
	Mobile   string `json:"mobile"`
	Sex      int32  `json:"sex"`
}

type UpdateUserReq struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
	Sex      int32  `json:"sex"`
}

type UpdateUserStatusReq struct {
	ID     int64 `json:"id"`
	Status int32 `json:"status"`
}

type DeleteUserReq struct {
	ID int64 `json:"id"`
}
