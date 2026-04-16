package do

type LarkUserInfo struct {
	Name            string `json:"name"`
	EnName          string `json:"en_name"`
	AvatarURL       string `json:"avatar_url"`
	AvatarThumb     string `json:"avatar_thumb"`
	AvatarMiddle    string `json:"avatar_middle"`
	AvatarBig       string `json:"avatar_big"`
	OpenID          string `json:"open_id"`
	UnionID         string `json:"union_id"`
	Email           string `json:"email"`
	EnterpriseEmail string `json:"enterprise_email"`
	UserID          string `json:"user_id"`
	Mobile          string `json:"mobile"`
	TenantKey       string `json:"tenant_key"`
	EmployeeNo      string `json:"employee_no"`
}

type LarkUserAccessToken struct {
	Code              int64  `json:"code"`
	AccessToken       string `json:"access_token"`
	TenantAccessToken string `json:"tenant_access_token"`
	ExpiresIn         int64  `json:"expires_in"`
	ErrCode           int64  `json:"error"`
	ErrMsg            string `json:"error_description"`
}
