package dto

import "github.com/JunLang-7/mall/common"

type AddCategoryReq struct {
	Name     string `json:"name"`
	ParentID int64  `json:"parent_id"`
	Level    int32  `json:"level"`
	Sort     int32  `json:"sort"`
}

type UpdateCategoryReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DeleteCategoryReq struct {
	IDs []int64 `json:"ids"`
}

type UpdateSort struct {
	ID       int64 `json:"id"`
	Sort     int32 `json:"sort"`
	Level    int32 `json:"level"`
	ParentID int64 `json:"parent_id"`
}
type UpdateCategorySortReq []UpdateSort

type ListCategoryReq struct {
	common.Pager
}

type CategoryDto struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Sort     int32  `json:"sort"`
	ParentID int64  `json:"parent_id"`
	Level    int32  `json:"level"`
}
