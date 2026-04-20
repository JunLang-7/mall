package do

type PermCreate struct {
	AdminUserID int64
	Code        string
	Type        int32
	Name        string
	PagePath    string
	ParentID    int64
	Sort        int32
	Desc        string
}

type PermUpdate struct {
	ID int64
	PermCreate
}

type PermUpdateList struct {
	List []PermUpdate
}

type PermDelete struct {
	ID int64
}
