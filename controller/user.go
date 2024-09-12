package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"k8s-platform/dao/mysql"
	"k8s-platform/model"
	"k8s-platform/service"
)

var User user

type user struct{}

func (u *user) SignupHandler(c *gin.Context) {
	param := new(model.ParamSignup)
	if err := c.ShouldBindJSON(param); err != nil {
		zap.L().Error("SignupHandler with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}

	if err := service.User.Signup(param); err != nil {
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

func (u *user) LoginHandler(c *gin.Context) {
	param := new(model.ParamLogin)
	if err := c.ShouldBindJSON(param); err != nil {
		zap.L().Error("LoginHandler with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}

	data, err := service.User.Login(param)

	if err != nil {
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"user_id":  data.UserID,
		"username": data.Username,
		"token":    data.Token,
	})

}
