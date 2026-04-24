package do

type GetTempSecret struct {
	Scene    string
	FileName string
	FileSize int64
	FileType string
	ClientIP string
}

type TempSecret struct {
	SecretID      string
	SecretKey     string
	SecurityToken string
	ExpireTime    int64
	StartTime     int64
	Bucket        string
	Region        string
	Key           string
	FileURL       string
}

type GetPreviewUrl struct {
	Keys   []string
	Expire int64
}

type DeleteFile struct {
	Keys []string
}

type AddUploadFile struct {
	Scene    string
	FileKey  string
	FileName string
	FileSize int64
	FileType string
	ClientIP string
	UserID   int64
	UserType int32
}
