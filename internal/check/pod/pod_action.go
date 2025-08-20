package check

import (
	"bigagent/internal/kubernetes"
	"context"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (d *AbnormalPod) FetchPods(fieldSelector string) ([]coreV1.Pod, error) {
	// 获取k8s客户端
	k8sOp := d.k()
	c, _, err := k8sOp.Client(d.Cluster)
	if err != nil {
		return nil, err
	}
	if c == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("FetchPods.ClientSet.ByclusterName.empty[clusterName:%v]", d.Cluster)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return nil, err
	}

	// 获取ctx做超时
	ctx, cancelFunc1 := GetK8sTwContext()
	defer cancelFunc1()

	// 配置pod字段选择器
	options := metaV1.ListOptions{
		FieldSelector: fieldSelector,
	}

	// 获取pods
	pods, err := c.CoreV1().Pods("").List(ctx, options)
	if err != nil {
		klog.Errorf("FetchPods.listPod.err[clusterName:%v][fieldSelector:%v][err:%v]",
			d.Cluster,
			fieldSelector,
			err,
		)
	}
	return pods.Items, err
}

// CheckPodTerminating 判断传入的旧 Pod 是否仍处于终止状态且未被重建
func (d *AbnormalPod) GetPodTerminating(oldPod *coreV1.Pod) bool {

	// 获取k8s客户端
	k8sOp := d.k()
	c, _, err := k8sOp.Client(d.Cluster)
	if err != nil {
		return false
	}
	if c == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("FetchPods.ClientSet.ByclusterName.empty[clusterName:%v]", d.Cluster)
		klog.Errorf(msg)
		return false
	}

	// 将该pod重新获取一遍
	newPod, _ := c.CoreV1().Pods(oldPod.Namespace).Get(d.ctx, oldPod.Name, metaV1.GetOptions{})

	if newPod == nil {
		// 已经被删掉
		return false
	}
	if newPod.UID != oldPod.UID {
		// 重新创建了
		return false
	}

	//pod被k8s删除事会被打上删除时间戳，若为零值，说明 Pod 没（或正在）被删除
	if newPod.DeletionTimestamp.IsZero() {
		return false
	}

	return true
}

// GetPodTerminating 判断传入的旧 Pod 是否仍处于终止状态且未被重建 (外部函数调用版)
func GetPodTerminating(ctx context.Context, op func() *kubernetes.DefaultK8sOperator, Cluster string, oldPod *coreV1.Pod) bool {

	// 获取k8s客户端
	k8sOp := op()
	c, _, err := k8sOp.Client(Cluster)
	if err != nil {
		return false
	}
	if c == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("FetchPods.ClientSet.ByclusterName.empty[clusterName:%v]", Cluster)
		klog.Errorf(msg)
		return false
	}

	// 将该pod重新获取一遍
	newPod, _ := c.CoreV1().Pods(oldPod.Namespace).Get(ctx, oldPod.Name, metaV1.GetOptions{})

	if newPod == nil {
		// 已经被删掉
		return false
	}
	if newPod.UID != oldPod.UID {
		// 重新创建了
		return false
	}

	//pod被k8s删除事会被打上删除时间戳，若为零值，说明 Pod 没（或正在）被删除
	if newPod.DeletionTimestamp.IsZero() {
		return false
	}

	return true
}
