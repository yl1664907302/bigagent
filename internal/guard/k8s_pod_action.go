package guard

import (
	"context"
	"encoding/json"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

func (gr *Guarder) EvictPod(clusterName, podNs, podName string) error {

	//	clusterName-->clientSet的map  找到这个clientSet
	kClient := gr.ClientSetMap[clusterName]
	if kClient == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("EvictPod.ClientSet.ByclusterName.empty[clusterName:%v][podName:%v]", clusterName, podName)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return err
	}

	// 先去get pod

	// 先获取这个超时的ctx
	ctx1, cancelFunc1 := gr.GetK8sTwContext()
	defer cancelFunc1()
	// 在cordon前需要先获取一下这个节点，拿到节点对象
	pod, err := kClient.CoreV1().Pods(podNs).Get(ctx1, podName, v1.GetOptions{})
	if err != nil {
		klog.Errorf("EvictPod.getPod.err[clusterName:%v][podNs:%v][podName:%v][err:%v]",
			clusterName,
			podNs,
			podName,
			err,
		)
		return err
	}

	// 先获取这个超时的ctx
	ctx2, cancelFunc2 := gr.GetK8sTwContext()
	defer cancelFunc2()
	var gracePeriodSeconds int64 = 0
	var value *int64
	value = &gracePeriodSeconds
	// 类似删除的时候 grace-period=0
	//kubectl delete pod web-0 --grace-period=0
	deletePolicy := v1.DeletePropagationBackground
	err = kClient.CoreV1().Pods(podNs).Delete(ctx2, pod.Name, v1.DeleteOptions{
		GracePeriodSeconds: value,
		PropagationPolicy:  &deletePolicy,
	})
	if err != nil {
		klog.Errorf("EvictPod.Delete.err[clusterName:%v][podNs:%v][podName:%v][err:%v]",
			clusterName,
			podNs,
			podName,
			err)
		return err
	}
	klog.Infof("EvictPod.Delete.success[clusterName:%v][podNs:%v][podName:%v]",
		clusterName,
		podNs,
		podName)
	return nil
}

// pod 重新创建
func (gr *Guarder) CheckPodTerminating(clusterName string, oldPod *coreV1.Pod) bool {

	kClient := gr.ClientSetMap[clusterName]
	if kClient == nil {
		return false
	}

	newPod, _ := kClient.CoreV1().Pods(oldPod.Namespace).Get(context.Background(), oldPod.Name, v1.GetOptions{})

	if newPod == nil {
		// 已经被删掉
		return false
	}
	if newPod.UID != oldPod.UID {
		// 重新创建了
		return false
	}

	if newPod.DeletionTimestamp.IsZero() {
		return false
	}

	return true
}

type RemoveStringValue struct {
	Op   string `json:"op"`
	Path string `json:"path"`
}

func (gr *Guarder) RemovePodFinalizer(clusterName, podNs, podName string) error {

	//	clusterName-->clientSet的map  找到这个clientSet
	kClient := gr.ClientSetMap[clusterName]
	if kClient == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("RemovePodFinalizer.ClientSet.ByclusterName.empty[clusterName:%v][podName:%v]", clusterName, podName)
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
	ctx2, cancelFunc2 := gr.GetK8sTwContext()
	defer cancelFunc2()
	_, err := kClient.CoreV1().Pods(podNs).Patch(ctx2, podName, types.JSONPatchType, data, v1.PatchOptions{})
	return err
}

func (gr *Guarder) FetchPods(clusterName, fieldSelector string) ([]coreV1.Pod, error) {

	//	clusterName-->clientSet的map  找到这个clientSet
	kClient := gr.ClientSetMap[clusterName]
	if kClient == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("FetchPods.ClientSet.ByclusterName.empty[clusterName:%v]", clusterName)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return nil, err
	}

	// 先去get pod

	// 先获取这个超时的ctx
	ctx1, cancelFunc1 := gr.GetK8sTwContext()
	defer cancelFunc1()
	listOps := v1.ListOptions{
		//LabelSelector:        "",
		FieldSelector: fieldSelector,
	}
	// 在cordon前需要先获取一下这个节点，拿到节点对象
	pods, err := kClient.CoreV1().Pods("").List(ctx1, listOps)
	if err != nil {
		klog.Errorf("FetchPods.listPod.err[clusterName:%v][fieldSelector:%v][err:%v]",
			clusterName,
			fieldSelector,
			err,
		)

	}
	return pods.Items, err

}

func (gr *Guarder) TestEvictPod() {
	err := gr.EvictPod("cpu-compute-01", "default", "dns-core-test")
	fmt.Println(err)
}
