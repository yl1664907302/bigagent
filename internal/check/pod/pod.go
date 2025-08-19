package check

import (
	"bigagent/internal/kubernetes"
	result2 "bigagent/internal/result"
	"bigagent/internal/utils"
	"context"
	corev1 "k8s.io/api/core/v1"
)

// AbnormalPod 定义“重启不正常”的 Pod 结构体
type AbnormalPod struct {
	k            func() *kubernetes.DefaultK8sOperator `json:"-"`
	ctx          context.Context                       `json:"-"`
	Cluster      string                                `json:"cluster"`
	Namespace    string                                `json:"namespace"`
	Pod          string                                `json:"pod"`
	Terminating  bool                                  `json:"terminating,omitempty"` // 是否仍处于终止状态
	Container    string                                `json:"container"`
	RestartCount int32                                 `json:"restart_count"`
	Reason       string                                `json:"reason"`
	Message      string                                `json:"message"`
}

func NewAbnormalPod(ctx context.Context, k func() *kubernetes.DefaultK8sOperator, cluster string, namespace string) *AbnormalPod {
	return &AbnormalPod{k: k, ctx: ctx, Cluster: cluster, Namespace: namespace}
}

// Check 结果整合
func (d *AbnormalPod) Check(args ...interface{}) result2.Result {
	var severity string
	report, err := d.getPodAbnormalRestarts(d.k, d.Cluster, d.Namespace, 3)
	if err != nil {
		utils.DefaultLogger.Errorf(err.Error())
	}
	switch len(report) {

	case 0:
		severity = "info"
	case 1:
		severity = "warn"
	default:
		severity = "critical"
	}
	// 返回异常核心报告
	return result2.NewResultPod(result2.Base{
		Cluster: d.Cluster,
		Items:   report,
	}, nil, "PodAbnormalRestarts", severity, len(report), report)
}

func (d *AbnormalPod) CheckPodTerminating(args ...interface{}) result2.Result {
	if len(args) > 0 {
		if pod, ok := args[0].(*corev1.Pod); ok {
			return d.getPodTerminating(d.Cluster, pod)
		}
	}
	return result2.NewResultPod(result2.Base{Cluster: d.Cluster, Items: nil}, nil, "PodTerminating", "info", 0, nil)
}

// newAbnormalPod 创建异常 Pod 记录
func newAbnormalPod(cluster, namespace, podName string, cs corev1.ContainerStatus, reason, message string) AbnormalPod {
	return AbnormalPod{
		Cluster:      cluster,
		Namespace:    namespace,
		Pod:          podName,
		Container:    cs.Name,
		RestartCount: cs.RestartCount,
		Reason:       reason,
		Message:      message,
	}
}
