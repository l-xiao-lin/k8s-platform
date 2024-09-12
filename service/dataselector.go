package service

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"
)

type dataSelector struct {
	GenericDataList []DataCell
	dataSelectQuery *DataSelectQuery
}

type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

// Len 实现排序接口
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

// Sort 实现sort排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// Filter 过滤功能
func (d *dataSelector) Filter() *dataSelector {
	if d.dataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	var filteredList []DataCell
	for _, value := range d.GenericDataList {
		matches := true
		objName := value.GetName()
		if !strings.Contains(objName, d.dataSelectQuery.FilterQuery.Name) {
			matches = false
			continue
		}
		if matches {
			filteredList = append(filteredList, value)
		}
	}
	d.GenericDataList = filteredList
	return d
}

// Paginate 分页功能
func (d *dataSelector) Paginate() *dataSelector {
	page := d.dataSelectQuery.PaginateQuery.Page
	limit := d.dataSelectQuery.PaginateQuery.Limit
	if page <= 0 || limit <= 0 {
		return d
	}
	startIndex := limit * (page - 1)
	endIndex := limit * page

	//处理最后一页小于endIndex

	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}

	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d

}

type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deploymentCell) GetName() string {
	return d.Name
}

type serviceCell corev1.Service

func (s serviceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s serviceCell) GetName() string {
	return s.Name
}
