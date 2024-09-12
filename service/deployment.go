package service

import (
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"
)

var Deployment deployment

var ErrorNoFoundRes = errors.New("未查到相关的资源")

type deployment struct{}

func (d *deployment) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i])
	}
	return cells
}

func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}

type DeploymentsResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"containerPort"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
}

type DeploysNp struct {
	Namespace string `json:"namespace"`
	DeployNum int    `json:"deploy_num"`
}

// GetDeployments 获取deployment列表
func (d *deployment) GetDeployments(filterName, namespace string, limit, page int) (deploymentsResp *DeploymentsResp, err error) {
	deploymentList, err := K8s.clientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error("get deploymentList failed", zap.Error(err))
		return nil, err
	}
	selectableData := dataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	if total == 0 {
		zap.L().Error("未找到相关的deployment资源", zap.Error(err))
		return nil, ErrorNoFoundRes
	}
	data := filtered.Sort().Paginate()

	deployments := d.fromCells(data.GenericDataList)
	return &DeploymentsResp{
		Items: deployments,
		Total: total,
	}, nil
}

// GetDeploymentDetail 获取deployment详情
func (d *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = K8s.clientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error("get deployment failed", zap.Error(err))
		return nil, err
	}
	return deployment, nil
}

// ScaleDeployment 更新副本数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replica int32, err error) {
	scale, err := K8s.clientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error("GetScale failed", zap.Error(err))
		return 0, err
	}
	//修改副本数
	scale.Spec.Replicas = int32(scaleNum)

	//更新副本数
	newScale, err := K8s.clientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error("UpdateScale failed", zap.Error(err))
		return 0, err
	}

	return newScale.Spec.Replicas, nil
}

// CreateDeployment 创建deployment
func (d *deployment) CreateDeployment(data *DeployCreate) (err error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "http",
									ContainerPort: data.ContainerPort,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							LivenessProbe:  nil,
							ReadinessProbe: nil,
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

	if data.HealthCheck {
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}

		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}

		deployment.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}

		deployment.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(data.Cpu),
			corev1.ResourceMemory: resource.MustParse(data.Memory),
		}

	}
	_, err = K8s.clientSet.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error("create deployment failed", zap.Error(err))
		return err
	}
	return

}

// DeleteDeployment 删除deployment
func (d *deployment) DeleteDeployment(deploymentName, namespace string) (err error) {
	err = K8s.clientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error("delete deployment failed", zap.Error(err))
		return
	}
	return
}

// RestartDeployment 重启deployment
func (d *deployment) RestartDeployment(deploymentName, namespace string) (err error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						map[string]interface{}{
							"name": deploymentName,
							"env": []map[string]interface{}{
								map[string]interface{}{
									"name":  "RESTART_TIME",
									"value": strconv.FormatInt(time.Now().Unix(), 10),
								},
							},
						},
					},
				},
			},
		},
	}

	patchByte, err := json.Marshal(patchData)
	if err != nil {
		zap.L().Error("json Marshal patchData failed", zap.Error(err))
		return
	}

	_, err = K8s.clientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, types.StrategicMergePatchType, patchByte, metav1.PatchOptions{})
	if err != nil {
		zap.L().Error("restart deployment failed", zap.Error(err))
		return
	}
	return

}

// UpdateDeployment 更新deployment
func (d *deployment) UpdateDeployment(namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}
	if err = json.Unmarshal([]byte(content), deploy); err != nil {
		zap.L().Error("json Unmarshal failed", zap.Error(err))
		return
	}

	_, err = K8s.clientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error("update deployment failed", zap.Error(err))
		return
	}
	return
}

// GetDeployNumPerNp 获取每个namespace中的deployment数量
func (d *deployment) GetDeployNumPerNp() (deploysNps []*DeploysNp, err error) {
	namespaceList, err := K8s.clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {
		deploymentList, err := K8s.clientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		deploysNp := &DeploysNp{
			Namespace: namespace.Name,
			DeployNum: len(deploymentList.Items),
		}
		deploysNps = append(deploysNps, deploysNp)
	}
	return deploysNps, nil
}
