package decision

import (
	"bigagent/internal/kubernetes"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"log"
)

// AbnormalPod 定义“重启不正常”的 Pod 结构体
type AbnormalPod struct {
	k            func() *kubernetes.DefaultK8sOperator
	ctx          context.Context
	Cluster      string // 集群名称
	Namespace    string
	Pod          string
	Container    string
	RestartCount int32
	Reason       string
	Message      string
}

func NewAbnormalPod(k func() *kubernetes.DefaultK8sOperator, ctx context.Context, cluster string, namespace string) *AbnormalPod {
	return &AbnormalPod{k: k, ctx: ctx, Cluster: cluster, Namespace: namespace}
}

func (d *AbnormalPod) StartCheck() error {
	report, err := d.CheckPodAbnormalRestarts(d.ctx, d.k, d.Cluster, d.Namespace, 1)
	for _, pod := range report {
		log.Printf("重启异常的pod为%s", pod.Pod)
	}
	return err
}

// CheckPodAbnormalRestarts 使用 DefaultK8sOperator 判定 Pod 是否“重启不正常”
// - namespace 为空则扫描所有命名空间
// - restartThreshold: 重启次数阈值（如 3）
func (d *AbnormalPod) CheckPodAbnormalRestarts(ctx context.Context, op func() *kubernetes.DefaultK8sOperator, cluster, namespace string, restartThreshold int32) ([]AbnormalPod, error) {
	if op() == nil {
		return nil, fmt.Errorf("k8s 操作手没创建！")
	}
	// 通过 op.ListPods 列表 Pod（确保使用你的 op）
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
					Namespace:    p.Namespace,
					Pod:          p.Name,
					Container:    cs.Name,
					RestartCount: cs.RestartCount,
					Reason:       "HighRestartCount",
					Message:      "",
				})
			}
		}
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
