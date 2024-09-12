package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s-platform/service"
	"strconv"
)

var Workflow workflow

type workflow struct{}

func (w *workflow) CreateWorkflow(c *gin.Context) {
	param := new(service.WorkflowCreate)
	if err := c.ShouldBindJSON(param); err != nil {
		zap.L().Error("CreateWorkflow invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	if err := service.Workflow.CreateWorkflow(param); err != nil {
		zap.L().Error("service.Workflow.CreateWorkflow failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)

}

func (w *workflow) DeleteWorkflow(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	if err = service.Workflow.DeleteWorkflow(id); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
