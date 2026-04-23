package do

import "github.com/JunLang-7/mall/common"

type AddCategory struct {
	Name     string
	Level    int32
	ParentID int64
	Sort     int32
}

type UpdateCategory struct {
	ID   int64
	Name string
}

type DeleteCategory struct {
	ID int64
}

type UpdateSort struct {
	ID       int64
	Sort     int32
	Level    int32
	ParentID int64
}

type UpdateCategorySort []UpdateSort

type ListCategory struct {
	common.Pager
}
