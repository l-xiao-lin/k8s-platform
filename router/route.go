package router

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"k8s-platform/controller"
	"k8s-platform/logger"
	"k8s-platform/middlewares"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(ginzap.Ginzap(logger.LG, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.LG, true))

	v1 := r.Group("/api/k8s/")

	//用户注册登录
	v1.POST("/signup", controller.User.SignupHandler)
	v1.POST("/login", controller.User.LoginHandler)

	v1.Use(middlewares.JwtAuthMiddleware())
	{
		//pod
		v1.GET("/pods", controller.Pod.GetPods)
		v1.GET("/pod/detail", controller.Pod.GetPodDetail)
		v1.DELETE("/pod/del", controller.Pod.DeletePod)
		v1.PUT("/pod/update", controller.Pod.UpdatePod)
		v1.GET("/pod/container", controller.Pod.GetPodContainer)
		v1.GET("/pod/log", controller.Pod.GetPodLog)
		v1.GET("/pod/numnp", controller.Pod.GetPodNumPerNp)

		//deployment
		v1.GET("/deployments", controller.Deployment.GetDeployments)
		v1.GET("/deployment/detail", controller.Deployment.GetDeploymentDetail)
		v1.PUT("/deployment/scale", controller.Deployment.ScaleDeployment)
		v1.POST("/deployment/create", controller.Deployment.CreateDeployment)
		v1.DELETE("/deployment/del", controller.Deployment.DeleteDeployment)
		v1.PUT("/deployment/restart", controller.Deployment.RestartDeployment)
		v1.PUT("/deployment/update", controller.Deployment.UpdateDeployment)
		v1.GET("/deployment/numnp", controller.Deployment.GetDeployNumPerNp)

		//service
		v1.POST("/service/create", controller.Servicev1.CreateService)
		v1.DELETE("/service/del", controller.Servicev1.DeleteService)

		//ingress
		v1.POST("/ingress/create", controller.Ingress.CreateIngress)
		v1.DELETE("/ingress/del", controller.Ingress.DeleteIngress)

		//workflow
		v1.POST("/workflow/create", controller.Workflow.CreateWorkflow)
		v1.DELETE("/workflow/del/:id", controller.Workflow.DeleteWorkflow)
	}

	return r
}
