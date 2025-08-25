package polling

import (
	checknode "bigagent/internal/check/node"
	checkpod "bigagent/internal/check/pod"
	decisionpod "bigagent/internal/decision/pod"
	"bigagent/internal/kubernetes"
	prom "bigagent/internal/prometheus"
	"bigagent/internal/result"
	"bigagent/internal/utils"
	"context"
)

// KubePolling k8s轮询器
type KubePolling struct{}

func NewKubePolling() *KubePolling { return &KubePolling{} }

func (k *KubePolling) RunPolling(fn PollFunc) result.Result { return fn() }

// RunAll 集成执行 main 中的所有检查与决策
func (k *KubePolling) RunAll(
	ctx context.Context,
	clusters []kubernetes.KubeCluster,
	k8sOpFactory func() *kubernetes.DefaultK8sOperator,
	promFactory func(clusterName string) (*prom.PromClient, error),
) {
	for _, kc := range clusters {
		clusterName := kc.Name
		pc, err := promFactory(clusterName)
		if err != nil || pc == nil {
			utils.DefaultLogger.WithField("cluster", clusterName).WithError(err).Error("Prometheus 客户端创建失败")
			continue
		}

		//==================== pod ===============
		// 需要强制删除的 Pod (防止循环依赖，链式调用和下面的不同)
		if dec := decisionpod.NewAbnormalPodDecision(
			checkpod.NewAbnormalPod(ctx, k8sOpFactory, clusterName, "").CheckPodNeedDelete(),
		).JudgeNeedDelete(); dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", clusterName).WithField("检查", clusterName).Info("pod检查无异常")
		}

		//==================== node ===============
		// 按 PromQL 的通用检查 file_system_read_only
		dec := checknode.NewAbnormalNode(ctx, k8sOpFactory, clusterName, pc).CheckNodeByProm("file_system_read_only").JudgeDeadNode()
		if dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", clusterName).WithField("检查", "file_system_read_only").Info("节点检查无异常")
		}

		// 按 PromQL 的通用检查 arp_too_many
		dec = checknode.NewAbnormalNode(ctx, k8sOpFactory, clusterName, pc).CheckNodeByProm("arp_too_many").JudgeDeadNode()
		if dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", clusterName).WithField("检查", "arp_too_many").Info("节点检查无异常")
		}

		// 基于数据库记录的自愈恢复（Uncordon）
		dec = checknode.NewAbnormalNode(ctx, k8sOpFactory, clusterName, pc).CheckCoNodeToUnCordon().JudgeDeadNode()
		if dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", clusterName).Info("无需恢复的异常节点")
		}

		// 对掉线的节点进行不可调度（Cordon）
		dec = checknode.NewAbnormalNode(ctx, k8sOpFactory, clusterName, pc).CheckDownNode().JudgeDownNodeToCordon()
		if dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", clusterName).Info("无掉线的异常节点")
		}

		// 对 ntp 异常的节点进行不可调度（Cordon）
		dec = checknode.NewAbnormalNode(ctx, k8sOpFactory, clusterName, pc).CheckNtpNode().JudgeNtpNodeToCordon()
		if dec != nil {
			dec()
		} else {
			utils.DefaultLogger.WithField("cluster", clusterName).Info("无ntp异常节点")
		}
	}
}
