package service

import (
	"k8s-platform/dao/mysql"
	"k8s-platform/model"
	"k8s-platform/pkg/jwt"
	"k8s-platform/pkg/snowflake"
)

var User user

type user struct{}

func (u *user) Signup(p *model.ParamSignup) (err error) {
	//1、判断用户是否已存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return
	}
	//2、雪花算法生成user_id
	userID := snowflake.GenID()

	//3、构建user数据
	user := &model.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//4、保存进数据库
	return mysql.InsertUser(user)

}

func (u *user) Login(p *model.ParamLogin) (user *model.User, err error) {
	//1、查询用户是否存在并且密码是否正确
	user = &model.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err = mysql.Login(user); err != nil {
		return nil, err
	}

	//2、生成token
	token, err := jwt.GetToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}
