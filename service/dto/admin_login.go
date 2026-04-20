package dto

type GetVerifyCaptchaReq struct {
	Once string `json:"once"`
	Time int64  `json:"ts"`
	Sign string `json:"sign"` // 秘钥固定加密： md5(once+daqing2025+ts) 转小写
}
type GetVerifyCaptchaResp struct {
	Key            string `json:"key"`
	ImageBs64      string `json:"image_base64"`       // 包含“data:image/jpeg;base64
	TitleImageBs64 string `json:"title_image_base64"` // 滑块图片，包含“data:image/jpeg;base64
	TitleHeight    int    `json:"title_height"`       // 滑块图片高
	TitleWidth     int    `json:"title_width"`        // 滑块图片宽
	TitleX         int    `json:"title_x"`            // 滑块图的x坐标
	TitleY         int    `json:"title_y"`            // 滑块图的y坐标
	Expire         int64  `json:"expire"`             // 过期时间
}

type CheckCaptchaReq struct {
	Key    string `json:"key"`
	SlideX int    `json:"slide_x"`
	SlideY int    `json:"slide_y"`
}

type CheckCaptchaResp struct {
	Ticket string `json:"ticket"`
	Expire int64  `json:"expire"`
}

type MobileLoginReq struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Ticket   string `json:"ticket"`
}

type LoginResp struct {
	Token string       `json:"token"`
	User  AdminUserDto `json:"user"`
}

type LarkQrCodeLoginReq struct {
	AppCode     int32  `json:"app_code"`
	Code        string `json:"code"`
	RedirectUrl string `json:"redirect_url"`
}

type GetSmsCodeVerifyReq struct {
	Scene  string `json:"scene"` // login, register, reset_password
	Mobile string `json:"mobile"`
	Ticket string `json:"ticket"`
}

type MobileVerifyCodeLoginReq struct {
	Mobile     string `json:"mobile"`
	VerifyCode string `json:"verify_code"`
}
