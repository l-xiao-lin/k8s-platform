package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"k8s-platform/model"
)

var secret = "cisco46589"

// CheckUserExist 判断用户是否存在
func CheckUserExist(username string) (err error) {
	var count int64
	sqlStr := `select count(user_id) from user where username=?`
	if err = db.Get(&count, sqlStr, username); err != nil {
		return
	}
	if count > 0 {
		err = ErrorUserExist
	}
	return
}

// InsertUser 添加用户
func InsertUser(user *model.User) (err error) {
	//1、对密码进行加密
	user.Password = encryptPassword(user.Password)

	//2、插入表中
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *model.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id,username,password from user where username=?`
	if err = db.Get(user, sqlStr, user.Username); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorUserNotExist
		}
		return
	}
	password := encryptPassword(oPassword)
	if user.Password != password {
		return ErrorInvalidPassword
	}
	return
}
