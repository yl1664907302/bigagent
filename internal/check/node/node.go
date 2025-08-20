package check

import (
	"bigagent/internal/kubernetes"
	prom "bigagent/internal/prometheus"
	"bigagent/internal/result"
	"bigagent/internal/utils"
	"context"
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

func (n *AbnormalNode) CheckNodeByProm(args ...interface{}) result.Result {
	var ql string
	if args != nil {
		ql = args[0].(string)
	}
	var results result.Results
	results.Count = 0
	results.Items = nil
	//ql1 := global.V.GetString("check_ql_map.file_system_read_only")
	//ql := global.V.GetString("check_ql_map.arp_too_many")
	// 临时 map，用于存放多个查询语句
	//qlMap := map[string]string{
	//"file_system_read_only": ql1,
	//	"arp_too_many": ql2,
	//}
	// 设置Prometheus查询超时时间为5秒
	//tw := 5
	promeResults, err := n.PromC.InstantQuery(ql)
	if err != nil {
		utils.DefaultLogger.Errorf("check_ql_map.arp_too_many:%v", err)
		result.NewResults(result.Base{Cluster: n.Cluster, Items: nil}, nil, "check_ql_map.arp_too_many", "info", 0, nil)
	}

	// vecs 代表这个集群中 有多少个 异常的ntp 节点
	num := len(promeResults)
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
			"check_ql_map.arp_too_many",
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
