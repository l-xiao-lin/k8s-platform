package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
)

var codeMsgCode = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "无效的密码",
	CodeServerBusy:      "服务器繁忙",
	CodeNeedLogin:       "请登录",
	CodeInvalidToken:    "无效的token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgCode[c]
	if !ok {
		msg = codeMsgCode[CodeServerBusy]
	}
	return msg
}
