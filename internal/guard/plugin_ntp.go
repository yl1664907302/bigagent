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

const (
	PLUGIN_NTP        = "ntp_time"
	PLUGIN_NTP_REASON = "ntp_time"
	MOD_VOLUME_REASON = "vm_diff_csi"
)

func TickerModel(ctxAll context.Context) error {
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-t.C:
			klog.Infof("TickerModel.ticker.run....")
			// 这里写你真实的工作
		case <-ctxAll.Done():
			klog.Infof("TickerModel.ticker.receive.quit.signal....")
			return nil
		}
	}
}

// wait的封装
func (gr *Guarder) RunNtpCheckManager(ctx context.Context) error {
	// (ctxAll,执行的方法，每隔多长时间)
	go wait.UntilWithContext(ctx, gr.RunNtpCheck, time.Duration(gr.Cg.PluginNtpC.CheckIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Infof("RunNtpCheckManager.exit.receive_quit_signal")
	return nil

}

// 模块真实的执行方法
func (gr *Guarder) RunNtpCheck(ctx context.Context) {
	var (
		wg sync.WaitGroup
	)

	// 遍历所有集群的prometheus-api map
	for clusterName, promAddr := range gr.Cg.K8sComputePromAddrMap {
		if _, exists := gr.Cg.PluginNtpC.EnabledClusters[clusterName]; !exists {
			klog.Infof("[RunNtpCheck.config.disable][cluster:%v]", clusterName)
			continue
		}
		clusterName := clusterName
		promAddr := promAddr

		wg.Add(1)
		go func() {
			defer wg.Done()
			gr.RunNtpCheckOnCluster(clusterName, promAddr)
		}()
	}
	wg.Wait()
}

func (gr *Guarder) RunNtpCheckOnCluster(clusterName, promAddr string) {
	klog.Infof("RunNtpCheckOnCluster.start[clusterName:%v]", clusterName)

	// 先进行ntp query check_ql
	ql := gr.Cg.PluginNtpC.CheckQl
	tw := gr.Cg.PluginNtpC.QueryPromTimeOutSeconds
	vecs, err := gr.PromInstantQuery(ql, promAddr, tw)
	if err != nil {
		klog.Errorf("RunNtpCheckOnCluster.PromInstantQuery.err[clusterName:%v][ql:%v][err:%v]", clusterName, clusterName, err)
		return
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

		klog.Infof("[RunNtpCheckOnCluster.QueryRes.print][clusterName:%v][%d/%d][nodeName:%v][value:%v]",
			clusterName,
			index+1,
			num,
			nodeName,
			value,
		)
		gr.PluginNtpDealOneNode(clusterName, nodeName, ip)
	}

}

func (gr *Guarder) PluginNtpDealOneNode(clusterName, nodeName, ip string) {
	toDayStr := time.Now().Format("2006-01-02")
	// 判断是否在db中
	nodeDb := model.NodeMaintenance{
		Name:        nodeName,
		Ip:          ip,
		ClusterName: clusterName,
		ModuleName:  PLUGIN_NTP,
		Reason:      PLUGIN_NTP,
		FirstDate:   toDayStr,
		RecoveryQl:  gr.Cg.PluginNtpC.RecoveryQl,
	}
	ok, _ := nodeDb.CheckExist()
	if ok {
		klog.Infof("PluginNtpDealOneNode.Already.Deal[cluster:%v][node:%v]", clusterName, nodeName)
		return
	}

	// 判断每日限流

	todayMNodes, err := model.GetDailyNodeMaintenances(toDayStr, PLUGIN_NTP, clusterName)
	if err != nil {
		klog.Errorf("PluginNtpDealOneNode.getGetDailyNodeMaintenances.err:%v", err)
		return
	}

	// 如果限流被拦截 那也不能做
	if len(todayMNodes) > gr.Cg.PluginNtpC.CordonDailyLimit {
		msg := fmt.Sprintf("[日期:%v][集群:%v]\n][模块:%v 达到每日限流停止操作:%v][节点:%v]",

			toDayStr,
			clusterName,
			PLUGIN_NTP,
			gr.Cg.PluginNtpC.CordonDailyLimit,
			nodeName,
		)

		klog.Infof(msg)
		utils.DingDingMsgDirectSend(gr.Cg.PluginNtpC.ImDingDingC, msg)
		return
	}

	// 先cordon它
	err = gr.CordonNode(nodeName, clusterName, PLUGIN_NTP)
	cordonRes := "成功"
	if err != nil {
		cordonRes = fmt.Sprintf("失败：%v", err)
	}
	msg := fmt.Sprintf("[集群:%v]\n[ntp插件cordon节点：%v结果：%v ]\n[今日操作数:%v]",
		clusterName,
		nodeName,
		cordonRes,
		len(todayMNodes)+1,
	)
	klog.Infof(msg)
	// 先通知一下
	utils.DingDingMsgDirectSend(gr.Cg.PluginNtpC.ImDingDingC, msg)

	_, err = nodeDb.AddOrGetOne()
	if err != nil {
		klog.Errorf("PluginNtpDealOneNode.addToDb.err[node:%v][err:%v]", nodeDb, err)

	}
	klog.Infof("PluginNtpDealOneNode.addToDb.success[node:%v][err:%v]", nodeDb.Name, err)

}
