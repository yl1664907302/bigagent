package check

import (
	"bigagent/internal/check/result"
	"bigagent/internal/kubernetes"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"log"
)

// AbnormalPod 定义“重启不正常”的 Pod 结构体
type AbnormalPod struct {
	k            func() *kubernetes.DefaultK8sOperator `json:"-"`
	ctx          context.Context                       `json:"-"`
	Cluster      string                                `json:"cluster"`
	Namespace    string                                `json:"namespace"`
	Pod          string                                `json:"pod"`
	Container    string                                `json:"container"`
	RestartCount int32                                 `json:"restart_count"`
	Reason       string                                `json:"reason"`
	Message      string                                `json:"message"`
}

func NewAbnormalPod(ctx context.Context, k func() *kubernetes.DefaultK8sOperator, cluster string, namespace string) *AbnormalPod {
	return &AbnormalPod{k: k, ctx: ctx, Cluster: cluster, Namespace: namespace}
}

func (d *AbnormalPod) Check() (result.Result, error) {
	var severity string
	report, err := d.getPodAbnormalRestarts(d.ctx, d.k, d.Cluster, d.Namespace, 1)
	if err != nil {
		return nil, err
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
	return result.NewResultPod(result.Base{
		Cluster: d.Cluster,
		Items:   report,
	}, nil, "PodAbnormalRestarts", severity, len(report), report), nil
}

func (d *AbnormalPod) getPodAbnormalRestarts(ctx context.Context, op func() *kubernetes.DefaultK8sOperator, cluster, namespace string, restartThreshold int32) ([]AbnormalPod, error) {
	if op() == nil {
		return nil, fmt.Errorf("k8s 操作手没创建！")
	}

	pods, err := op().ListPods(ctx, cluster, namespace, "")
	if err != nil {
		return nil, err
	}
	var out []AbnormalPod
	for _, p := range pods {
		// 跳过已完成
		if p.Status.Phase == corev1.PodSucceeded {
			continue
		}
		// 聚合 init 与 app 容器
		statuses := make([]corev1.ContainerStatus, 0, len(p.Status.InitContainerStatuses)+len(p.Status.ContainerStatuses))
		statuses = append(statuses, p.Status.InitContainerStatuses...)
		statuses = append(statuses, p.Status.ContainerStatuses...)

		for _, cs := range statuses {
			// Waiting 异常
			if cs.State.Waiting != nil && isBadWaitingReason(cs.State.Waiting.Reason) {
				out = append(out, AbnormalPod{
					Cluster:      cluster,
					Namespace:    p.Namespace,
					Pod:          p.Name,
					Container:    cs.Name,
					RestartCount: cs.RestartCount,
					Reason:       cs.State.Waiting.Reason,
					Message:      cs.State.Waiting.Message,
				})
				continue
			}
			// 非零退出 + 重启次数过高
			if term := cs.LastTerminationState.Terminated; term != nil {
				if term.ExitCode != 0 && cs.RestartCount >= restartThreshold {
					out = append(out, AbnormalPod{
						Cluster:      cluster,
						Namespace:    p.Namespace,
						Pod:          p.Name,
						Container:    cs.Name,
						RestartCount: cs.RestartCount,
						Reason:       term.Reason,
						Message:      term.Message,
					})
					continue
				}
				// OOMKilled 直接判异常
				if term.Reason == "OOMKilled" {
					out = append(out, AbnormalPod{
						Cluster:      cluster,
						Namespace:    p.Namespace,
						Pod:          p.Name,
						Container:    cs.Name,
						RestartCount: cs.RestartCount,
						Reason:       term.Reason,
						Message:      term.Message,
					})
					continue
				}
			}
			// 重启次数过高
			if cs.RestartCount >= restartThreshold {
				out = append(out, AbnormalPod{
					Cluster:      cluster,
					Namespace:    p.Namespace,
					Pod:          p.Name,
					Container:    cs.Name,
					RestartCount: cs.RestartCount,
					Reason:       "HighRestartCount",
					Message:      "Container has restarted too many times",
				})
			}
		}
	}
	for _, pod := range out {
		log.Println(pod.Reason)
	}
	return out, nil
}

func isBadWaitingReason(reason string) bool {
	switch reason {
	case "CrashLoopBackOff", "ImagePullBackOff", "ErrImagePull", "CreateContainerError":
		return true
	default:
		return false
	}
}
