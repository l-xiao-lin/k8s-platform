package service

import (
	"go.uber.org/zap"
	"k8s-platform/setting"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	clientSet *kubernetes.Clientset
}

func (k *k8s) Init() (err error) {
	conf, err := clientcmd.BuildConfigFromFlags("", setting.Conf.Kubeconfig)
	if err != nil {
		zap.L().Error("创建k8s配置失败", zap.Error(err))
		return
	}
	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		zap.L().Error("创建k8s clientSet失败", zap.Error(err))
	} else {
		zap.L().Info("创建k8s clientSet成功")
	}
	k.clientSet = clientSet
	return
}
