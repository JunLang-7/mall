package do

import "github.com/JunLang-7/mall/common"

type AddRole struct {
	AdminUserID int64
	Name        string
	Desc        string
}

type UpdateRole struct {
	AdminUserID int64
	ID          int64
	Name        string
	Desc        string
	Status      int32
}

type ListRole struct {
	common.Pager
	NameKw string
	Status int32
}
