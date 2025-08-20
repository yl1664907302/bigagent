package check

import (
	"bigagent/internal/kubernetes"
	"bigagent/internal/result"
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

func (d *AbnormalPod) CheckPodNeedDelete(args ...interface{}) result.Result {
	if len(args) > 0 {
		d.Cluster = args[0].(string)
		pods, err := d.FetchPods("status.phase!=Running")
		if err != nil {
			utils.DefaultLogger.Infof("集群：%s ,未发现非Running的pod", d.Cluster)
			return result.NewResultPod(result.Base{Cluster: d.Cluster, Items: nil}, nil, "PodNeedDelete", "info", 1, nil)
		}
		return result.NewResultPod(result.Base{Cluster: d.Cluster, Items: pods}, nil, "PodNeedDelete", "critical", 1, nil)
	}
	return result.NewResultPod(result.Base{Cluster: d.Cluster, Items: nil}, nil, "PodNeedDelete", "info", 0, nil)
}

func (d *AbnormalPod) CheckPodTerminating(args ...interface{}) result.Result {
	if len(args) > 0 {
		if pod, ok := args[0].(*corev1.Pod); ok {
			if d.GetPodTerminating(pod) {
				return result.NewResultPod(result.Base{
					Cluster: d.Cluster,
					Items:   []AbnormalPod{newAbnormalPod(d.Cluster, pod.Namespace, pod.Name, corev1.ContainerStatus{}, "Terminating", "Pod is terminating")},
				}, nil, "PodTerminating", "critical", 0, nil)
			}
		}
	}
	return result.NewResultPod(result.Base{Cluster: d.Cluster, Items: nil}, nil, "PodTerminating", "info", 0, nil)
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
