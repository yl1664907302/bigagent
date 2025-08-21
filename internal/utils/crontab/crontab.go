package crontab

import (
	"bigagent/internal/config/global"
	"bigagent/internal/scrape/machine"
	"bigagent/internal/strategy"
	"bigagent/internal/utils"

	"github.com/robfig/cron/v3"

	// added for kube polling integration
	"bigagent/internal/kubernetes"
	"bigagent/internal/polling"
	prom "bigagent/internal/prometheus"
	"context"
)

// CronTask Crontab执行的任务列表
func cronTask() {
	//开始采集
	machine.SmpMa = machine.NewSmpMachine()
	//更新通知
	machine.NotifySmpMachineAddressChange()
}

// ScrapeCrontab 初始化采集crontab任务
func ScrapeCrontab() {
	if strategy.Agents == nil {
		utils.DefaultLogger.Warn("agent策略为空，暂停采集任务")
		return
	}
	machine.SmpMa = machine.NewSmpMachine()
	var crontabRule string
	if global.V.GetString("collection_frequency") == "" {
		crontabRule = "@every 3s"
	}
	crontabRule = "@every " + global.V.GetString("collection_frequency")
	c := cron.New()
	c.Start()

	addFunc, err := c.AddFunc(crontabRule, cronTask)
	if err != nil {
		utils.DefaultLogger.Error("定时任务启动异常：", err)
		return
	}
	utils.DefaultLogger.Info("定时任务启动成功,EntryID：", addFunc)
}

// KubePollingCrontab 按规则执行传入任务（如 KubePolling.RunAll）
func KubePollingCrontab(rule string, task func()) {
	crontabRule := rule
	if crontabRule == "" {
		freq := global.V.GetString("collection_frequency")
		if freq == "" {
			crontabRule = "@every 3s"
		} else {
			crontabRule = "@every " + freq
		}
	}
	c := cron.New()
	c.Start()
	id, err := c.AddFunc(crontabRule, task)
	if err != nil {
		utils.DefaultLogger.Error("KubePolling 定时任务启动异常：", err)
		return
	}
	utils.DefaultLogger.Info("KubePolling 定时任务启动成功,EntryID：", id)
}

// StartKubePolling 集中执行：发现集群、构造操作器与轮询器，并按配置定时运行
func StartKubePolling(ctx context.Context) {
	clusters := kubeClustersFromConfig()
	if len(clusters) == 0 {
		utils.DefaultLogger.WithField("config_key", "k8s_configs").Warn("未配置任何 k8s 集群")
		return
	}
	k8sOp := makeK8sOperator(clusters)
	poller := polling.NewKubePolling()
	KubePollingCrontab("", func() {
		poller.RunAll(ctx, clusters, k8sOp, promClientForCluster)
	})
}

func kubeClustersFromConfig() []kubernetes.KubeCluster {
	var clusters []kubernetes.KubeCluster
	kubeClientMap := global.V.GetStringMapString("k8s_configs")
	if len(kubeClientMap) == 0 {
		return clusters
	}
	for name, cfg := range kubeClientMap {
		if cfg == "" {
			continue
		}
		clusters = append(clusters, kubernetes.KubeCluster{
			Name:       name,
			Kubeconfig: cfg,
		})
	}
	return clusters
}

func makeK8sOperator(clusters []kubernetes.KubeCluster) func() *kubernetes.DefaultK8sOperator {
	return func() *kubernetes.DefaultK8sOperator {
		return &kubernetes.DefaultK8sOperator{Clusters: clusters}
	}
}

func promClientForCluster(clusterName string) (*prom.PromClient, error) {
	promMap := global.V.GetStringMapString("k8s_compute_prom_addr_map")
	addr := promMap[clusterName]
	return prom.NewPromClient(addr, global.V.GetInt("recovery_conf.query_prom_time_out_seconds"))
}
