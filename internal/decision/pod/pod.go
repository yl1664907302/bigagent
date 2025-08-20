package decision

import (
	check "bigagent/internal/check/pod"
	"bigagent/internal/kubernetes"
	"bigagent/internal/result"
	"bigagent/internal/utils"
	"context"
	"fmt"
	"github.com/gammazero/workerpool"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"time"
)

// AbnormalPodDecision Pod 异常裁决器
type AbnormalPodDecision struct {
	K         func() *kubernetes.DefaultK8sOperator
	Ctx       context.Context
	Cluster   string
	Namespace string
	result    result.Result
}

func NewAbnormalPodDecision(res result.Result) *AbnormalPodDecision {
	return &AbnormalPodDecision{result: res}
}

// JudgeWith 通用的裁决入口，接受不同的恢复逻辑
func (d *AbnormalPodDecision) JudgeWith(
	recoveryFn func(result.Result) error,
) func() {
	score, err := d.result.SetScore()
	if err != nil {
		utils.DefaultLogger.Errorf("设置分数失败pod：%s", err.Error())
		return nil
	}
	if score > 2 {
		return func() {
			utils.DefaultLogger.Warnf("发现疑是异常pod，开始处理")
			if err := recoveryFn(d.result); err != nil {
				utils.DefaultLogger.Errorf("失败处理pod异常：%s", err.Error())
			}
		}
	}
	return nil
}

// ========================== 各类恢复逻辑 ==========================

// recoveryPods 处理多 Pod 异常（重启过多等场景）
func (d *AbnormalPodDecision) recoveryPods(res result.Result) error {
	item, err := res.GetItem()
	if err != nil {
		return err
	}
	pods, ok := item.([]check.AbnormalPod)
	if !ok {
		return nil
	}
	for _, pod := range pods {
		utils.DefaultLogger.Warnf("处理pod异常: 集群 %s, 命名空间 %s, Pod %s, 容器 %s, 重启次数 %d, 原因 %s, 信息 %s",
			pod.Cluster, pod.Namespace, pod.Pod, pod.Container, pod.RestartCount, pod.Reason, pod.Message)
	}
	return nil
}

// recoveryTem 处理单 Pod 处于 Terminating 状态
func (d *AbnormalPodDecision) recoveryTem(res result.Result) error {
	item, err := res.GetItem()
	if err != nil {
		return err
	}
	pod, ok := item.(check.AbnormalPod)
	if !ok {
		return nil
	}
	if pod.Terminating {
		utils.DefaultLogger.Warnf("正在删除pod %s", pod.Pod)
	} else {
		utils.DefaultLogger.Warnf("Pod %s 失效，无效处理", pod.Pod)
	}
	return nil
}

// recoveryNeedDelete 处理需要强制删除的 Pod
func (d *AbnormalPodDecision) recoveryNeedDelete(res result.Result) error {
	item, err := res.GetItem()
	if err != nil {
		return err
	}
	pods, ok := item.([]v1.Pod)
	if !ok {
		return nil
	}

	wp := workerpool.New(20)
	for _, pod := range pods {
		//log.Println(pod)
		wp.Submit(func() {
			//log.Println(pod.DeletionTimestamp.IsZero())
			if !pod.DeletionTimestamp.IsZero() {
				// 等待5秒
				time.Sleep(time.Duration(5) * time.Second)
				checkOk := check.GetPodTerminating(d.Ctx, d.K, d.Cluster, &pod)
				if !checkOk {
					// 正常的删除
					return
				}
				var err error
				action := ""
				switch len(pod.Finalizers) {
				case 0:
					// 尝试强制删除
					err = d.EvictPod(pod.Namespace, pod.Name)
					action = "尝试强制删除"
				default:
					// 说明有finalizer

					//if time.Now().Sub(pod.DeletionTimestamp.Time) < time.Hour {
					//	return
					//
					//}

					// 异常的删除  调用pod的patch remove finalizer
					err = d.RemovePodFinalizer(pod.Namespace, pod.Name)
					action = "尝试patch-finalizer"

				}
				res := "成功"
				if err != nil {
					res = fmt.Sprintf("失败 %v", err)
				}
				msg := fmt.Sprintf("[集群:%v][ns:%v][pod:%v]\n[长期删不掉的pod 执行动作:%v 结果:%v]",
					d.Cluster,
					pod.Namespace,
					pod.Name,
					action,
					res)
				klog.Infof(msg)
				//toDayStr := time.Now().Format("2006-01-02")
				//util.DingDingMsgDirectSend(gr.Cg.AbnormalPodCleanC.ImDingDingC, msg)
				//podDb := models.PodMaintenance{
				//	Name:        pod.Name,
				//	Ip:          pod.Status.PodIP,
				//	NodeName:    pod.Spec.NodeName,
				//	PodNs:       pod.Namespace,
				//	ClusterName: clusterName,
				//	ModuleName:  "abnormal_pod_clean",
				//	Reason:      action,
				//	FirstDate:   toDayStr,
				//}
				//_, err = podDb.AddOrGetOne()
				klog.Infof("podDb.AddOrGetOne.err:%v[pod%v]", err, pod.Name)
			}
			utils.DefaultLogger.Warnf("Pod %s 正常，未被k8s执行删除", pod.Name)
		})

	}
	wp.StopWait()
	return nil
}

// ========================== 对外暴露的裁决方法 ==========================

// Judge Pod 异常处理（多 Pod）
func (d *AbnormalPodDecision) Judge() func() {
	return d.JudgeWith(d.recoveryPods)
}

// JudgeTem Pod 终止处理
func (d *AbnormalPodDecision) JudgeTem() func() {
	return d.JudgeWith(d.recoveryTem)
}

// JudgeNeedDelete Pod 强制删除处理
func (d *AbnormalPodDecision) JudgeNeedDelete() func() {
	return d.JudgeWith(d.recoveryNeedDelete)
}
