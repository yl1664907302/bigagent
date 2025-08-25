package check

import (
	"bigagent/internal/config/global"
	decision "bigagent/internal/decision/node"
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

const (
	PLUGIN_NTP        = "ntp_time"
	PLUGIN_NTP_REASON = "ntp_time"
	MOD_VOLUME_REASON = "vm_diff_csi"
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

func (n *AbnormalNode) CheckNtpNode(args ...interface{}) *decision.AbnormalNodeDecision {
	// 先进行ntp query check_ql
	ql := global.V.GetString("plugin_ntp.check_ql")
	vecs, err := n.PromC.InstantQuery(ql)
	if err != nil {
		klog.Errorf("RunNtpCheckOnCluster.PromInstantQuery.err[clusterName:%v][ql:%v][err:%v]", n.Cluster, n.Cluster, err)
		return decision.NewAbnormalNodeDecision(
			&result.Results{
				Base:      result.Base{Cluster: n.Cluster, Items: nil},
				Count:     0,
				Severity:  "info",
				CheckName: "CheckNtpNode",
			})
	}

	// vecs 代表这个集群中 有多少个 异常的ntp 节点
	//num := len(vecs)
	var node2Ip map[string]string
	for _, vec := range vecs {
		vec := vec
		labelMap := vec.Metric
		nodeName := string(labelMap["node"])
		ip := string(labelMap["intance"])
		if nodeName == "" {
			continue
		}
		node2Ip[nodeName] = ip
	}
	if len(node2Ip) == 0 {
		return decision.NewAbnormalNodeDecision(
			&result.Results{
				Base:      result.Base{Cluster: n.Cluster, Items: nil},
				Count:     0,
				Severity:  "info",
				CheckName: "CheckNtpNode",
			})
	}
	return decision.NewAbnormalNodeDecision(
		&result.Results{
			Base:      result.Base{Cluster: n.Cluster, Items: nil},
			Count:     1,
			Severity:  "critical",
			CheckName: "CheckNtpNode",
		})
}

func (n *AbnormalNode) CheckDownNode(args ...interface{}) *decision.AbnormalNodeDecision {
	nodeDownCheckResMap := map[string]int{}
	thresholdNum := len(global.V.GetStringSlice("node_down.check_qls"))
	// 节点宕机比较严重：多条件的check
	// 单1条件 up 的问题：组件down了，但节点没down
	// 我们这里选的规则 一般都是 节点必备组件 1个的daemonset
	//     - avg_over_time(up{job="kube-proxy"}[1d])==0
	//    - avg_over_time(kube_node_status_condition{condition="Ready",status="unknown"}[1d])==1
	//    - avg_over_time(up{job="node-exporter"}[1d])==0
	for _, ql := range global.V.GetStringSlice("node_down.check_qls") {
		ql := ql
		res, err := n.PromC.InstantQuery(ql)
		if err != nil {
			klog.Errorf("NodeDownCheckOneCluster.PromInstantQuery.err[cluster:%v][err:%v]", n.Cluster, err)
			continue
		}
		for _, v := range res {
			v := v
			//klog.Infof("NodeDownCheckOneCluster.QueryMetricInstantFloat.print[cluster:%v][metrics:%v]", clusterName, v.Metric)

			nodeName := string(v.Metric["node"])
			if nodeName == "" {
				klog.Errorf("NodeDownCheckOneCluster.QueryMetricInstantFloat.nodeName.empty[cluster:%v][v.Metric:%v]", n.Cluster, v.Metric)
				continue
			}
			nodeDownCheckResMap[nodeName]++

		}
	}

	realDownNodeWithIps := map[string]string{}

	// check res
	for node, num := range nodeDownCheckResMap {

		if num < thresholdNum {
			klog.Infof("NodeDownCheckOneCluster.node.not.reach.threshold.[cluster:%v][node:%v][thresholdNum:%v][num:%v]", n.Cluster, node, thresholdNum, num)
			continue
		}
		// 到这里说明 配置3条check_ql 都已触发 说明这个节点是真的down了
		// 但是后面比如 要做重启等操作 需要ip
		nodeNameToIpQl := fmt.Sprintf(global.V.GetString("node_down.node_name_to_ip_ql"), node)
		res, err := n.PromC.InstantQuery(nodeNameToIpQl)
		if err != nil {
			klog.Errorf("NodeDownCheckOneCluster.gr.Cg.NodeDownC.NodeNameToIpQl.err[cluster:%v][nodeNameToIpQl:%v][err:%v]", n.Cluster, nodeNameToIpQl, err)
			continue
		}
		ipStr := ""
		for _, v := range res {

			v := v
			ipStr = string(v.Metric["instance"])
		}
		realDownNodeWithIps[node] = ipStr
	}
	if len(realDownNodeWithIps) == 0 {
		return decision.NewAbnormalNodeDecision(
			&result.Results{
				Base:      result.Base{Cluster: n.Cluster, Items: nil},
				Count:     0,
				Severity:  "info",
				CheckName: "CheckDownNode",
			})
	}
	return decision.NewAbnormalNodeDecision(
		&result.Results{
			Base:      result.Base{Cluster: n.Cluster, Items: nil},
			Count:     1,
			Severity:  "critical",
			CheckName: "CheckDownNode",
		})
}

func (n *AbnormalNode) CheckCoNodeToUnCordon(args ...interface{}) *decision.AbnormalNodeDecision {
	// 数据从db 中 node维护记录

	toCheckNodes, err := model.GetToRecoveryNodeMaintenances(global.V.GetInt("recovery_conf.check_day_num"))
	if err != nil {
		klog.Errorf("RunRecoveryCheckManager.GetToRecoveryNodeMaintenances.err:%v", err)
		return decision.NewAbnormalNodeDecision(
			&result.Results{
				Base:      result.Base{Cluster: n.Cluster, Items: nil},
				Count:     0,
				Severity:  "info",
				CheckName: "CoNodeToUnCordon",
			})
	}
	if len(toCheckNodes) == 0 {
		return decision.NewAbnormalNodeDecision(
			&result.Results{
				Base:      result.Base{Cluster: n.Cluster, Items: nil},
				Count:     0,
				Severity:  "info",
				CheckName: "CoNodeToUnCordon",
			})
	}

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
		return decision.NewAbnormalNodeDecision(
			&result.Results{
				Base:      result.Base{Cluster: n.Cluster, Items: nil},
				Count:     0,
				Severity:  "info",
				CheckName: "CoNodeToUnCordon",
			})
	}
	return decision.NewAbnormalNodeDecision(
		&result.Results{
			Base:      result.Base{Cluster: n.Cluster, Items: nil},
			Count:     1,
			Severity:  "critical",
			CheckName: "CoNodeToUnCordon",
		})
}

func (n *AbnormalNode) CheckNodeByProm(args ...interface{}) *decision.AbnormalNodeDecision {
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
		return decision.NewAbnormalNodeDecision(
			&result.Results{
				Base:      result.Base{Cluster: n.Cluster, Items: nil},
				Count:     0,
				Severity:  "info",
				CheckName: checkleixing,
			})
	}

	// vecs 代表这个集群中 有多少个 异常的ntp 节点
	num := len(promeResults)
	if num == 0 {
		utils.DefaultLogger.WithField("cluster", n.Cluster).WithField("Prometheus查询结果数量", "promResultNum").Info(global.V.GetString("check_ql_map."+checkleixing) + ",z结果为空")
		return decision.NewAbnormalNodeDecision(
			&result.Results{
				Base:      result.Base{Cluster: n.Cluster, Items: nil},
				Count:     0,
				Severity:  "info",
				CheckName: checkleixing,
			})
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
	return decision.NewAbnormalNodeDecision(
		&result.Results{
			Base:      result.Base{Cluster: n.Cluster, Items: nil},
			Count:     1,
			Severity:  "critical",
			CheckName: checkleixing,
		})
}
