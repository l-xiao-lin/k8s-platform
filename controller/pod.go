package controller

import (
	"k8s-platform/service"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var Pod pod

type pod struct{}

func (p *pod) GetPods(c *gin.Context) {
	params := new(ParamPodList)
	if err := c.ShouldBind(params); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		zap.L().Error("service.Pod.GetPods", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}

func (p *pod) GetPodDetail(c *gin.Context) {
	param := new(ParamPodDetail)
	if err := c.ShouldBind(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := service.Pod.GetPodDetail(param.PodName, param.Namespace)

	if err != nil {
		zap.L().Error("service.Pod.GetPodDetail", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func (p *pod) DeletePod(c *gin.Context) {
	param := new(ParamPodDelete)
	if err := c.ShouldBindJSON(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	if err := service.Pod.DeletePod(param.PodName, param.Namespace); err != nil {
		zap.L().Error("service.Pod.DeletePod failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func (p *pod) UpdatePod(c *gin.Context) {
	param := new(ParamPodUpdate)
	if err := c.ShouldBindJSON(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	if err := service.Pod.UpdatePod(param.Namespace, param.Content); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func (p *pod) GetPodContainer(c *gin.Context) {
	param := new(ParamPodContainer)
	if err := c.ShouldBind(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := service.Pod.GetPodContainer(param.PodName, param.Namespace)
	if err != nil {
		zap.L().Error("service.Pod.GetPodContainer failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func (p *pod) GetPodLog(c *gin.Context) {
	param := new(ParamPodLog)
	if err := c.ShouldBind(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := service.Pod.GetPodLog(param.ContainerName, param.PodName, param.Namespace)
	if err != nil {
		zap.L().Error("service.Pod.GetPodLog failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}

func (p *pod) GetPodNumPerNp(c *gin.Context) {
	data, err := service.Pod.GetPodNumPerNp()
	if err != nil {
		zap.L().Error("service.Pod.GetPodNumPerNp failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
