package do

import "github.com/JunLang-7/mall/common"

type CreateUser struct {
	AdminUserID int64  `json:"admin_user_id"`
	Name        string `json:"name"`
	NickName    string `json:"nick_name"`
	Mobile      string `json:"mobile"`
	Sex         int32  `json:"sex"`
}

type UpdateUser struct {
	AdminUserID int64  `json:"admin_user_id"`
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	NickName    string `json:"nick_name"`
	Sex         int32  `json:"sex"`
}

type UpdateUserStatus struct {
	AdminUserID int64 `json:"admin_user_id"`
	ID          int64 `json:"id"`
	Status      int32 `json:"status"`
}

type UpdateUserPassword struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
}

type ListUsers struct {
	Name   string       `json:"name"`
	Mobile string       `json:"mobile"`
	RoleID int64        `json:"role_id"`
	Status int32        `json:"status"`
	Pager  common.Pager `json:"pager"`
}
