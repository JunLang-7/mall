package dto

import "github.com/JunLang-7/mall/common"

type PermissionDto struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Type     int32  `json:"type"` // 1:  2
	Name     string `json:"name"`
	PagePath string `json:"page_path"`
	ParentID int64  `json:"parent_id"` // ID
	Status   int32  `json:"status"`    // 1 -1
	Sort     int32  `json:"sort"`
	Desc     string `json:"desc"`
}

type PermissionListResp struct {
	Pager common.Pager `json:"pager"`
	Total int64        `json:"total"`
	List  []*PermissionDto
}
