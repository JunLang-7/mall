package dto

import "github.com/JunLang-7/mall/common"

type CreateCourseReq struct {
	Name           string   `json:"name"`
	CoursePrice    int64    `json:"course_price"`
	ServiceTime    int32    `json:"service_time"`
	LearnTime      int32    `json:"learn_time"`
	Sort           int32    `json:"sort"`
	Features       []string `json:"features"`
	UpdateStatus   int32    `json:"update_status"`
	CoverKey       string   `json:"cover_key"`
	DetailCoverKey string   `json:"detail_cover_key"`
	Detail         string   `json:"detail"`
}

type CourseDto struct {
	common.CreateUpdateName
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	CoursePrice    int64    `json:"course_price"`
	ServiceTime    int32    `json:"service_time"`
	LearnTime      int32    `json:"learn_time"`
	Status         int32    `json:"status"`
	Sort           int32    `json:"sort"`
	Features       []string `json:"features"`
	UpdateStatus   int32    `json:"update_status"`
	HasPurchased   bool     `json:"has_purchased"`
	CoverKey       string   `json:"cover_key"`
	CoverURL       string   `json:"cover_url"`
	DetailCoverKey string   `json:"detail_cover_key"`
	DetailCOverURL string   `json:"detail_cover_url"`
	Detail         string   `json:"detail"`
	CreateBy       int64    `json:"create_by"`
	CreateAt       int64    `json:"create_at"`
	UpdateBy       int64    `json:"update_by"`
	UpdateAt       int64    `json:"update_at"`
}

type CourseInfoReq struct {
	ID int64 `form:"id"`
}

type UpdateCourseReq struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	CoursePrice    int64    `json:"course_price"`
	ServiceTime    int32    `json:"service_time"`
	LearnTime      int32    `json:"learn_time"`
	Sort           int32    `json:"sort"`
	Features       []string `json:"features"`
	UpdateStatus   int32    `json:"update_status"`
	CoverKey       string   `json:"cover_key"`
	DetailCoverKey string   `json:"detail_cover_key"`
	Detail         string   `json:"detail"`
}

type UpdateCourseStatusReq struct {
	ID     int64 `form:"id"`
	Status int32 `form:"status"`
}

type CourseListReq struct {
	common.Pager
	ID              int64  `form:"id"`
	NameKW          string `form:"name_kw"`
	CreateStartTime int64  `form:"create_start_time"`
	CreateEndTime   int64  `form:"create_end_time"`
	UpdateStartTime int64  `form:"update_start_time"`
	UpdateEndTime   int64  `form:"update_end_time"`
	UpdateStatus    int32  `form:"update_status"`
	Status          int32  `form:"status"`
	IsRecommend     bool   `form:"is_recommend"`
}

type CourseListResp struct {
	common.Pager
	List  []*CourseDto `json:"list"`
	Total int64        `json:"total"`
}

type AddCatalogReq struct {
	CourseID int64  `json:"course_id"`
	Name     string `json:"name"`
	ParentID int64  `json:"parent_id"`
	Sort     int32  `json:"sort"`
	Level    int32  `json:"level"`
}

type UpdateCatalogReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DeleteCatalogReq struct {
	ID int64 `json:"id"`
}

type UpdateCatalogSortDto struct {
	ID       int64            `json:"id"`
	Sort     int32            `json:"sort"`
	ParentID int64            `json:"parent_id"`
	Level    int32            `json:"level"`
	Lessons  []*common.IDSort `json:"lessons"`
}

type CatalogLessonDto struct {
	ID            int64  `json:"id"`
	LessonID      int64  `json:"lesson_id"`
	Name          string `json:"name"`
	LessonName    string `json:"lesson_name"`
	Detail        string `json:"detail"`
	VideoURL      string `json:"video_url"`
	VideoFileName string `json:"video_file_name"`
	Duration      int32  `json:"duration"`
	Status        int32  `json:"status"`
	ShowTime      int64  `json:"show_time"`
	EnableTrial   bool   `json:"enable_trial"`
}

type CatalogInfoReq struct {
	CourseID int64 `form:"course_id"`
}

type CatalogDto struct {
	ID          int64               `json:"id"`
	ParentID    int64               `json:"parent_id"`
	Level       int32               `json:"level"`
	Name        string              `json:"name"`
	CourseID    int64               `json:"course_id"`
	Sort        int32               `json:"sort"`
	Lessons     []*CatalogLessonDto `json:"lessons"`
	LessonCount int32               `json:"lesson_count"`
}

type CatalogInfoResp struct {
	TotalDuration int64         `json:"total_duration"`
	LessonCount   int32         `json:"lesson_count"`
	Catalogs      []*CatalogDto `json:"catalogs"`
}

type AddCatalogLessonReq struct {
	CourseID  int64   `json:"course_id"`
	CatalogID int64   `json:"catalog_id"`
	LessonIDs []int64 `json:"lesson_ids"`
}

type RemoveCatalogLessonReq struct {
	IDs []int64 `json:"ids"`
}

type UpdateCatalogLessonReq struct {
	ID          int64  `json:"id"`
	EnableTrial int32  `json:"enable_trial"`
	Name        string `json:"name"`
	ShowTime    int64  `json:"show_time"`
}
