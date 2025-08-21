package decision

import (
	check "bigagent/internal/check/pod"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (d *AbnormalNodeDecision) CordonNode(nodeName, clusterName, modelName string) error {
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

	// 在cordon前需要先获取一下这个节点，拿到节点对象
	node, err := c.CoreV1().Nodes().Get(ctx, nodeName, v1.GetOptions{})
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
	anno["disable-by"] = "bigagent"
	anno["disable-module"] += fmt.Sprintf("|%s", modelName)

	node.Annotations = anno
	// 先获取这个超时的ctx
	ctx2, cancelFunc2 := check.GetK8sTwContext()
	defer cancelFunc2()

	// 更新这个node
	node, err = c.CoreV1().Nodes().Update(ctx2, node, v1.UpdateOptions{})
	if err != nil {
		klog.Errorf("CordonNode.updateNode.err[clusterName:%v][nodeName:%v][modelName:%v][err:%v]", clusterName, nodeName, modelName, err)
		return err
	}
	klog.Infof("CordonNode.updateNode.success[clusterName:%v][nodeName:%v][modelName:%v]", clusterName, nodeName, modelName)
	return nil
}

func (d *AbnormalNodeDecision) UnCordonNode(nodeName, clusterName, modelName string) error {
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
	ctx1, cancelFunc1 := check.GetK8sTwContext()
	defer cancelFunc1()
	// 在cordon前需要先获取一下这个节点，拿到节点对象
	node, err := c.CoreV1().Nodes().Get(ctx1, nodeName, v1.GetOptions{})
	if err != nil {
		klog.Errorf("UnCordonNode.getNode.err[clusterName:%v][nodeName:%v][modelName:%v][err:%v]", clusterName, nodeName, modelName, err)
		return err
	}

	// 设置unschedule 字段
	node.Spec.Unschedulable = false

	// 先获取这个超时的ctx
	ctx2, cancelFunc2 := check.GetK8sTwContext()
	defer cancelFunc2()

	// 更新这个node
	node, err = c.CoreV1().Nodes().Update(ctx2, node, v1.UpdateOptions{})
	if err != nil {
		klog.Errorf("UnCordonNode.updateNode.err[clusterName:%v][nodeName:%v][modelName:%v][err:%v]", clusterName, nodeName, modelName, err)
		return err
	}
	klog.Infof("UnCordonNode.updateNode.success[clusterName:%v][nodeName:%v][modelName:%v]", clusterName, nodeName, modelName)
	return nil
}
