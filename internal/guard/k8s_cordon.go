package guard

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// 我们需要根据  clientSet 对象操作这个集群
// clusterName-->clientSet的map
// 我得知道是哪个模块使它cordon的
func (gr *Guarder) CordonNode(nodeName, clusterName, modelName string) error {

	//	clusterName-->clientSet的map  找到这个clientSet
	kClient := gr.ClientSetMap[clusterName]
	if kClient == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("CordonNode.ClientSet.ByclusterName.empty[clusterName:%v][nodeName:%v][modelName:%v]", clusterName, nodeName, modelName)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return err
	}

	// 先获取这个超时的ctx
	ctx1, cancelFunc1 := gr.GetK8sTwContext()
	defer cancelFunc1()
	// 在cordon前需要先获取一下这个节点，拿到节点对象
	node, err := kClient.CoreV1().Nodes().Get(ctx1, nodeName, v1.GetOptions{})
	if err != nil {
		klog.Errorf("CordonNode.getNode.err[clusterName:%v][nodeName:%v][modelName:%v][err:%v]", clusterName, nodeName, modelName, err)
		return err
	}

	// 设置unschedule 字段
	node.Spec.Unschedulable = true
	anno := node.Annotations
	if anno == nil {
		anno = make(map[string]string)
	}
	anno["disable-by"] = "k8s-cluster-guard"
	anno["disable-module"] += fmt.Sprintf("|%s", modelName)

	node.Annotations = anno
	// 先获取这个超时的ctx
	ctx2, cancelFunc2 := gr.GetK8sTwContext()
	defer cancelFunc2()

	// 更新这个node
	node, err = kClient.CoreV1().Nodes().Update(ctx2, node, v1.UpdateOptions{})
	if err != nil {
		klog.Errorf("CordonNode.updateNode.err[clusterName:%v][nodeName:%v][modelName:%v][err:%v]", clusterName, nodeName, modelName, err)
		return err
	}
	klog.Infof("CordonNode.updateNode.success[clusterName:%v][nodeName:%v][modelName:%v]", clusterName, nodeName, modelName)
	return nil
}

func (gr *Guarder) UnCordonNode(nodeName, clusterName, modelName string) error {

	//	clusterName-->clientSet的map  找到这个clientSet
	kClient := gr.ClientSetMap[clusterName]
	if kClient == nil {
		// 打印错误和返回都需要用
		msg := fmt.Sprintf("UnCordonNode.ClientSet.ByclusterName.empty[clusterName:%v][nodeName:%v][modelName:%v]", clusterName, nodeName, modelName)
		err := fmt.Errorf(msg)
		klog.Errorf(msg)
		return err
	}

	// 先获取这个超时的ctx
	ctx1, cancelFunc1 := gr.GetK8sTwContext()
	defer cancelFunc1()
	// 在cordon前需要先获取一下这个节点，拿到节点对象
	node, err := kClient.CoreV1().Nodes().Get(ctx1, nodeName, v1.GetOptions{})
	if err != nil {
		klog.Errorf("UnCordonNode.getNode.err[clusterName:%v][nodeName:%v][modelName:%v][err:%v]", clusterName, nodeName, modelName, err)
		return err
	}

	// 设置unschedule 字段
	node.Spec.Unschedulable = false

	// 先获取这个超时的ctx
	ctx2, cancelFunc2 := gr.GetK8sTwContext()
	defer cancelFunc2()

	// 更新这个node
	node, err = kClient.CoreV1().Nodes().Update(ctx2, node, v1.UpdateOptions{})
	if err != nil {
		klog.Errorf("UnCordonNode.updateNode.err[clusterName:%v][nodeName:%v][modelName:%v][err:%v]", clusterName, nodeName, modelName, err)
		return err
	}
	klog.Infof("UnCordonNode.updateNode.success[clusterName:%v][nodeName:%v][modelName:%v]", clusterName, nodeName, modelName)
	return nil
}

//func (gr *Guarder) TestCordon() {
//
//	gr.CordonNode("k8s-node01", "cpu-compute-01", "ntp-check")
//
//}
//
//func (gr *Guarder) TestUnCordon() {
//
//	gr.UnCordonNode("k8s-node01", "cpu-compute-01", "ntp-check")
//
//}
