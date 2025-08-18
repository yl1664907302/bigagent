package guard

import (
	model "bigagent/internal/model/k8s"
	"bigagent/internal/utils"
	"context"
	"fmt"
	"github.com/gammazero/workerpool"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

// wait的封装
func (gr *Guarder) RunModulePodCleanManager(ctx context.Context) error {
	// (ctxAll,执行的方法，每隔多长时间)
	go wait.UntilWithContext(ctx, gr.RunModuleCleanPod, time.Duration(gr.Cg.AbnormalPodCleanC.CheckIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Infof("RunModulePodManager.exit.receive_quit_signal")
	return nil

}

func (gr *Guarder) RunModuleCleanPod(ctx context.Context) {
	var (
		wg sync.WaitGroup
	)

	// 遍历所有集群的prometheus-api map
	for clusterName := range gr.Cg.AbnormalPodCleanC.EnabledClusters {

		clusterName := clusterName

		wg.Add(1)
		go func() {
			defer wg.Done()
			gr.RunModuleCleanPodOnCluster(clusterName)
		}()
	}
	wg.Wait()
}

func (gr *Guarder) RunModuleCleanPodOnCluster(clusterName string) {

	// 需要先获取异常的pod
	pods, err := gr.FetchPods(clusterName, gr.Cg.AbnormalPodCleanC.FieldSelector)
	if err != nil {
		return
	}
	num := len(pods)
	if num == 0 {
		klog.Infof("RunModuleCleanPodOnCluster.abnormalPod.zero")
	}

	// 遍历pod

	wp := workerpool.New(20)

	for _, pod := range pods {

		wp.Submit(func() {
			gr.RunModuleCleanDealWithOnePod(pod, clusterName)

		})

	}

	wp.StopWait()
}

func (gr *Guarder) RunModuleCleanDealWithOnePod(pod coreV1.Pod, clusterName string) {
	// 先来判断是否是要删除
	if !pod.DeletionTimestamp.IsZero() {

		time.Sleep(time.Duration(gr.Cg.AbnormalPodCleanC.DoubleCheckSecSeconds) * time.Second)
		checkOk := gr.CheckPodTerminating(clusterName, &pod)
		if !checkOk {
			// 正常的删除
			return
		}
		var err error
		action := ""
		switch len(pod.Finalizers) {
		case 0:

			// 尝试强制删除
			err = gr.EvictPod(clusterName, pod.Namespace, pod.Name)
			action = "尝试强制删除"

		default:
			// 说明有finalizer

			//if time.Now().Sub(pod.DeletionTimestamp.Time) < time.Hour {
			//	return
			//
			//}
			// 异常的删除  调用pod的patch remove finalizer
			err = gr.RemovePodFinalizer(clusterName, pod.Namespace, pod.Name)
			action = "尝试patch-finalizer"

		}
		res := "成功"
		if err != nil {
			res = fmt.Sprintf("失败 %v", err)
		}
		msg := fmt.Sprintf("[集群:%v][ns:%v][pod:%v]\n[长期删不掉的pod 执行动作:%v 结果:%v]",
			clusterName,
			pod.Namespace,
			pod.Name,
			action,
			res)
		klog.Infof(msg)
		toDayStr := time.Now().Format("2006-01-02")
		utils.DingDingMsgDirectSend(gr.Cg.AbnormalPodCleanC.ImDingDingC, msg)
		podDb := model.PodMaintenance{
			Name:        pod.Name,
			Ip:          pod.Status.PodIP,
			NodeName:    pod.Spec.NodeName,
			PodNs:       pod.Namespace,
			ClusterName: clusterName,
			ModuleName:  "abnormal_pod_clean",
			Reason:      action,
			FirstDate:   toDayStr,
		}
		_, err = podDb.AddOrGetOne()
		klog.Infof("podDb.AddOrGetOne.err:%v[pod%v]", err, pod.Name)

		return

	}

	//switch pod.Status.Phase {
	//case coreV1.PodRunning, coreV1.PodSucceeded, coreV1.PodFailed:
	//	return
	//
	//default:
	//	msg := fmt.Sprintf("[pod_info][cluster:%v][pod_ns:%v][pod_name:%v][phase:%v][node_name:%v]",
	//		clusterName,
	//		pod.Namespace,
	//		pod.Name,
	//		pod.Status.Phase,
	//		pod.Spec.NodeName,
	//	)
	//
	//	klog.Infof(msg)
	//}
}
