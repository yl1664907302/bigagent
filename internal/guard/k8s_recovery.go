package guard

import (
	model "bigagent/internal/model/k8s"
	"bigagent/internal/utils"
	"context"
	"fmt"
	"github.com/gammazero/workerpool"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"time"
)

// wait的封装
func (gr *Guarder) RunRecoveryCheckManager(ctx context.Context) error {
	// (ctxAll,执行的方法，每隔多长时间)
	go wait.UntilWithContext(ctx, gr.RunRecoveryCheck, time.Duration(gr.Cg.RecoveryC.CheckIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Infof("RunRecoveryCheckManager.exit.receive_quit_signal")
	return nil

}

func (gr *Guarder) RunRecoveryCheck(ctx context.Context) {
	// 数据从db 中 node维护记录

	toCheckNodes, err := model.GetToRecoveryNodeMaintenances(gr.Cg.RecoveryC.CheckDayNum)
	if err != nil {
		klog.Errorf("RunRecoveryCheckManager.GetToRecoveryNodeMaintenances.err:%v", err)
		return
	}
	if len(toCheckNodes) == 0 {
		klog.Infof("RunRecoveryCheckManager.GetToRecoveryNodeMaintenances.zero")
		return
	}
	klog.Infof("RunRecoveryCheckManager.GetToRecoveryNodeMaintenances.num:%v", len(toCheckNodes))

	// 遍历节点记录 判断它是否已经符合自愈的条件了
	// 查询ql做到，每隔模块的ql 是不是不一样，
	// 那我们应该联想到：模块处理的时候就应该把 自愈的ql写入 db中

	wp := workerpool.New(10)
	//realNodes:=make( []*models.NodeMaintenance,0)

	for _, toCheckNode := range toCheckNodes {
		toCheckNode := toCheckNode
		wp.Submit(func() {
			gr.RecoveryDealOneNode(toCheckNode)

		})
	}
	wp.StopWait()

}

func (gr *Guarder) RecoveryDealOneNode(nodeObj *model.NodeMaintenance) {
	//  先拿到集群的 api地址
	apiAddr := gr.Cg.K8sComputePromAddrMap[nodeObj.ClusterName]
	if apiAddr == "" {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("RecoveryDealOneNode.getapiAddr.K8sComputePromAddrMap.empty[clusterName:%v][nodeObj.name:%v]", nodeObj.ClusterName, nodeObj.Name)
		klog.Errorf(msg)
		return
	}

	// 查询 前拼接sql
	ql := fmt.Sprintf(nodeObj.RecoveryQl, nodeObj.Name)
	vecs, err := gr.PromInstantQuery(ql, apiAddr, gr.Cg.RecoveryC.QueryPromTimeOutSeconds)
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
		nodeObj.ClusterName,
		nodeObj.Name,
		nodeObj.ModuleName,
	)
	// 开始自愈 设置 可调度
	err = gr.UnCordonNode(nodeObj.Name, nodeObj.ClusterName, nodeObj.ModuleName)
	recoveryRes := "成功"
	if err != nil {
		recoveryRes = fmt.Sprintf("失败：%v", err)
	}
	// 发通知
	msg := fmt.Sprintf(
		"[自愈检查 %v：恢复节点调度][集群:%v][模块:%v]\n[节点:%v][时间:%v]\n"+
			"[自愈检查ql:%v]",
		recoveryRes,
		nodeObj.ClusterName,
		nodeObj.ModuleName,
		nodeObj.Name,
		nodeObj.CreateTime,
		ql,
	)

	utils.DingDingMsgDirectSend(gr.Cg.RecoveryC.ImDingDingC, msg)
	// 设置这个节点已经自愈过了
	ok, err := nodeObj.MarkRecovery()
	klog.Infof("nodeObj.MarkRecovery.res.print[ok:%v][err:%v]", ok, err)

}
