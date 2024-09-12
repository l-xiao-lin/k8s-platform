package mysql

import "errors"

var (
	ErrorInvalidId       = errors.New("无效的ID")
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("无效的密码")
)
