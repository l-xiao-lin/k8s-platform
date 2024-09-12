package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"k8s-platform/logger"
	"k8s-platform/setting"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct{}

func (p *pod) toCells(std []corev1.Pod) []DataCell {

	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//先断言 确认是podCell类型 再转换
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}

type PodsNp struct {
	Namespace string `json:"namespace"`
	PodNum    int    `json:"podNum"`
}

// GetPods 获取Pod列表，并实现过滤、排序、分页
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	podList, err := K8s.clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.LG.Error("GetPods failed", zap.Error(err))
		return nil, err
	}

	selectableData := dataSelector{
		GenericDataList: p.toCells(podList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)

	//排序 分页
	data := filtered.Sort().Paginate()

	//将[]DataCell 类型的pod列表转为corev1.Pod列表
	pods := p.fromCells(data.GenericDataList)

	return &PodsResp{
		Items: pods,
		Total: total,
	}, nil
}

// GetPodDetail 获取pod详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error("GetPodDetail failed", zap.Error(err))
		return nil, err
	}
	return pod, nil
}

// DeletePod 删除Pod
func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.clientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error("DeletePod failed", zap.Error(err))
		return err
	}
	return
}

// UpdatePod 更新Pod
func (p *pod) UpdatePod(namespace, content string) (err error) {
	var newPod corev1.Pod
	if err = json.Unmarshal([]byte(content), &newPod); err != nil {
		zap.L().Error("json Unmarshal failed", zap.Error(err))
		return
	}
	_, err = K8s.clientSet.CoreV1().Pods(namespace).Update(context.TODO(), &newPod, metav1.UpdateOptions{})
	if err != nil {
		zap.L().Error("update pod failed", zap.Error(err))
		return err
	}
	return
}

// GetPodContainer 获取Pod中的容器名
func (p *pod) GetPodContainer(podName, namespace string) (containers []string, err error) {
	pod, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		zap.L().Error("GetPodDetail failed", zap.Error(err))
		return
	}

	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

// GetPodLog 获取pod日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {

	lineLimit := int64(setting.Conf.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}

	req := K8s.clientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		zap.L().Error("获取pod日志失败", zap.Error(err))
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		zap.L().Error("复制podLog失败", zap.Error(err))
		return "", err
	}
	return buf.String(), nil
}

// GetPodNumPerNp 获取每个namespace中的pod数量
func (p *pod) GetPodNumPerNp() (podsNps []*PodsNp, err error) {
	namespaceList, err := K8s.clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error("get namespace failed", zap.Error(err))
		return nil, err
	}
	for _, namespace := range namespaceList.Items {

		podList, err := K8s.clientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			zap.L().Error("get podList failed", zap.Error(err))
			return nil, err
		}
		podsNp := &PodsNp{
			Namespace: namespace.Name,
			PodNum:    len(podList.Items),
		}
		podsNps = append(podsNps, podsNp)
	}
	return podsNps, nil
}
