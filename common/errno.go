package common

type Errno struct {
	Code   int
	Msg    string
	ErrMsg string
}

func (e *Errno) Error() string {
	return e.Msg
}

func (e *Errno) WithMsg(msg string) *Errno {
	e.Msg = e.Msg + "," + msg
	return e
}

func (e *Errno) WithErr(err error) *Errno {
	var msg string
	if err != nil {
		msg = err.Error()
	}
	e.ErrMsg = e.Msg + "," + msg
	return e
}

var (
	OK            = Errno{Code: 200, Msg: "OK"}
	ServerErr     = Errno{Code: 500, Msg: "Internal Server Error"}
	ParamErr      = Errno{Code: 400, Msg: "Param Error"}
	AuthErr       = Errno{Code: 401, Msg: "Auth Error"}
	PermissionErr = Errno{Code: 403, Msg: "Permission Error"}

	DataBaseErr = Errno{Code: 10000, Msg: "DataBase Error"}
	RedisErr    = Errno{Code: 10001, Msg: "Redis Error"}

	UserNotFoundErr   = Errno{Code: 11001, Msg: "User not found"}
	InvalidCaptchaErr = Errno{Code: 110002, Msg: "Captcha Verification Error"}
)
