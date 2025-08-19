package check

import (
	"bigagent/internal/check/result"
	"bigagent/internal/kubernetes"
	"bigagent/internal/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func (d *AbnormalPod) Check(args ...interface{}) result.Result {
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
	return result.NewResultPod(result.Base{
		Cluster: d.Cluster,
		Items:   report,
	}, nil, "PodAbnormalRestarts", severity, len(report), report)
}

func (d *AbnormalPod) CheckPodTerminating(args ...interface{}) result.Result {
	if len(args) > 0 {
		if pod, ok := args[0].(*corev1.Pod); ok {
			return d.getPodTerminating(d.k, d.Cluster, pod)
		}
	}
	return result.NewResultPod(result.Base{Cluster: d.Cluster, Items: nil}, nil, "PodTerminating", "info", 0, nil)
}

// CheckPodTerminating 判断传入的旧 Pod 是否仍处于终止状态且未被重建，整合为 Result 返回
func (d *AbnormalPod) getPodTerminating(op func() *kubernetes.DefaultK8sOperator, cluster string, oldPod *corev1.Pod) result.Result {
	severity := "critical"
	var items interface{}
	var count int
	if op == nil || oldPod == nil {
		return result.NewResultPod(result.Base{Cluster: cluster, Items: nil}, nil, "PodTerminating", severity, 0, nil)
	}
	k8sOp := op()
	cs, _, err := k8sOp.Client(cluster)
	if err != nil {
		return result.NewResultPod(result.Base{Cluster: cluster, Items: nil}, nil, "PodTerminating", severity, 0, nil)
	}
	// 获取该pod的最新状态
	newPod, err := cs.CoreV1().Pods(oldPod.Namespace).Get(d.ctx, oldPod.Name, metav1.GetOptions{})
	if err != nil {
		utils.DefaultLogger.Errorf("获取oldPod %s/%s 失败: %v", oldPod.Namespace, oldPod.Name, err)
	}
	stillTerminating := true
	count = 0
	if newPod == nil {
		severity = "warn"
		stillTerminating = false
	} else if newPod.UID != oldPod.UID {
		severity = "warn"
		stillTerminating = false
	} else if newPod.DeletionTimestamp.IsZero() {
		severity = "warn"
		stillTerminating = false
	}
	items = AbnormalPod{
		Cluster:     cluster,
		Namespace:   oldPod.Namespace,
		Pod:         oldPod.Name,
		Terminating: stillTerminating,
	}
	return result.NewResultPod(result.Base{Cluster: cluster, Items: items}, nil, "PodTerminating", severity, count, items)
}

func (d *AbnormalPod) getPodAbnormalRestarts(
	op func() *kubernetes.DefaultK8sOperator,
	cluster, namespace string,
	restartThreshold int32,
) ([]AbnormalPod, error) {

	k8sOp := op()
	if k8sOp == nil {
		return nil, fmt.Errorf("k8s 操作手没创建！")
	}

	// 只获取 Running 状态的 Pod
	pods, err := k8sOp.ListPods(d.ctx, cluster, namespace, "status.phase=Running")
	if err != nil {
		return nil, err
	}

	var out []AbnormalPod

	for _, p := range pods {
		if p.Status.Phase == corev1.PodSucceeded {
			continue // 跳过已完成的 Pod
		}

		// 聚合 Init 容器 + 普通容器
		statuses := append(p.Status.InitContainerStatuses, p.Status.ContainerStatuses...)

		for _, cs := range statuses {
			// 1. Waiting 异常
			if cs.State.Waiting != nil && isBadWaitingReason(cs.State.Waiting.Reason) {
				out = append(out, newAbnormalPod(cluster, p.Namespace, p.Name, cs, cs.State.Waiting.Reason, cs.State.Waiting.Message))
				continue
			}

			// 2. Terminated 异常
			if term := cs.LastTerminationState.Terminated; term != nil {
				switch {
				case term.Reason == "OOMKilled":
					out = append(out, newAbnormalPod(
						cluster,
						p.Namespace,
						p.Name,
						cs,
						term.Reason,
						term.Message))
					continue
				case term.ExitCode != 0 && cs.RestartCount >= restartThreshold:
					out = append(out, newAbnormalPod(
						cluster,
						p.Namespace,
						p.Name,
						cs,
						term.Reason,
						term.Message))
					continue
				}
			}

			// 3. 高重启次数异常
			if cs.RestartCount >= restartThreshold {
				out = append(out, newAbnormalPod(
					cluster,
					p.Namespace,
					p.Name,
					cs,
					"HighRestartCount",
					"Container has restarted too many times"))
			}
		}
	}
	return out, nil
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

func isBadWaitingReason(reason string) bool {
	switch reason {
	case "CrashLoopBackOff", "ImagePullBackOff", "ErrImagePull", "CreateContainerError":
		return true
	default:
		return false
	}
}
