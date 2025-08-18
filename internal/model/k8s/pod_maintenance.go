package model

import (
	"fmt"
	"time"
)

type PodMaintenance struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`         // pod主机名
	Ip          string    `json:"ip"`           // podip
	NodeName    string    `json:"node_name"`    // 节点名称\r\n
	PodNs       string    `json:"pod_ns"`       // ns\r\n
	ClusterName string    `json:"cluster_name"` // 集群名称
	ModuleName  string    `json:"module_name"`  // 是哪个模块标记它的
	Reason      string    `json:"reason"`       // 具体原因
	HasRestart  int       `json:"has_restart"`  // =1 代表已经重启过了
	FirstDate   string    `json:"first_date"`   // 维护第一次的天的日期
	CreateTime  time.Time `json:"create_time" xorm:"create_time created"`
	UpdateTime  time.Time `json:"update_time" xorm:"update_time updated"`
}

// 插入
// 插入一条记录
func (obj *PodMaintenance) AddOne() (int64, error) {
	rowAffect, err := DB.InsertOne(obj)
	return rowAffect, err
}

// 检查记录是否存在
func (obj *PodMaintenance) CheckExist() (bool, error) {
	exist, err := DB.Exist(obj)
	return exist, err
}

// 基于检查存在和 add封装方法
func (obj *PodMaintenance) AddOrGetOne() (int64, error) {
	exist, _ := obj.CheckExist()
	if exist {
		return obj.Id, nil
	}

	rowAffect, err := obj.AddOne()
	return rowAffect, err
}

// 查询
// 根据部分条件查询一条记录
func (obj *PodMaintenance) GetOne() (*PodMaintenance, error) {
	has, err := DB.Get(obj)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	//mountWhere := fmt.Sprintf("cluster_id=%d", obj.ClusterId)
	//total, _ := MountGetCount(mountWhere)
	//obj.NodeNum = total

	return obj, nil
}

// 更新记录
func (obj *PodMaintenance) Update() (bool, error) {
	rowAffected, err := DB.Where(fmt.Sprintf("id=%d", obj.Id)).Update(obj)
	if err != nil {
		return false, err
	}
	if rowAffected > 0 {
		return true, nil
	}
	return false, nil
}

// 查询所有
func GetAllPodMaintenances() ([]*PodMaintenance, error) {
	var objs []*PodMaintenance
	session := DB.Where(fmt.Sprintf("id>%d", 0))
	err := session.Find(&objs)
	return objs, err
}

// 根据部分条件查询
func PodMaintenanceGetByName(name string, moduleName string) (*PodMaintenance, error) {

	obj := PodMaintenance{
		Name:       name,
		ModuleName: moduleName,
	}
	// select * from node_man where name='xxx' and module_name='xx' limit 1;
	return obj.GetOne()

}

// 外部传条件 加前端分页
func PodMaintenanceGetManyWithLimit(limit, offset int, orderByColumn string, where string, args ...interface{}) ([]PodMaintenance, error) {
	var objs []PodMaintenance
	err := DB.Where(where, args...).Limit(limit, offset).Desc(orderByColumn).Find(&objs)
	if err != nil {
		return nil, err
	}
	return objs, nil

}
