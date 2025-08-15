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
	report, err := d.getDeploymentPodAbnormal(d.ctx, d.k, d.Cluster, d.Namespace, "")
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

func (d *AbnormalPod) getPodAbnormalRestarts(
	ctx context.Context,
	op func() *kubernetes.DefaultK8sOperator,
	cluster, namespace string,
	restartThreshold int32,
) ([]AbnormalPod, error) {

	k8sOp := op()
	if k8sOp == nil {
		return nil, fmt.Errorf("k8s 操作手没创建！")
	}

	// 只获取 Running 状态的 Pod
	pods, err := k8sOp.ListPods(ctx, cluster, namespace, "status.phase=Running")
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

// 内存缓存：记录每个容器连续异常次数
var podFailedCounter = make(map[string]int)

// 连续异常阈值
const consecutiveThreshold = 3

func (d *AbnormalPod) getDeploymentPodAbnormal(
	ctx context.Context,
	op func() *kubernetes.DefaultK8sOperator,
	cluster, namespace, deploymentName string,
) ([]AbnormalPod, error) {

	k8sOp := op()
	if k8sOp == nil {
		return nil, fmt.Errorf("k8s 操作手没创建！")
	}

	// 1. 获取 Deployment 的 Pods
	dm, err := k8sOp.GetDeployment(ctx, cluster, namespace, deploymentName)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("获取 Deployment 失败: %v", err)
	}

	podList, err := k8sOp.ListPodsByDeployment(ctx, cluster, namespace, dm.Name)
	if err != nil {
		return nil, err
	}

	var out []AbnormalPod

	for _, p := range podList.Items {
		if p.Status.Phase != corev1.PodRunning {
			continue // 只关注 Running Pod
		}

		statuses := append(p.Status.InitContainerStatuses, p.Status.ContainerStatuses...)

		for _, cs := range statuses {
			key := fmt.Sprintf("%s/%s/%s", p.Namespace, p.Name, cs.Name)
			isAbnormal := false

			// 判断异常条件
			if cs.State.Waiting != nil && isBadWaitingReason(cs.State.Waiting.Reason) {
				isAbnormal = true
			} else if term := cs.LastTerminationState.Terminated; term != nil {
				if term.Reason == "OOMKilled" || (term.ExitCode != 0) {
					isAbnormal = true
				}
			} else if cs.RestartCount > 0 {
				isAbnormal = true
			}

			// 更新计数
			if isAbnormal {
				podFailedCounter[key]++
			} else {
				podFailedCounter[key] = 0
			}

			// 达到连续异常阈值才记录
			if podFailedCounter[key] >= consecutiveThreshold {
				out = append(out, newAbnormalPod(
					cluster,
					p.Namespace,
					p.Name,
					cs,
					"ConsecutiveRestart",
					fmt.Sprintf("Container has abnormal state for %d consecutive times", consecutiveThreshold),
				))
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
