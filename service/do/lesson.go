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
	IDs []int64
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

type LessonChapter struct {
	Name          string
	BeginPosition int64
	EndPosition   int64
}

type Attachment struct {
	FileKey    string `json:"file_key"`
	OriginName string `json:"origin_name"`
}

type CreateLesson struct {
	UserID        int64
	Name          string
	Detail        string
	CategoryID    int64
	VideoKey      string
	VideoFileName string
	Attachments   []Attachment
	Duration      int32
	Chapters      []LessonChapter
}

type UpdateLesson struct {
	UserID        int64
	ID            int64
	Name          string
	Detail        string
	CategoryID    int64
	VideoKey      string
	VideoFileName string
	Attachments   []Attachment
	Duration      int32
	Chapters      []LessonChapter
}

type UpdateLessonStatus struct {
	UserID int64
	ID     int64
	Status int32
}

type MoveLesson struct {
	UserID     int64
	LessonIDs  []int64
	CategoryID int64
}

type ListLesson struct {
	common.Pager
	ID              int64
	NameKw          string
	CategoryIDs     []int64
	Status          int32
	StartCreateTime int64
	EndCreateTime   int64
	StartUpdateTime int64
	EndUpdateTime   int64
}
