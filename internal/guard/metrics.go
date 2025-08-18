package guard

import (
	model "bigagent/internal/model/k8s"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
	"strings"
)

var (

	// 模块状态 ：开启 关闭
	ModuleInfo = prometheus.NewDesc(
		"k8s_cluster_guard_module_info",
		"k8s_cluster_guard_module_info",
		[]string{
			"module_name",
			"enable",
			"clusters",
		},
		nil)

	// 集群状态 ：开启 关闭
	ClusterInfo = prometheus.NewDesc(
		"k8s_cluster_guard_cluster_info",
		"k8s_cluster_guard_cluster_info",
		[]string{
			"cluster_name",
			"prometheus_api",
		},
		nil)

	NodeCordonInfo = prometheus.NewDesc(
		"k8s_guard_node_cordon_info",
		"k8s_guard_node_cordon_info",
		[]string{
			"target_cluster",
			"target_node",
			"module",
			"cordon_reason",
			"first_date",
			"create_time",
			"update_time",
			"ip",
			"sn",
			"swp_ticket_url",
			"has_restart",
			"has_repair",
			"has_recovery",
		},
		nil)

	PodDeleteInfo = prometheus.NewDesc(
		"k8s_guard_delete_pod_info",
		"k8s_guard_delete_pod_info",
		[]string{
			"target_cluster",
			"target_node",
			"pod_name",
			"pod_ns",
			"module",
			"reason",
			"first_date",
			"create_time",
			"update_time",
		},
		nil)
)

func (gr *Guarder) Collect(ch chan<- prometheus.Metric) {

	keyFunc := func(m map[string]string) string {
		tmp := []string{}
		for k := range m {
			tmp = append(tmp, k)
		}
		return strings.Join(tmp, ",")

	}

	// 基础信息从配置文件中拿
	// ntp/通用/宕机/pod清理
	// xx模块|是否开启 --> 集群1,集群2
	modInfoMap := map[string]string{}
	enableMap := map[bool]string{
		true:  "开启",
		false: "关闭",
	}

	modNtp := gr.Cg.PluginNtpC
	modNtpKey := fmt.Sprintf("ntp模块|%s", enableMap[modNtp.Enable])
	modInfoMap[modNtpKey] = keyFunc(modNtp.EnabledClusters)

	modCommon := gr.Cg.ModuleCommonC
	modCommonKey := fmt.Sprintf("异常通用模块|%s", enableMap[modCommon.Enable])
	modInfoMap[modCommonKey] = keyFunc(modCommon.EnabledClusters)

	modPod := gr.Cg.AbnormalPodCleanC
	modPodKey := fmt.Sprintf("异常pod清理模块|%s", enableMap[modPod.Enable])
	modInfoMap[modPodKey] = keyFunc(modPod.EnabledClusters)

	modNodeDown := gr.Cg.NodeDownC
	modNodeDownKey := fmt.Sprintf("节点宕机模块|%s", enableMap[modNodeDown.Enable])
	modInfoMap[modNodeDownKey] = keyFunc(modNodeDown.EnabledClusters)

	for k, v := range modInfoMap {
		tmp := strings.Split(k, "|")
		modName, enable := tmp[0], tmp[1]
		mc := prometheus.MustNewConstMetric(ModuleInfo,
			prometheus.GaugeValue, 1,
			modName,
			enable,
			v,
		)
		ch <- mc
	}

	for clusterName, api := range gr.Cg.K8sComputePromAddrMap {
		mc := prometheus.MustNewConstMetric(ClusterInfo,
			prometheus.GaugeValue, 1,
			clusterName,
			api,
		)
		ch <- mc
	}

	nodeFunc := func() {

		klog.Info("[Guarder.Metric.Collect.models.StartWith.AllNodeMaintenances]")
		allNodes, err := model.GetAllNodeMaintenances()
		if err != nil {
			klog.Errorf("[Guarder.Metric.Collect.models.GetAllNodeMaintenances.err][err:%v]", err)
			return
		}
		if len(allNodes) == 0 {
			klog.Infof("[Guarder.Metric.Collect.models.GetAllNodeMaintenances.zero][err:%v]", err)
			return
		}

		for _, n := range allNodes {
			n := n
			cTimeStr := n.CreateTime.Local().Format("2006-01-02 15:04:05")
			uTimeStr := n.UpdateTime.Local().Format("2006-01-02 15:04:05")
			mc := prometheus.MustNewConstMetric(NodeCordonInfo,
				prometheus.GaugeValue, 1,
				n.ClusterName,
				n.Name,
				n.ModuleName,
				n.Reason,
				n.FirstDate,
				cTimeStr,
				uTimeStr,
				n.Ip,
				n.Sn,
				n.RepairTicketUrl,
				fmt.Sprintf("%d", n.HasRestart),
				fmt.Sprintf("%d", n.HasRepair),
				fmt.Sprintf("%d", n.HasRecovery),
			)
			ch <- mc

		}

	}

	podFunc := func() {
		klog.Info("[Guarder.Metric.Collect.models.StartWith.AllPodMaintenances]")
		allPods, err := model.GetAllPodMaintenances()
		if err != nil {
			klog.Errorf("[Guarder.Metric.Collect.models.GetAllPodMaintenances.err][err:%v]", err)
			return
		}
		if len(allPods) == 0 {
			klog.Infof("[Guarder.Metric.Collect.models.GetAllPodMaintenances.zero][err:%v]", err)
			return
		}

		for _, p := range allPods {
			p := p
			cTimeStr := p.CreateTime.Local().Format("2006-01-02 15:04:05")
			uTimeStr := p.UpdateTime.Local().Format("2006-01-02 15:04:05")
			mc := prometheus.MustNewConstMetric(PodDeleteInfo,
				prometheus.GaugeValue, 1,

				p.ClusterName,
				p.NodeName,
				p.Name,
				p.PodNs,
				p.ModuleName,
				p.Reason,
				p.FirstDate,
				cTimeStr,
				uTimeStr,
			)
			ch <- mc

		}

	}
	nodeFunc()
	podFunc()

}
func (gr *Guarder) Describe(ch chan<- *prometheus.Desc) {
	klog.Infof("[Guarder.Metric.Collect.Describe.start]")
	ch <- ModuleInfo
	ch <- ClusterInfo
	ch <- NodeCordonInfo
	ch <- PodDeleteInfo

}
