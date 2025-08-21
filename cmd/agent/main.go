package main

import (
	"bigagent/inits"
	checknode "bigagent/internal/check/node"
	check "bigagent/internal/check/pod"
	"bigagent/internal/config/global"
	decisionnode "bigagent/internal/decision/node"
	decision "bigagent/internal/decision/pod"
	"bigagent/internal/kubernetes"
	prom "bigagent/internal/prometheus"
	"bigagent/internal/utils"
	"context"
	"os"
)

func init() {
	// 创建一个通道来接收系统信号
	sigs := make(chan os.Signal, 1)
	inits.InitCmd(sigs)
	inits.LoggerInit()
	inits.InstallIfNotExists([]string{"git", "wget", "curl", "gcc", "make", "jq", "https://pkg.osquery.io/rpm/osquery-5.11.0-1.linux.x86_64.rpm"})
	inits.AgentRegister()
	err := inits.InitDB()
	if err != nil {
		utils.DefaultLogger.WithError(err).Error("数据库初始化失败")
	}
	//inits.InitOsqueryClient()
	inits.Crontab()
	inits.ListerChannel()
	utils.DefaultLogger.WithField("version", "20250730").Info("Agent 版本")
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

func promClientForCluster(clusterName string) (*prom.PromClient, error) {
	promMap := global.V.GetStringMapString("k8s_compute_prom_addr_map")
	addr := promMap[clusterName]
	return prom.NewPromClient(addr, global.V.GetInt("recovery_conf.query_prom_time_out_seconds"))
}

func makeK8sOperator(clusters []kubernetes.KubeCluster) func() *kubernetes.DefaultK8sOperator {
	return func() *kubernetes.DefaultK8sOperator {
		return &kubernetes.DefaultK8sOperator{Clusters: clusters}
	}
}

func main() {
	ctx := context.Background()

	clusters := kubeClustersFromConfig()
	if len(clusters) == 0 {
		utils.DefaultLogger.WithField("config_key", "k8s_configs").Warn("未配置任何 k8s 集群")
		inits.RunG()
		inits.Hander(global.V.GetString("system.addr"))
		return
	}

	k8sOp := makeK8sOperator(clusters)

	for _, kc := range clusters {
		// Prometheus 客户端
		pc, err := promClientForCluster(kc.Name)
		if err != nil || pc == nil {
			utils.DefaultLogger.WithField("cluster", kc.Name).WithError(err).Error("Prometheus 客户端创建失败")
			continue
		}

		// Node 异常检测（按 PromQL 类型）
		if dec := decision.NewAbnormalPodDecision(
			check.NewAbnormalPod(ctx, k8sOp, kc.Name, "").CheckPodNeedDelete(),
		).JudgeNeedDelete(); dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", kc.Name).WithField("检查", kc.Name).Info("pod检查无异常")
		}

		// Node 异常检测（按 PromQL 类型）
		if dec := decisionnode.NewAbnormalNodeDecision(
			checknode.NewAbnormalNode(ctx, k8sOp, kc.Name, pc).CheckNodeByProm("file_system_read_only"),
		).JudgeDeadNode(); dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", kc.Name).WithField("检查", "file_system_read_only").Info("节点检查无异常")
		}

		if dec := decisionnode.NewAbnormalNodeDecision(
			checknode.NewAbnormalNode(ctx, k8sOp, kc.Name, pc).CheckNodeByProm("arp_too_many"),
		).JudgeDeadNode(); dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", kc.Name).WithField("检查", "arp_too_many").Info("节点检查无异常")
		}

		// 基于数据库记录的自愈恢复（Uncordon）
		if dec := decisionnode.NewAbnormalNodeDecision(
			checknode.NewAbnormalNode(ctx, k8sOp, kc.Name, pc).CheckCoNodeToUnCordon(),
		).JudgeDeadNode(); dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", kc.Name).Info("无需恢复的异常节点")
		}
	}

	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
