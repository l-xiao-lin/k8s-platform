package controller

type ParamPodList struct {
	FilterName string `form:"filterName"`
	Namespace  string `form:"namespace"`
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
}

type ParamPodDetail struct {
	PodName   string `form:"pod_name"`
	Namespace string `form:"namespace"`
}

type ParamPodDelete struct {
	PodName   string `json:"pod_name"`
	Namespace string `json:"namespace"`
}

type ParamPodUpdate struct {
	Namespace string `json:"namespace"`
	Content   string `json:"content"`
}

type ParamPodContainer struct {
	PodName   string `form:"pod_name"`
	Namespace string `form:"namespace"`
}

type ParamPodLog struct {
	ContainerName string `form:"container_name"`
	PodName       string `form:"pod_name"`
	Namespace     string `form:"namespace"`
}

type ParamDeploymentList struct {
	FilterName string `form:"filterName"`
	Namespace  string `form:"namespace"`
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
}

type ParamDeploymentDetail struct {
	DeploymentName string `form:"deployment_name"`
	Namespace      string `form:"namespace"`
}

type ParamDeploymentScale struct {
	DeploymentName string `json:"deployment_name"`
	Namespace      string `json:"namespace"`
	ScaleNum       int    `json:"scale_num"`
}

type ParamDeploymentDel struct {
	DeploymentName string `json:"deployment_name"`
	Namespace      string `json:"namespace"`
}

type ParamDeploymentRestart struct {
	DeploymentName string `json:"deployment_name"`
	Namespace      string `json:"namespace"`
}

type ParamDeploymentUpdate struct {
	Namespace string `json:"namespace"`
	Content   string `json:"content"`
}

type ParamServiceDel struct {
	ServiceName string `json:"service_name"`
	Namespace   string `json:"namespace"`
}

type ParamIngressDel struct {
	IngressName string `json:"ingress_name"`
	Namespace   string `json:"namespace"`
}
