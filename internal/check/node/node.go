package check

import (
	"bigagent/internal/config/global"
	"bigagent/internal/kubernetes"
	model "bigagent/internal/model/k8s"
	prom "bigagent/internal/prometheus"
	"bigagent/internal/result"
	"bigagent/internal/utils"
	"context"
	"fmt"
	"github.com/gammazero/workerpool"
	"k8s.io/klog/v2"
)

type AbnormalNode struct {
	K         func() *kubernetes.DefaultK8sOperator `json:"-"`
	Ctx       context.Context                       `json:"-"`
	Cluster   string                                `json:"cluster"`
	Namespace string                                `json:"namespace"`
	Node      string                                `json:"node"`
	Message   string                                `json:"message"`
	PromC     *prom.PromClient
}

func NewAbnormalNode(ctx context.Context, k func() *kubernetes.DefaultK8sOperator, cluster string, p *prom.PromClient) *AbnormalNode {
	return &AbnormalNode{K: k, Ctx: ctx, Cluster: cluster, PromC: p}
}
func (n *AbnormalNode) CheckCoNodeToUnCordon(args ...interface{}) result.Result {
	// 数据从db 中 node维护记录

	toCheckNodes, err := model.GetToRecoveryNodeMaintenances(global.V.GetInt("recovery_conf.check_day_num"))
	if err != nil {
		klog.Errorf("RunRecoveryCheckManager.GetToRecoveryNodeMaintenances.err:%v", err)
		return result.NewResults(result.Base{Cluster: n.Cluster, Items: nil}, nil, "CoNodeToUnCordon", "info ", 0, err)
	}
	if len(toCheckNodes) == 0 {
		klog.Infof("RunRecoveryCheckManager.GetToRecoveryNodeMaintenances.zero")
		return result.NewResults(result.Base{Cluster: n.Cluster, Items: nil}, nil, "CoNodeToUnCordon", "info", 0, err)
	}
	klog.Infof("RunRecoveryCheckManager.GetToRecoveryNodeMaintenances.num:%v", len(toCheckNodes))

	// 遍历节点记录 判断它是否已经符合自愈的条件了
	// 查询ql做到，每隔模块的ql 是不是不一样，
	// 那我们应该联想到：模块处理的时候就应该把 自愈的ql写入 db中

	wp := workerpool.New(10)
	var toCheckNodeList []model.NodeMaintenance
	for _, toCheckNode := range toCheckNodes {
		toCheckNode := toCheckNode
		wp.Submit(func() {
			// 查询 前拼接sql
			ql := fmt.Sprintf(toCheckNode.RecoveryQl, toCheckNode.Name)
			vecs, err := n.PromC.InstantQuery(ql)
			if err != nil {
				klog.Errorf("RecoveryDealOneNode.PromInstantQuery.err[ql:%v][err:%v]", ql, err)
				return
			}
			// 结果只能有一条，因为我们已经指定了这个nodeName
			if len(vecs) != 1 {
				klog.Errorf("RecoveryDealOneNode.PromInstantQuery.vecs.err[ql:%v][vecs:%v]", ql, vecs)
				return
			}
			klog.Infof("RecoveryDealOneNode.checkOk[cluster:%v][node:%v][plugin:%v]",
				toCheckNode.ClusterName,
				toCheckNode.Name,
				toCheckNode.ModuleName,
			)
			toCheckNodeList = append(toCheckNodeList, *toCheckNode)
			// 开始自愈 设置 可调度
			//err = n.UnCordonNode(toCheckNode.Name, toCheckNode.ClusterName, toCheckNode.ModuleName)
			//recoveryRes := "成功"
			//if err != nil {
			//	recoveryRes = fmt.Sprintf("失败：%v", err)
			//}
			// 发通知
			//msg := fmt.Sprintf(
			//	"[自愈检查 %v：恢复节点调度][集群:%v][模块:%v]\n[节点:%v][时间:%v]\n"+
			//		"[自愈检查ql:%v]",
			//	recoveryRes,
			//	toCheckNode.ClusterName,
			//	toCheckNode.ModuleName,
			//	toCheckNode.Name,
			//	toCheckNode.CreateTime,
			//	ql,
			//)
			//
			//util.DingDingMsgDirectSend(gr.Cg.RecoveryC.ImDingDingC, msg)
			// 设置这个节点已经自愈过了
			ok, err := toCheckNode.MarkRecovery()
			klog.Infof("nodeObj.MarkRecovery.res.print[ok:%v][err:%v]", ok, err)
		})
	}
	wp.StopWait()
	if len(toCheckNodeList) == 0 {
		return result.NewResults(result.Base{Cluster: n.Cluster, Items: nil}, nil, "CoNodeToUnCordon", "info", 1, nil)
	}
	return result.NewResults(result.Base{Cluster: n.Cluster, Items: toCheckNodeList}, nil, "CoNodeToUnCordon", "critical", 1, nil)
}

func (n *AbnormalNode) CheckNodeByProm(args ...interface{}) result.Result {
	var checkleixing string
	if args != nil {
		checkleixing = args[0].(string)
	}
	var results result.Results
	results.Count = 1
	results.Items = nil
	results.CheckName = checkleixing
	//ql1 := global.V.GetString("check_ql_map.file_system_read_only")
	//ql := global.V.GetString("check_ql_map.arp_too_many")
	// 临时 map，用于存放多个查询语句
	//qlMap := map[string]string{
	//"file_system_read_only": ql1,
	//	"arp_too_many": ql2,
	//}
	// 设置Prometheus查询超时时间为5秒
	//tw := 5
	promeResults, err := n.PromC.InstantQuery(global.V.GetString("check_ql_map." + checkleixing))
	if err != nil {
		utils.DefaultLogger.Errorf("check_ql_map:%v", err)
		return result.NewResults(result.Base{Cluster: n.Cluster, Items: nil}, nil, checkleixing, "info", 0, nil)
	}

	// vecs 代表这个集群中 有多少个 异常的ntp 节点
	num := len(promeResults)
	if num == 0 {
		//utils.DefaultLogger.Infof("集群：%s ,未发现异常的节点", n.Cluster)
		return result.NewResults(result.Base{Cluster: n.Cluster, Items: nil}, nil, checkleixing, "info", 0, nil)
	}
	for index, vec := range promeResults {
		vec := vec

		// node名称与ip存为map
		labelMap := vec.Metric
		nodeName := string(labelMap["node"])
		ip := string(labelMap["intance"])
		results = result.Results{}
		results.Base.Node2Ip[nodeName] = ip
		if nodeName == "" {
			continue
		}
		value := vec.Value.String()

		klog.Infof("[RunCommonModuleOnCluster.QueryRes.print][clusterName:%v][moduleName:%v][%d/%d][nodeName:%v][value:%v]",
			n.Cluster,
			checkleixing,
			index+1,
			num,
			nodeName,
			value,
		)
		//gr.CommonModuleDealOneNode(moduleName, clusterName, nodeName, ip)
	}
	results.Severity = "critical"
	return &results
}
