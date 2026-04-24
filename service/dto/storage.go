package dto

type GetTempSecretReq struct {
	Scene    string `json:"scene"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
	ClientIP string `json:"client_ip"`
}
type TempSecretResp struct {
	SecretID      string `json:"secret_id"`
	SecretKey     string `json:"secret_key"`
	SecurityToken string `json:"security_token"`
	ExpiredTime   int64  `json:"expired_time"`
	Bucket        string `json:"bucket"`
	Region        string `json:"region"`
	Key           string `json:"key"`
	FileURL       string `json:"file_url"`
}
