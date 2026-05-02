package do

import (
	"github.com/JunLang-7/mall/common"
)

type CreateCourse struct {
	UserID         int64
	Name           string
	CoursePrice    int64
	ServiceTime    int32
	LearnTime      int32
	Sort           int32
	Features       []string
	UpdateStatus   int32
	CoverKey       string
	DetailCoverKey string
	Detail         string
}

type UpdateCourse struct {
	UserID         int64
	ID             int64
	Name           string
	CoursePrice    int64
	ServiceTime    int32
	LearnTime      int32
	Sort           int32
	Features       []string
	UpdateStatus   int32
	CoverKey       string
	DetailCoverKey string
	Detail         string
}

type UpdateCourseStatus struct {
	UserID int64
	ID     int64
	Status int32
}

type CourseList struct {
	common.Pager
	ID              int64
	NameKW          string
	CreateStartTime int64
	CreateEndTime   int64
	UpdateStartTime int64
	UpdateEndTime   int64
	UpdateStatus    int32
	Status          int32
}
