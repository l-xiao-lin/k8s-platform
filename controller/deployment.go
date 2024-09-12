package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s-platform/service"
)

var Deployment deployment

type deployment struct{}

func (d *deployment) GetDeployments(c *gin.Context) {
	param := new(ParamDeploymentList)
	if err := c.ShouldBind(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := service.Deployment.GetDeployments(param.FilterName, param.Namespace, param.Limit, param.Page)
	if err != nil {
		zap.L().Error("service.Deployment.GetDeployments failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func (d *deployment) GetDeploymentDetail(c *gin.Context) {
	param := new(ParamDeploymentDetail)
	if err := c.ShouldBind(param); err != nil {
		zap.L().Error("ParamDeploymentDetail failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
	}

	data, err := service.Deployment.GetDeploymentDetail(param.DeploymentName, param.Namespace)
	if err != nil {
		zap.L().Error("service.Deployment.GetDeploymentDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//将 *appsv1.Deployment 转成字符切片
	jsonStr, _ := json.Marshal(data)
	ResponseSuccess(c, string(jsonStr))

}

func (d *deployment) ScaleDeployment(c *gin.Context) {
	param := new(ParamDeploymentScale)
	if err := c.ShouldBindJSON(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := service.Deployment.ScaleDeployment(param.DeploymentName, param.Namespace, param.ScaleNum)
	if err != nil {
		zap.L().Error("service.Deployment.ScaleDeployment failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func (d *deployment) CreateDeployment(c *gin.Context) {
	param := new(service.DeployCreate)
	if err := c.ShouldBindJSON(param); err != nil {
		zap.L().Error("CreateDeployment invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	if err := service.Deployment.CreateDeployment(param); err != nil {
		zap.L().Error("service.Deployment.CreateDeployment failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}

func (d *deployment) DeleteDeployment(c *gin.Context) {
	param := new(ParamDeploymentDel)
	if err := c.ShouldBindJSON(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	if err := service.Deployment.DeleteDeployment(param.DeploymentName, param.Namespace); err != nil {
		zap.L().Error("service.Deployment.DeleteDeployment failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func (d *deployment) RestartDeployment(c *gin.Context) {
	param := new(ParamDeploymentRestart)
	if err := c.ShouldBindJSON(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	if err := service.Deployment.RestartDeployment(param.DeploymentName, param.Namespace); err != nil {
		zap.L().Error("service.Deployment.RestartDeployment failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func (d *deployment) UpdateDeployment(c *gin.Context) {
	param := new(ParamDeploymentUpdate)
	if err := c.ShouldBindJSON(param); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	if err := service.Deployment.UpdateDeployment(param.Namespace, param.Content); err != nil {
		zap.L().Error("service.Deployment.UpdateDeployment failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}

func (d *deployment) GetDeployNumPerNp(c *gin.Context) {
	data, err := service.Deployment.GetDeployNumPerNp()
	if err != nil {
		zap.L().Error("service.Deployment.GetDeployNumPerNp failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}
