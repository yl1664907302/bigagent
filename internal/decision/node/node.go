package decision

import (
	"bigagent/internal/config/global"
	"bigagent/internal/kubernetes"
	model "bigagent/internal/model/k8s"
	"bigagent/internal/result"
	"bigagent/internal/utils"
	"context"
	"fmt"
	"github.com/gammazero/workerpool"
	"k8s.io/klog/v2"
	"time"
)

// AbnormalNodeDecision Pod 异常裁决器
type AbnormalNodeDecision struct {
	K         func() *kubernetes.DefaultK8sOperator
	Ctx       context.Context
	Cluster   string
	Namespace string
	result    result.Result
}

func NewAbnormalNodeDecision(res result.Result) *AbnormalNodeDecision {
	return &AbnormalNodeDecision{result: res}
}

// JudgeWith 通用的裁决入口，接受不同的恢复逻辑
func (d *AbnormalNodeDecision) JudgeWith(
	recoveryFn func(result.Result) error,
) func() {
	score, err := d.result.SetScore()
	if err != nil {
		utils.DefaultLogger.Errorf("设置分数失败node：%s", err.Error())
		return nil
	}
	if score > 2 {
		return func() {
			utils.DefaultLogger.Warnf("发现疑是异常node，开始处理")
			if err := recoveryFn(d.result); err != nil {
				utils.DefaultLogger.Errorf("失败处理node异常：%s", err.Error())
			}
		}
	}
	return nil
}

// ========================== 各类恢复逻辑 ==========================

// recoveryNeedDelete 处理需要强制删除的 Pod
func (d *AbnormalNodeDecision) recoveryDeadNode(res result.Result) error {
	recovery_ql_map := make(map[string]string)
	cordon_daily_limit_map := make(map[string]int)
	recovery_ql_map["recovery_ql_map"] = global.V.GetString("recovery_ql_map.file_system_read_only")
	recovery_ql_map["recovery_ql_map"] = global.V.GetString("recovery_ql_map.arp_too_many")
	cordon_daily_limit_map["cordon_daily_limit_map"] = global.V.GetInt("cordon_daily_limit_map.file_system_read_only")
	cordon_daily_limit_map["cordon_daily_limit_map"] = global.V.GetInt("cordon_daily_limit_map.arp_too_many")

	wp := workerpool.New(20)
	n2Ip, err := d.result.GetNode2Ip()
	if err != nil {
		return err
	}
	for nodeName, _ := range n2Ip {
		wp.Submit(func() {
			toDayStr := time.Now().Format("2006-01-02")
			// 需要根据模块名拿到配置 好的 对的recoveryQl
			checkName, _ := d.result.GetCheckName()
			//recoveryQl := recovery_ql_map[checkName]
			// 判断是否在db中
			//nodeDb := model.NodeMaintenance{
			//	Name:        nodeName,
			//	Ip:          ip,
			//	ClusterName: clusterName,
			//	ModuleName:  moduleName,
			//	Reason:      moduleName,
			//	FirstDate:   toDayStr,
			//	RecoveryQl:  recoveryQl,
			//}
			//ok, _ := nodeDb.CheckExist()
			//if ok {
			//	klog.Infof("CommonModuleDealOneNode.Already.Deal[cluster:%v][node:%v][moduleName:%v]", clusterName, nodeName, moduleName)
			//	return
			//}

			// 判断每日限流
			// 首先拿到这个 模块的限流数字
			limitNum := cordon_daily_limit_map[checkName]

			todayMNodes, err := model.GetDailyNodeMaintenances(toDayStr, checkName, d.Cluster)
			if err != nil {
				klog.Errorf("CommonModuleDealOneNode.getGetDailyNodeMaintenances.err:%v", err)
				return
			}

			// 如果限流被拦截 那也不能做
			if len(todayMNodes) > limitNum {
				msg := fmt.Sprintf("[日期:%v][集群:%v]\n][模块:%v 达到每日限流停止操作:%v][节点:%v]",

					toDayStr,
					d.Cluster,
					checkName,
					limitNum,
					nodeName,
				)

				klog.Infof(msg)
				//util.DingDingMsgDirectSend(gr.Cg.ModuleCommonC.ImDingDingC, msg)
				return
			}

			// 先cordon它
			err = d.CordonNode(nodeName, d.Cluster, checkName)
			cordonRes := "成功"
			if err != nil {
				cordonRes = fmt.Sprintf("失败：%v", err)
			}
			msg := fmt.Sprintf("[集群:%v]\n[%s 插件cordon节点：%v结果：%v ]\n[今日操作数:%v]",
				d.Cluster,
				checkName,
				nodeName,
				cordonRes,
				len(todayMNodes)+1,
			)
			klog.Infof(msg)
			// 先通知一下
			//util.DingDingMsgDirectSend(gr.Cg.ModuleCommonC.ImDingDingC, msg)

			//_, err = nodeDb.AddOrGetOne()
			//if err != nil {
			//	klog.Errorf("CommonModuleDealOneNode.addToDb.err[moduleName:%v][node:%v][err:%v]", checkName, nodeDb, err)
			//
			//}
			//klog.Infof("CommonModuleDealOneNode.addToDb.success[moduleName:%v][node:%v][err:%v]", moduleName, nodeDb.Name, err)
		})
	}
	wp.StopWait()
	return nil
}

// ========================== 对外暴露的裁决方法 ==========================

// JudgeNeedDelete Pod 强制删除处理
func (d *AbnormalNodeDecision) JudgeDeadNode() func() {
	return d.JudgeWith(d.recoveryDeadNode)
}
