package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s-platform/service"
)

var Ingress ingress

type ingress struct{}

func (i *ingress) CreateIngress(c *gin.Context) {
	param := new(service.IngressCreate)
	if err := c.ShouldBindJSON(param); err != nil {
		zap.L().Error("CreateIngress invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	if err := service.Ingress.CreateIngress(param); err != nil {
		zap.L().Error("service.Ingress.CreateIngress failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}

func (i *ingress) DeleteIngress(c *gin.Context) {
	param := new(ParamIngressDel)
	if err := c.ShouldBindJSON(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	if err := service.Ingress.DeleteIngress(param.IngressName, param.Namespace); err != nil {
		zap.L().Error("service.Ingress.DeleteIngress failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
