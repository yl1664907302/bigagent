package guard

import (
	model "bigagent/internal/model/k8s"
	"bigagent/internal/utils"
	"context"
	"fmt"
	"github.com/gammazero/workerpool"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

const (
	MOD_NODE_DOWN        = "node_down"
	MOD_NODE_DOWN_REASON = "node_down_guard"
	POWER_ACTION_OFF     = "POWER_OFF"
	POWER_ACTION_ON      = "POWER_ON"
)

func (gr *Guarder) RunNodeDownCheckManager(ctx context.Context) error {
	go wait.UntilWithContext(ctx, gr.RunNodeDownCheck, time.Duration(gr.Cg.NodeDownC.CheckIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Infof("RunNodeDownGuardManager.exit.receive_quit_signal")
	return nil

}

func (gr *Guarder) RunNodeDownRestartManager(ctx context.Context) error {
	go wait.UntilWithContext(ctx, gr.RestartServerOrRepair, time.Duration(gr.Cg.NodeDownC.CheckIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Infof("RunNodeDownRestartManager.exit.receive_quit_signal")
	return nil

}

func (gr *Guarder) RunNodeDownCheck(ctx context.Context) {

	var (
		wg sync.WaitGroup
	)

	// 遍历prometheus-api地址
	for clusterName, promAddr := range gr.Cg.K8sComputePromAddrMap {
		clusterName := clusterName

		if _, exists := gr.Cg.NodeDownC.EnabledClusters[clusterName]; !exists {
			klog.Infof("[RunNodeDownCheck.config.disable][cluster:%v]", clusterName)
			continue
		}
		promAddr := promAddr
		wg.Add(1)
		go func() {
			defer wg.Done()
			gr.NodeDownCheckOneCluster(clusterName, promAddr)
		}()

	}

	wg.Wait()
}

func (gr *Guarder) NodeDownCheckOneCluster(clusterName string, promAddr string) {
	klog.Infof("NodeDownCheckOneCluster.QueryMetricInstantFloat.start[cluster:%v]", clusterName)

	nodeDownCheckResMap := map[string]int{}
	thresholdNum := len(gr.Cg.NodeDownC.CheckQls)
	// 节点宕机比较严重：多条件的check
	// 单1条件 up 的问题：组件down了，但节点没down
	// 我们这里选的规则 一般都是 节点必备组件 1个的daemonset
	//     - avg_over_time(up{job="kube-proxy"}[1d])==0
	//    - avg_over_time(kube_node_status_condition{condition="Ready",status="unknown"}[1d])==1
	//    - avg_over_time(up{job="node-exporter"}[1d])==0
	for _, ql := range gr.Cg.NodeDownC.CheckQls {
		ql := ql
		res, err := gr.PromInstantQuery(ql, promAddr, gr.Cg.NodeDownC.QueryPromTimeOutSeconds)
		if err != nil {
			klog.Errorf("NodeDownCheckOneCluster.PromInstantQuery.err[cluster:%v][err:%v]", clusterName, err)
			continue
		}
		for _, v := range res {
			v := v
			//klog.Infof("NodeDownCheckOneCluster.QueryMetricInstantFloat.print[cluster:%v][metrics:%v]", clusterName, v.Metric)

			nodeName := string(v.Metric["node"])
			if nodeName == "" {
				klog.Errorf("NodeDownCheckOneCluster.QueryMetricInstantFloat.nodeName.empty[cluster:%v][v.Metric:%v]", clusterName, v.Metric)
				continue
			}
			nodeDownCheckResMap[nodeName]++

		}
	}

	realDownNodeWithIps := map[string]string{}

	// check res
	for node, num := range nodeDownCheckResMap {

		if num < thresholdNum {
			klog.Infof("NodeDownCheckOneCluster.node.not.reach.threshold.[cluster:%v][node:%v][thresholdNum:%v][num:%v]", clusterName, node, thresholdNum, num)
			continue
		}
		// 到这里说明 配置3条check_ql 都已触发 说明这个节点是真的down了
		// 但是后面比如 要做重启等操作 需要ip
		nodeNameToIpQl := fmt.Sprintf(gr.Cg.NodeDownC.NodeNameToIpQl, node)
		res, err := gr.PromInstantQuery(nodeNameToIpQl, promAddr, gr.Cg.NodeDownC.QueryPromTimeOutSeconds)
		if err != nil {
			klog.Errorf("NodeDownCheckOneCluster.gr.Cg.NodeDownC.NodeNameToIpQl.err[cluster:%v][nodeNameToIpQl:%v][err:%v]", clusterName, nodeNameToIpQl, err)
			continue
		}
		ipStr := ""
		for _, v := range res {

			v := v
			ipStr = string(v.Metric["instance"])
		}
		//if ipStr == "" {
		//	klog.Errorf("NodeDownCheckOneCluster.QueryMetricInstantFloat.ipStr.empty[cluster:%v][node:%v]", clusterName, nodeNameToIpQl)
		//	continue
		//}

		realDownNodeWithIps[node] = ipStr
	}

	realDownNodesMsg := fmt.Sprintf("[%s]\n[集群:%v][宕机总数:%v]\n",
		gr.Cg.NodeDownC.ImDingDingC.Title,

		clusterName, len(realDownNodeWithIps))

	newFound := 0
	for nodeName, nodeIp := range realDownNodeWithIps {

		toDayStr := time.Now().Format("2006-01-02")
		// cordon node

		err := gr.CordonNode(nodeName, clusterName, MOD_NODE_DOWN)
		if err != nil {
			klog.Errorf("CordonNode.err[cluster:%v][node:%v][err:%v]", clusterName, nodeName, err)
			//continue
		}

		//  查询节点ip和sn号

		// 更新到维修db中
		nodeInDb := &model.NodeMaintenance{
			Name:        nodeName,
			ClusterName: clusterName,
			ModuleName:  MOD_NODE_DOWN,
			Ip:          nodeIp,
		}
		exist, _ := nodeInDb.CheckExist()
		if exist {
			continue
		}
		newFound++
		realDownNodesMsg += fmt.Sprintf("[node:%v][ip:%v][sn:]\n", nodeName, nodeIp)

		nodeInDb.Reason = MOD_NODE_DOWN_REASON
		nodeInDb.FirstDate = toDayStr
		_, err = nodeInDb.AddOne()

		if err != nil {
			klog.Errorf("NodeDownCheckOneCluster.abnormal.cordonNode.AddOrGetOne.err[cluster:%v][node:%v][err:%v]", clusterName, nodeName, err)

		}
		klog.Infof("NodeDownCheckOneCluster.abnormal.cordonNode.AddOrGetOne.success[cluster:%v][node:%v]", clusterName, nodeName)

	}
	if newFound > 0 {
		klog.Infof("NodeDownCheckOneCluster.get.realDownNodes.print[realDownNodesMsg:%v][cluster:%v][num:%v][detail:%v]", realDownNodesMsg, clusterName, len(realDownNodeWithIps), realDownNodeWithIps)

		utils.DingDingMsgDirectSend(gr.Cg.NodeDownC.ImDingDingC, realDownNodesMsg)

	}

}

func (gr *Guarder) RestartServerOrRepair(ctx context.Context) {

	// 思考 我的数据来自哪里：宕机模块标记的 坏节点
	nodes, err := model.NodeMaintenanceGetRestart()
	if err != nil {
		klog.Errorf("models.NodeMaintenanceGetRestart.err:%v", err)
		return
	}
	if len(nodes) == 0 {
		klog.Warning("models.NodeMaintenanceGetRestart.zero")
		return
	}

	klog.Infof("models.NodeMaintenanceGetRestart.num:%v", len(nodes))
	wp := workerpool.New(20)
	for _, node := range nodes {
		node := node
		//if node.Sn == "" {
		//	continue
		//}
		wp.Submit(func() {
			gr.TocTryRestartServer(node)
		})
	}
	wp.StopWait()
}

func (gr *Guarder) TocTryRestartServer(node *model.NodeMaintenance) {
	// 问题在于 一般机器 可能夯死了
	// ssh reboot -f 不好使
	// 你应该去对接你公司内部的 强制加电等接口重启物理机 or 公有云重启 虚拟机
	restartMsg := fmt.Sprintf(
		"[节点重启和维修模块 强制重启中][节点:%v][日期:%v][集群:%v]",
		node.Name,
		node.FirstDate,
		node.ClusterName,
	)

	utils.DingDingMsgDirectSend(gr.Cg.NodeDownC.ImDingDingC, restartMsg)
	// 这里写真实的动作

	// 将这个节点标记为已重启过
	node.MarkRestart()

}
