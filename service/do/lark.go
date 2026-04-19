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

type LarkTenantAccessToken struct {
	Code              int64  `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int64  `json:"expire"`
}

type SendLarkMsg struct {
	AppCode int32  // APP CODE
	OpenID  string // lark open_id
	IDType  string
	Content string // fmt.Sprintf("<b>手机验证码</b>\\n\\n手机号：%s \\n验证码：%s", req.Mobile, verifyCode)
}
