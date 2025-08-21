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

const (
	MOD_NODE_DOWN        = "node_down"
	MOD_NODE_DOWN_REASON = "node_down_guard"
	POWER_ACTION_OFF     = "POWER_OFF"
	POWER_ACTION_ON      = "POWER_ON"
	PLUGIN_NTP           = "ntp_time"
	PLUGIN_NTP_REASON    = "ntp_time"
	MOD_VOLUME_REASON    = "vm_diff_csi"
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

func (d *AbnormalNodeDecision) recoveryNtpNodeToCordon(res result.Result) error {
	node2Ip, err := d.result.GetNode2Ip()
	if err != nil {
		return err
	}
	var index = 0
	num := len(node2Ip)
	for node, ip := range node2Ip {
		klog.Infof("[RunNtpCheckOnCluster.QueryRes.print][clusterName:%v][%d/%d][nodeName:%v]",
			d.Cluster,
			index+1,
			num,
			node,
		)
		toDayStr := time.Now().Format("2006-01-02")
		// 判断是否在db中
		nodeDb := model.NodeMaintenance{
			Name:        node,
			Ip:          ip,
			ClusterName: d.Cluster,
			ModuleName:  PLUGIN_NTP,
			Reason:      PLUGIN_NTP,
			FirstDate:   toDayStr,
			RecoveryQl:  global.V.GetString("plugin_ntp.recovery_ql"),
		}
		ok, _ := nodeDb.CheckExist()
		if ok {
			klog.Infof("PluginNtpDealOneNode.Already.Deal[cluster:%v][node:%v]", d.Cluster, node)
			return nil
		}

		// 判断每日限流
		todayMNodes, err := model.GetDailyNodeMaintenances(toDayStr, PLUGIN_NTP, d.Cluster)
		if err != nil {
			return fmt.Errorf("PluginNtpDealOneNode.getGetDailyNodeMaintenances.err:%v", err)
		}

		// 如果限流被拦截 那也不能做
		if len(todayMNodes) > global.V.GetInt("plugin_ntp.cordon_daily_limit") {
			msg := fmt.Sprintf("[日期:%v][集群:%v]\n][模块:%v 达到每日限流停止操作:%v][节点:%v]",

				toDayStr,
				d.Cluster,
				PLUGIN_NTP,
				global.V.GetInt("plugin_ntp.cordon_daily_limit"),
				node,
			)

			klog.Infof(msg)
			//util.DingDingMsgDirectSend(gr.Cg.PluginNtpC.ImDingDingC, msg)
			return nil
		}

		// 先cordon它
		err = d.CordonNode(node, d.Cluster, PLUGIN_NTP)
		cordonRes := "成功"
		if err != nil {
			cordonRes = fmt.Sprintf("失败：%v", err)
		}
		msg := fmt.Sprintf("[集群:%v]\n[ntp插件cordon节点：%v结果：%v ]\n[今日操作数:%v]",
			d.Cluster,
			node,
			cordonRes,
			len(todayMNodes)+1,
		)
		klog.Infof(msg)
		// 先通知一下
		//util.DingDingMsgDirectSend(gr.Cg.PluginNtpC.ImDingDingC, msg)
		_, err = nodeDb.AddOrGetOne()
		if err != nil {
			klog.Errorf("PluginNtpDealOneNode.addToDb.err[node:%v][err:%v]", nodeDb, err)

		}
		klog.Infof("PluginNtpDealOneNode.addToDb.success[node:%v][err:%v]", nodeDb.Name, err)
	}
	return nil
}

func (d *AbnormalNodeDecision) recoveryDownNodeToCordon(res result.Result) error {
	item, err := res.GetItem()
	if err != nil {
		return err
	}
	realDownNodeWithIps := item.(map[string]string)
	realDownNodesMsg := fmt.Sprintf("[%s]\n[集群:%v][宕机总数:%v]\n",
		global.V.GetString("node_down.im_ding_ding.title"),

		d.Cluster, len(realDownNodeWithIps))
	newFound := 0
	for nodeName, nodeIp := range realDownNodeWithIps {

		toDayStr := time.Now().Format("2006-01-02")
		// cordon node

		err := d.CordonNode(nodeName, d.Cluster, MOD_NODE_DOWN)
		if err != nil {
			klog.Errorf("CordonNode.err[cluster:%v][node:%v][err:%v]", d.Cluster, nodeName, err)
			//continue
		}

		//  查询节点ip和sn号

		// 更新到维修db中
		nodeInDb := &model.NodeMaintenance{
			Name:        nodeName,
			ClusterName: d.Cluster,
			ModuleName:  MOD_NODE_DOWN,
			Ip:          nodeIp,
		}
		exist, _ := nodeInDb.CheckExist()
		if exist {
			continue
		}
		newFound++
		realDownNodesMsg += fmt.Sprintf("[node:%v][ip:%v][sn:]\n", nodeName, nodeIp)
		nodeInDb.Reason = MOD_NODE_DOWN_REASON
		nodeInDb.FirstDate = toDayStr
		_, err = nodeInDb.AddOne()

		if err != nil {
			klog.Errorf("NodeDownCheckOneCluster.abnormal.cordonNode.AddOrGetOne.err[cluster:%v][node:%v][err:%v]", d.Cluster, nodeName, err)

		}
		klog.Infof("NodeDownCheckOneCluster.abnormal.cordonNode.AddOrGetOne.success[cluster:%v][node:%v]", d.Cluster, nodeName)

	}
	if newFound > 0 {
		klog.Infof("NodeDownCheckOneCluster.get.realDownNodes.print[realDownNodesMsg:%v][cluster:%v][num:%v][detail:%v]", realDownNodesMsg, d.Cluster, len(realDownNodeWithIps), realDownNodeWithIps)
		//util.DingDingMsgDirectSend(gr.Cg.NodeDownC.ImDingDingC, realDownNodesMsg
	}
	return nil
}

func (d *AbnormalNodeDecision) recoveryCordonToUnCordon(res result.Result) error {
	item, err := res.GetItem()
	if err != nil {
		return err
	}
	nodes, ok := item.([]model.NodeMaintenance)
	if !ok {
		return nil
	}
	wp := workerpool.New(20)
	for _, node := range nodes {
		node := node
		wp.Submit(func() {
			err = d.UnCordonNode(node.Name, d.Cluster, node.ModuleName)
			if err != nil {
				klog.Errorf("recoveryCordonToUnCordon.UnCordonNode.err[clusterName:%v][nodeName:%v][err:%v]", d.Cluster, node.Name, err)
				return
			}
			klog.Infof("recoveryCordonToUnCordon.UnCordonNode.success[clusterName:%v][nodeName:%v]", d.Cluster, node.Name)
		})
	}
	wp.StopWait()
	return nil
}

// recoveryNeedDelete
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
	for nodeName, ip := range n2Ip {
		wp.Submit(func() {
			toDayStr := time.Now().Format("2006-01-02")
			// 需要根据模块名拿到配置 好的 对的recoveryQl
			checkName, _ := d.result.GetCheckName()
			recoveryQl := recovery_ql_map[checkName]
			// 幂等查询是否已经处理过了，在数据库已记录
			nodeDb := model.NodeMaintenance{
				Name:        nodeName,
				Ip:          ip,
				ClusterName: d.Cluster,
				ModuleName:  checkName,
				Reason:      checkName,
				FirstDate:   toDayStr,
				RecoveryQl:  recoveryQl,
			}
			ok, _ := nodeDb.CheckExist()
			if ok {
				klog.Infof("CommonModuleDealOneNode.Already.Deal[cluster:%v][node:%v][moduleName:%v]", d.Cluster, nodeName, checkName)
				return
			}

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

			_, err = nodeDb.AddOrGetOne()
			if err != nil {
				klog.Errorf("CommonModuleDealOneNode.addToDb.err[moduleName:%v][node:%v][err:%v]", checkName, nodeDb, err)

			}
			klog.Infof("CommonModuleDealOneNode.addToDb.success[moduleName:%v][node:%v][err:%v]", checkName, nodeDb.Name, err)
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

func (d *AbnormalNodeDecision) JudgeCordonToUnCordon() func() {
	return d.JudgeWith(d.recoveryCordonToUnCordon)
}

func (d *AbnormalNodeDecision) JudgeDownNodeToCordon() func() {
	return d.JudgeWith(d.recoveryDownNodeToCordon)
}

func (d *AbnormalNodeDecision) JudgeNtpNodeToCordon() func() {
	return d.JudgeWith(d.recoveryNtpNodeToCordon)
}
