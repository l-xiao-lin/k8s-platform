package service

import (
	"context"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var Servicev1 servicev1

type servicev1 struct{}

type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Type          string            `json:"type"`
	ContainerPort int32             `json:"container_port"`
	Port          int32             `json:"port"`
	NodePort      int32             `json:"node_port"`
	Label         map[string]string `json:"label"`
}

// CreateService 创建service
func (s servicev1) CreateService(data *ServiceCreate) (err error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Protocol: corev1.ProtocolTCP,
					Port:     data.Port,
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Label,
		},
		Status: corev1.ServiceStatus{},
	}
	if data.NodePort != 0 && data.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = data.NodePort
	}

	_, err = K8s.clientSet.CoreV1().Services(data.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error("create service failed", zap.Error(err))
		return
	}
	return
}

func (s *servicev1) DeleteService(serviceName, namespace string) (err error) {
	err = K8s.clientSet.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error("delete service failed", zap.Error(err))
		return
	}
	return

}
