package guard

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (gr *Guarder) DrainNode(nodeName, clusterName, modelName string) error {

	//	clusterName-->clientSet的map  找到这个clientSet
	kClient := gr.ClientSetMap[clusterName]
	if kClient == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("CordonNode.ClientSet.ByclusterName.empty[clusterName:%v][nodeName:%v][modelName:%v]", clusterName, nodeName, modelName)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return err
	}

	// 首先要找到这个node上的所有pod 类似 kubectl describe node
	// 先获取这个超时的ctx
	ctx1, cancelFunc1 := gr.GetK8sTwContext()
	defer cancelFunc1()
	// field-selector 字段选择器
	listOps := v1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	}
	pods, err := kClient.CoreV1().Pods("").List(ctx1, listOps)
	if err != nil {
		klog.Errorf("DrainNode.ListPod.err[clusterName:%v][nodeName:%v][modelName:%v][err:%v]",
			clusterName,
			nodeName,
			modelName,
			err)
		return err
	}
	num := len(pods.Items)
	if num == 0 {
		msg := fmt.Sprintf("can not found pod on node[clusterName:%v][nodeName:%v][modelName:%v]",
			clusterName,
			nodeName,
			modelName)

		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return err
	}

	// 拿到所有的pods之后遍历处理

	for index, pod := range pods.Items {
		// 遍历OwnerReferences数组
		dsManaged := false
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == "DaemonSet" {
				dsManaged = true
				break
			}
		}

		klog.Infof("[DrainNode.ListPod.print][%d/%d][clusterName:%v][nodeName:%v][modelName:%v][ns:%v][podName:%v][control_by_ds:%v]",

			index+1,
			num,
			clusterName,
			nodeName,
			modelName,
			pod.Namespace,
			pod.Name,
			dsManaged,
		)
		// 如果是被ds管理的忽略它
		if dsManaged {
			continue
		}
		// 删除pod
		err := gr.EvictPod(clusterName, pod.Namespace, pod.Name)
		if err != nil {
			klog.Error("[DrainNode.EvictPod.err][%d/%d][clusterName:%v][nodeName:%v][modelName:%v][ns:%v][podName:%v][control_by_ds:%v][err:%v]",

				index+1,
				num,
				clusterName,
				nodeName,
				modelName,
				pod.Namespace,
				pod.Name,
				dsManaged,
				err,
			)
			return err
		}

	}
	return nil
}

func (gr *Guarder) TestDrainNode() {
	gr.DrainNode("k8s-node02", "cpu-compute-01", "ntp-check")
}
