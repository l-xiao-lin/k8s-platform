package service

import (
	"go.uber.org/zap"
	"k8s-platform/dao/mysql"
	"k8s-platform/model"
)

var Workflow workflow

type workflow struct{}

type WorkflowCreate struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	Replicas      int32                  `json:"replicas"`
	Image         string                 `json:"image"`
	Label         map[string]string      `json:"label"`
	Cpu           string                 `json:"cpu"`
	Memory        string                 `json:"memory"`
	ContainerPort int32                  `json:"container_port"`
	HealthCheck   bool                   `json:"health_check"`
	HealthPath    string                 `json:"health_path"`
	Type          string                 `json:"type"`
	Port          int32                  `json:"port"`
	NodePort      int32                  `json:"node_port"`
	Hosts         map[string][]*HttpPath `json:"hosts"`
}

func GetServiceName(name string) string {
	return name + "-svc"
}

func GetIngressName(name string) string {
	return name + "-ing"
}

func (w *workflow) CreateWorkflow(data *WorkflowCreate) (err error) {
	//1、构造Workflow结构体，往数据库中添加数据
	var ingressName string
	if data.Type == "Ingress" {
		ingressName = GetIngressName(data.Name)
	} else {
		ingressName = ""
	}
	workflow := &model.Workflow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicas,
		Deployment: data.Name,
		Service:    GetServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
	}

	if err = mysql.Workflow.CreateWorkflow(workflow); err != nil {
		zap.L().Error("mysql.Workflow.CreateWorkflow failed", zap.Error(err))
		return
	}

	//2、添加k8s资源

	if err = createWorkflowRes(data); err != nil {
		zap.L().Error("createWorkflowRes failed", zap.Error(err))
		return
	}
	return

}

func createWorkflowRes(data *WorkflowCreate) (err error) {
	//1、创建deployment资源
	dc := &DeployCreate{
		Name:          data.Name,
		Namespace:     data.Namespace,
		Replicas:      data.Replicas,
		Image:         data.Image,
		Label:         data.Label,
		Cpu:           data.Cpu,
		Memory:        data.Memory,
		ContainerPort: data.ContainerPort,
		HealthCheck:   data.HealthCheck,
		HealthPath:    data.HealthPath,
	}
	if err = Deployment.CreateDeployment(dc); err != nil {
		zap.L().Error("Deployment.CreateDeployment failed", zap.Error(err))
		return
	}

	//2、创建service资源
	//这里必须定义一个新的变量serviceType，如果直接修改data.Type == "ClusterIP"的话，后面在创建ingress资源时还需要用到data.Type变量
	var serviceType string
	if data.Type == "Ingress" {
		serviceType = "ClusterIP"
	} else {
		serviceType = data.Type
	}

	sc := &ServiceCreate{
		Name:          GetServiceName(data.Name),
		Namespace:     data.Namespace,
		Type:          serviceType,
		ContainerPort: data.ContainerPort,
		Port:          data.Port,
		NodePort:      data.NodePort,
		Label:         data.Label,
	}
	if err = Servicev1.CreateService(sc); err != nil {
		zap.L().Error("Servicev1.CreateService failed", zap.Error(err))
		return
	}

	//3、创建ingress资源
	if data.Type == "Ingress" {
		ic := &IngressCreate{
			Name:      GetIngressName(data.Name),
			Namespace: data.Namespace,
			Label:     data.Label,
			Hosts:     data.Hosts,
		}
		if err = Ingress.CreateIngress(ic); err != nil {
			zap.L().Error("Ingress.CreateIngress failed", zap.Error(err))
			return
		}
	}
	return
}

func (w *workflow) DeleteWorkflow(id int) (err error) {
	//1、先从表中查询相关的workflow信息
	workflow, err := mysql.Workflow.GetWorkflowById(id)
	if err != nil {
		zap.L().Error("mysql.Workflow.GetWorkflowById failed", zap.Int("id", id), zap.Error(err))
		return
	}
	//2、删除k8s相关资源
	if err = deleteWorkflowRes(workflow); err != nil {
		zap.L().Error("deleteWorkflow failed", zap.Error(err))
		return
	}

	//3、删除表的数据
	if err = mysql.Workflow.DeleteWorkflow(id); err != nil {
		zap.L().Error("mysql.Workflow.DeleteWorkflow failed", zap.Error(err))
		return
	}
	return

}

func deleteWorkflowRes(data *model.Workflow) (err error) {
	//1、删除deployment资源
	if err = Deployment.DeleteDeployment(data.Deployment, data.Namespace); err != nil {
		zap.L().Error("Deployment.DeleteDeployment failed", zap.Error(err))
		return
	}
	//2、删除service资源

	if err = Servicev1.DeleteService(GetServiceName(data.Name), data.Namespace); err != nil {
		zap.L().Error("Servicev1.DeleteService failed", zap.Error(err))
		return
	}

	//3、删除ingress资源

	if data.Type == "Ingress" {
		if err = Ingress.DeleteIngress(GetIngressName(data.Name), data.Namespace); err != nil {
			zap.L().Error("Ingress.DeleteIngress failed", zap.Error(err))
			return
		}
	}
	return

}
