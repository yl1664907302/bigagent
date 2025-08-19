package check

import (
	"bigagent/internal/kubernetes"
	result2 "bigagent/internal/result"
	"bigagent/internal/utils"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (d *AbnormalPod) EvictPod(cluster, namespace, podName string) error {
	// 获取k8s客户端
	k8sOp := d.k()
	c, _, err := k8sOp.Client(cluster)
	if err != nil {
		return err
	}
	if c == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("FetchPods.ClientSet.ByclusterName.empty[clusterName:%v]", cluster)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return err
	}

	// 先获取这个超时的ctx
	ctx, cancelFunc1 := d.GetK8sTwContext()
	defer cancelFunc1()

	// 获取需要处理pod对象
	pod, err := c.CoreV1().Pods(namespace).Get(ctx, podName, metaV1.GetOptions{})
	if err != nil {
		klog.Errorf("EvictPod.getPod.err[clusterName:%v][podNs:%v][podName:%v][err:%v]",
			cluster,
			namespace,
			podName,
			err,
		)
		return err
	}

	// 先获取这个超时的ctx
	ctx2, cancelFunc2 := d.GetK8sTwContext()
	defer cancelFunc2()
	var gracePeriodSeconds int64 = 0

	// 开始删除该pod,并配置deletePolicy(级联删除策略)与GracePeriodSeconds(优雅终止的宽限时间)
	//kubectl delete pod web-0 --grace-period=0
	deletePolicy := metaV1.DeletePropagationBackground
	err = c.CoreV1().Pods(namespace).Delete(ctx2, pod.Name, metaV1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
		PropagationPolicy:  &deletePolicy,
	})
	if err != nil {
		klog.Errorf("EvictPod.Delete.err[clusterName:%v][podNs:%v][podName:%v][err:%v]",
			cluster,
			namespace,
			podName,
			err)
		return err
	}
	klog.Infof("EvictPod.Delete.success[clusterName:%v][podNs:%v][podName:%v]",
		cluster,
		namespace,
		podName)
	return nil
}

func (d *AbnormalPod) FetchPods(cluster string, fieldSelector string) ([]coreV1.Pod, error) {
	// 获取k8s客户端
	k8sOp := d.k()
	c, _, err := k8sOp.Client(cluster)
	if err != nil {
		return nil, err
	}
	if c == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("FetchPods.ClientSet.ByclusterName.empty[clusterName:%v]", cluster)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return nil, err
	}

	// 获取ctx做超时
	ctx, cancelFunc1 := d.GetK8sTwContext()
	defer cancelFunc1()

	// 配置pod字段选择器
	options := metaV1.ListOptions{
		FieldSelector: fieldSelector,
	}

	// 获取pods
	pods, err := c.CoreV1().Pods("").List(ctx, options)
	if err != nil {
		klog.Errorf("FetchPods.listPod.err[clusterName:%v][fieldSelector:%v][err:%v]",
			cluster,
			fieldSelector,
			err,
		)
	}
	return pods.Items, err
}

// CheckPodTerminating 判断传入的旧 Pod 是否仍处于终止状态且未被重建，整合为 Result 返回
func (d *AbnormalPod) getPodTerminating(cluster string, oldPod *coreV1.Pod) result2.Result {
	severity := "critical"
	var items interface{}
	var count int
	if d.k == nil || oldPod == nil {
		return result2.NewResultPod(result2.Base{Cluster: cluster, Items: nil}, nil, "PodTerminating", severity, 0, nil)
	}
	k8sOp := d.k()
	cs, _, err := k8sOp.Client(cluster)
	if err != nil {
		return result2.NewResultPod(result2.Base{Cluster: cluster, Items: nil}, nil, "PodTerminating", severity, 0, nil)
	}
	// 获取该pod的最新状态
	newPod, err := cs.CoreV1().Pods(oldPod.Namespace).Get(d.ctx, oldPod.Name, metaV1.GetOptions{})
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
	return result2.NewResultPod(result2.Base{Cluster: cluster, Items: items}, nil, "PodTerminating", severity, count, items)
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
		if p.Status.Phase == coreV1.PodSucceeded {
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
