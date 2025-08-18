package guard

import (
	model "bigagent/internal/model/k8s"
	"bigagent/internal/utils"
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

// wait的封装
func (gr *Guarder) RunCommonModuleManager(ctx context.Context) error {
	// (ctxAll,执行的方法，每隔多长时间)
	go wait.UntilWithContext(ctx, gr.RunCommonModule, time.Duration(gr.Cg.ModuleCommonC.CheckIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Infof("RunCommonModuleManager.exit.receive_quit_signal")
	return nil

}

// 模块真实的执行方法
func (gr *Guarder) RunCommonModule(ctx context.Context) {
	var (
		wg sync.WaitGroup
	)

	// 遍历所有集群的prometheus-api map
	for clusterName, promAddr := range gr.Cg.K8sComputePromAddrMap {
		if _, exists := gr.Cg.PluginNtpC.EnabledClusters[clusterName]; !exists {
			klog.Infof("[RunCommonModule.config.disable][cluster:%v]", clusterName)
			continue
		}
		clusterName := clusterName
		promAddr := promAddr

		wg.Add(1)
		go func() {
			defer wg.Done()
			gr.RunCommonModuleOnCluster(clusterName, promAddr)
		}()
	}
	wg.Wait()
}

func (gr *Guarder) RunCommonModuleOnCluster(clusterName, promAddr string) {
	klog.Infof("RunCommonModuleOnCluster.start[clusterName:%v]", clusterName)

	// 遍历多个ql
	for moduleName, checkQl := range gr.Cg.ModuleCommonC.CheckQlMap {
		// 先进行ntp query check_ql
		ql := checkQl
		tw := gr.Cg.ModuleCommonC.QueryPromTimeOutSeconds
		vecs, err := gr.PromInstantQuery(ql, promAddr, tw)
		if err != nil {
			klog.Errorf("RunCommonModuleOnCluster.PromInstantQuery.err[clusterName:%v][moduleName:%v][ql:%v][err:%v]",
				clusterName,
				moduleName,
				ql,
				err)
			continue
		}

		// vecs 代表这个集群中 有多少个 异常的ntp 节点
		num := len(vecs)
		for index, vec := range vecs {
			vec := vec
			labelMap := vec.Metric
			nodeName := string(labelMap["node"])
			ip := string(labelMap["intance"])
			if nodeName == "" {
				continue
			}
			value := vec.Value.String()

			klog.Infof("[RunCommonModuleOnCluster.QueryRes.print][clusterName:%v][moduleName:%v][%d/%d][nodeName:%v][value:%v]",
				clusterName,
				moduleName,
				index+1,
				num,
				nodeName,
				value,
			)
			gr.CommonModuleDealOneNode(moduleName, clusterName, nodeName, ip)
		}

	}

}

func (gr *Guarder) CommonModuleDealOneNode(moduleName, clusterName, nodeName, ip string) {
	toDayStr := time.Now().Format("2006-01-02")
	// 需要根据模块名拿到配置 好的 对的recoveryQl
	recoveryQl := gr.Cg.ModuleCommonC.RecoveryQlMap[moduleName]

	// 判断是否在db中

	nodeDb := model.NodeMaintenance{
		Name:        nodeName,
		Ip:          ip,
		ClusterName: clusterName,
		ModuleName:  moduleName,
		Reason:      moduleName,
		FirstDate:   toDayStr,
		RecoveryQl:  recoveryQl,
	}
	ok, _ := nodeDb.CheckExist()
	if ok {
		klog.Infof("CommonModuleDealOneNode.Already.Deal[cluster:%v][node:%v][moduleName:%v]", clusterName, nodeName, moduleName)
		return
	}

	// 判断每日限流
	// 首先拿到这个 模块的限流数字
	limitNum := gr.Cg.ModuleCommonC.CordonDailyLimitMap[moduleName]

	todayMNodes, err := model.GetDailyNodeMaintenances(toDayStr, moduleName, clusterName)
	if err != nil {
		klog.Errorf("CommonModuleDealOneNode.getGetDailyNodeMaintenances.err:%v", err)
		return
	}

	// 如果限流被拦截 那也不能做
	if len(todayMNodes) > limitNum {
		msg := fmt.Sprintf("[日期:%v][集群:%v]\n][模块:%v 达到每日限流停止操作:%v][节点:%v]",

			toDayStr,
			clusterName,
			moduleName,
			limitNum,
			nodeName,
		)

		klog.Infof(msg)
		utils.DingDingMsgDirectSend(gr.Cg.ModuleCommonC.ImDingDingC, msg)
		return
	}

	// 先cordon它
	err = gr.CordonNode(nodeName, clusterName, moduleName)
	cordonRes := "成功"
	if err != nil {
		cordonRes = fmt.Sprintf("失败：%v", err)
	}
	msg := fmt.Sprintf("[集群:%v]\n[%s 插件cordon节点：%v结果：%v ]\n[今日操作数:%v]",
		clusterName,
		moduleName,
		nodeName,
		cordonRes,
		len(todayMNodes)+1,
	)
	klog.Infof(msg)
	// 先通知一下
	utils.DingDingMsgDirectSend(gr.Cg.ModuleCommonC.ImDingDingC, msg)

	_, err = nodeDb.AddOrGetOne()
	if err != nil {
		klog.Errorf("CommonModuleDealOneNode.addToDb.err[moduleName:%v][node:%v][err:%v]", moduleName, nodeDb, err)

	}
	klog.Infof("CommonModuleDealOneNode.addToDb.success[moduleName:%v][node:%v][err:%v]", moduleName, nodeDb.Name, err)

}
