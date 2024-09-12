package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s-platform/service"
)

var Servicev1 servicev1

type servicev1 struct {
}

func (s *servicev1) CreateService(c *gin.Context) {
	param := new(service.ServiceCreate)
	if err := c.ShouldBindJSON(param); err != nil {
		zap.L().Error("CreateService invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	if err := service.Servicev1.CreateService(param); err != nil {
		zap.L().Error("service.Servicev1.CreateService failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}

func (s *servicev1) DeleteService(c *gin.Context) {
	param := new(ParamServiceDel)
	if err := c.ShouldBindJSON(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	if err := service.Servicev1.DeleteService(param.ServiceName, param.Namespace); err != nil {
		zap.L().Error("service.Servicev1.DeleteService failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}
