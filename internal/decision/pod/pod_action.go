package decision

import (
	check "bigagent/internal/check/pod"
	"encoding/json"
	"fmt"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

type RemoveStringValue struct {
	Op   string `json:"op"`
	Path string `json:"path"`
}

// RemovePodFinalizer 移除 Pod 的 finalizer (收尾操作移除)
func (d *AbnormalPodDecision) RemovePodFinalizer(podNs, podName string) error {

	// 获取k8s客户端
	k8sOp := d.K()
	c, _, err := k8sOp.Client(d.Cluster)
	if err != nil {
		return err
	}
	if c == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("FetchPods.ClientSet.ByclusterName.empty[clusterName:%v]", d.Cluster)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return err
	}

	var payloads []interface{}
	removeFinalizer := RemoveStringValue{
		Op:   "remove",
		Path: "/metadata/finalizers",
	}

	payloads = append(payloads, removeFinalizer)
	data, _ := json.Marshal(payloads)

	// 先获取这个超时的ctx
	ctx, cancelFunc := check.GetK8sTwContext()
	defer cancelFunc()

	// 对pod以打补丁的方式进行更新，告诉 Kubernetes API：“移除对象 metadata.finalizers 字段”。
	_, err = c.CoreV1().Pods(podNs).Patch(ctx, podName, types.JSONPatchType, data, metaV1.PatchOptions{})
	return err
}

func (d *AbnormalPodDecision) EvictPod(namespace, podName string) error {
	// 获取k8s客户端
	k8sOp := d.K()
	c, _, err := k8sOp.Client(d.Cluster)
	if err != nil {
		return err
	}
	if c == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("FetchPods.ClientSet.ByclusterName.empty[clusterName:%v]", d.Cluster)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return err
	}

	// 先获取这个超时的ctx
	ctx, cancelFunc1 := check.GetK8sTwContext()
	defer cancelFunc1()

	// 获取需要处理pod对象
	pod, err := c.CoreV1().Pods(namespace).Get(ctx, podName, metaV1.GetOptions{})
	if err != nil {
		klog.Errorf("EvictPod.getPod.err[clusterName:%v][podNs:%v][podName:%v][err:%v]",
			d.Cluster,
			namespace,
			podName,
			err,
		)
		return err
	}

	// 先获取这个超时的ctx
	ctx2, cancelFunc2 := check.GetK8sTwContext()
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
			d.Cluster,
			namespace,
			podName,
			err)
		return err
	}
	klog.Infof("EvictPod.Delete.success[clusterName:%v][podNs:%v][podName:%v]",
		d.Cluster,
		namespace,
		podName)
	return nil
}
