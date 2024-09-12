package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Msg  interface{} `json:"msg"`
	Code ResCode     `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, ResponseData{
		Msg:  code.Msg(),
		Code: code,
		Data: nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Msg:  msg,
		Code: code,
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Msg:  CodeSuccess.Msg(),
		Code: CodeSuccess,
		Data: data,
	})
}
