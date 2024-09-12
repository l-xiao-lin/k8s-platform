package service

import (
	"context"
	"go.uber.org/zap"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Ingress ingress

type ingress struct{}

type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Label     map[string]string      `json:"label"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
}

type HttpPath struct {
	Path        string `json:"path"`
	PathType    string `json:"path_type"`
	ServiceName string `json:"service_name"`
	ServicePort int32  `json:"service_port"`
}

// CreateIngress 创建ingress
func (i *ingress) CreateIngress(data *IngressCreate) (err error) {

	var ingressRules []v1.IngressRule
	var httpIngressPaths []v1.HTTPIngressPath
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec:   v1.IngressSpec{},
		Status: v1.IngressStatus{},
	}

	for key, value := range data.Hosts {
		ir := v1.IngressRule{
			Host:             key,
			IngressRuleValue: v1.IngressRuleValue{HTTP: &v1.HTTPIngressRuleValue{Paths: nil}}}

		for _, httpPath := range value {
			pathType := v1.PathType(httpPath.PathType)
			pathTypeP := &pathType
			path := v1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: pathTypeP,
				Backend: v1.IngressBackend{
					Service: &v1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: v1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
					Resource: nil,
				},
			}

			httpIngressPaths = append(httpIngressPaths, path)
		}
		ir.IngressRuleValue.HTTP.Paths = httpIngressPaths
		ingressRules = append(ingressRules, ir)
	}

	ingress.Spec.Rules = ingressRules

	_, err = K8s.clientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error("create ingress failed", zap.Error(err))
		return
	}
	return
}

// DeleteIngress 删除ingress
func (i *ingress) DeleteIngress(ingressName, namespace string) (err error) {

	if err = K8s.clientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{}); err != nil {
		zap.L().Error("delete ingress failed", zap.Error(err))
		return
	}
	return
}
