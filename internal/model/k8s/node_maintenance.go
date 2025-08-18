package model

import (
	"fmt"
	"k8s.io/klog/v2"
	"time"
)

type NodeMaintenance struct {
	Id                     int64             `json:"id"`
	Name                   string            `json:"name"`              // 节点主机名
	Ip                     string            `json:"ip"`                // 节点ip
	Sn                     string            `json:"sn"`                // sn\r\n
	ClusterName            string            `json:"cluster_name"`      // 集群名称
	ModuleName             string            `json:"module_name"`       // 是哪个模块标记它的
	Reason                 string            `json:"reason"`            // 具体原因
	HasRestart             int               `json:"has_restart"`       // =1 代表已经重启过了
	HasRepair              int               `json:"has_repair"`        // =1 代表已经发出维修工单了
	HasRecovery            int               `json:"has_recovery"`      // =1 代表已经发出维修工单了
	RepairTicketUrl        string            `json:"repair_ticket_url"` // 节点保修的工单
	MultiRecoveryModuleMap map[string]string `xorm:"-"`                 // 不需要建表 ，hold model_name --> recovery_ql

	RecoveryQl string    `json:"recovery_ql"` // 自愈的ql
	FirstDate  string    `json:"first_date"`  // 维护第一次的天的日期
	CreateTime time.Time `json:"create_time" xorm:"create_time created"`
	UpdateTime time.Time `json:"update_time" xorm:"update_time updated"`
}

// 插入
// 插入一条记录
func (obj *NodeMaintenance) AddOne() (int64, error) {
	rowAffect, err := DB.InsertOne(obj)
	return rowAffect, err
}

// 检查记录是否存在
func (obj *NodeMaintenance) CheckExist() (bool, error) {
	exist, err := DB.Exist(obj)
	return exist, err
}

// 基于检查存在和 add封装方法
func (obj *NodeMaintenance) AddOrGetOne() (int64, error) {
	exist, _ := obj.CheckExist()
	if exist {
		return obj.Id, nil
	}

	rowAffect, err := obj.AddOne()
	return rowAffect, err
}

// 查询
// 根据部分条件查询一条记录
func (obj *NodeMaintenance) GetOne() (*NodeMaintenance, error) {
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
func (obj *NodeMaintenance) Update() (bool, error) {
	rowAffected, err := DB.Where(fmt.Sprintf("id=%d", obj.Id)).Update(obj)
	if err != nil {
		return false, err
	}
	if rowAffected > 0 {
		return true, nil
	}
	return false, nil
}

// 更新记录
func (obj *NodeMaintenance) MarkRecovery() (bool, error) {
	obj.HasRecovery = 1
	return obj.Update()
}

// 更新记录
func (obj *NodeMaintenance) MarkRestart() (bool, error) {
	obj.HasRestart = 1
	return obj.Update()
}

// 查询所有
func GetAllNodeMaintenances() ([]*NodeMaintenance, error) {
	var objs []*NodeMaintenance
	session := DB.Where(fmt.Sprintf("id>%d", 0))
	err := session.Find(&objs)
	return objs, err
}

// 根据部分条件查询
func NodeMaintenanceGetByName(name string, moduleName string) (*NodeMaintenance, error) {

	obj := NodeMaintenance{
		Name:       name,
		ModuleName: moduleName,
	}
	// select * from node_man where name='xxx' and module_name='xx' limit 1;
	return obj.GetOne()

}

func GetDailyNodeMaintenances(dateStr string, mod string, cluster string) ([]*NodeMaintenance, error) {

	var objs []*NodeMaintenance
	session := DB.Where(fmt.Sprintf("first_date='%s' and module_name='%s' and cluster_name='%s'  ", dateStr, mod, cluster))
	err := session.Find(&objs)

	return objs, err

}

func GetToRecoveryNodeMaintenances(dayNum int) ([]*NodeMaintenance, error) {

	var objs []*NodeMaintenance
	// 还没有自愈过，同时时间也符合要求
	session := DB.Where(fmt.Sprintf("has_recovery=0 and first_date >  DATE_SUB(CURDATE(),INTERVAL %d DAY ) ", dayNum))
	err := session.Find(&objs)

	return objs, err

}

// 获取待 发起 power 请求的机器
func NodeMaintenanceGetRestart() ([]*NodeMaintenance, error) {

	var objs []*NodeMaintenance
	session := DB.Where(fmt.Sprintf("has_restart=0 and has_repair=0 and module_name='%s' ", "node_down"))
	err := session.Find(&objs)
	return objs, err

}

// 外部传条件 加前端分页
func NodeMaintenanceGetManyWithLimit(limit, offset int, orderByColumn string, where string, args ...interface{}) ([]NodeMaintenance, error) {
	var objs []NodeMaintenance
	err := DB.Where(where, args...).Limit(limit, offset).Desc(orderByColumn).Find(&objs)
	if err != nil {
		return nil, err
	}
	return objs, nil

}

func NodeMaintenanceTestInsert() {
	toDayStr := time.Now().Format("2006-01-02")
	obj := NodeMaintenance{
		Name:        "test-node01",
		Ip:          "1.1.1.1",
		Sn:          "xwdwdwe1",
		ClusterName: "ads-aliyun-hangzhou01-ego-train-live-01",
		ModuleName:  "ntp",
		Reason:      "time wrong",
		HasRestart:  0,
		HasRepair:   0,
		FirstDate:   toDayStr,
	}
	id, err := obj.AddOrGetOne()
	klog.Infof("NodeMaintenanceTestInsert.addone[id:%v][err:%v]", id, err)
	// 查询
	objDb, err := obj.GetOne()
	klog.Infof("NodeMaintenanceTestInsert.GetOne[objDb:%v][err:%v]", objDb, err)

	// 更新
	obj.Name = "002"
	ok, err := obj.Update()
	klog.Infof("NodeMaintenanceTestInsert.Update[ok:%v][err:%v]", ok, err)

}
