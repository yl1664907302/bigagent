package main

import (
	"bigagent/inits"
	check2 "bigagent/internal/check/node"
	check "bigagent/internal/check/pod"
	"bigagent/internal/config/global"
	decision2 "bigagent/internal/decision/node"
	"bigagent/internal/decision/pod"
	"bigagent/internal/kubernetes"
	prom "bigagent/internal/prometheus"
	"bigagent/internal/utils"
	"context"
	"os"
)

func init() {
	// 创建一个通道来接收系统信号
	sigs := make(chan os.Signal, 1)
	inits.InitCmd(sigs)
	inits.LoggerInit()
	inits.InstallIfNotExists([]string{"git", "wget", "curl", "gcc", "make", "jq", "https://pkg.osquery.io/rpm/osquery-5.11.0-1.linux.x86_64.rpm"})
	inits.AgentRegister()
	//inits.InitOsqueryClient()
	inits.Crontab()
	inits.ListerChannel()
	utils.DefaultLogger.Info("当前代码版本为：", "20250730")
}

func main() {
	// 测试异常pod清理
	func() {
		// 生成参数
		ctx := context.Background()
		clusters := []kubernetes.KubeCluster{
			{
				Name:       "test",
				Kubeconfig: "kubeconfig1.yaml",
			},
		}

		//  创建一个轮询链路，妈的写的潦草
		recovery := decision.NewAbnormalPodDecision(check.NewAbnormalPod(ctx, func() *kubernetes.DefaultK8sOperator {
			return &kubernetes.DefaultK8sOperator{
				Clusters: clusters,
			}
		}, "test", "").CheckPodNeedDelete("test")).JudgeNeedDelete()
		if recovery == nil {
			utils.DefaultLogger.Info("无法获取自愈器，或没有需要处理的 Pod 异常")
		} else {
			recovery()
		}
	}()

	// 测试异常node清理
	func() {
		// 生成参数
		ctx := context.Background()
		clusters := []kubernetes.KubeCluster{
			{
				Name:       "test",
				Kubeconfig: "kubeconfig1.yaml",
			},
		}
		// 装配Prometheus客户端
		var proms []*prom.PromClient
		promCmap := global.V.GetStringMapString("k8s_compute_prom_addr_map")
		for _, v := range promCmap {
			c, err := prom.NewPromClient(v, global.V.GetInt("recovery_conf.query_prom_time_out_seconds"))
			if err != nil {
				utils.DefaultLogger.Errorf("NewPromClient err: %v", err)
			}
			proms = append(proms, c)
		}

		for _, c := range proms {
			//  创建一个轮询链路，妈的写的潦草
			recovery := decision2.NewAbnormalNodeDecision(check2.NewAbnormalNode(ctx, func() *kubernetes.DefaultK8sOperator {
				return &kubernetes.DefaultK8sOperator{
					Clusters: clusters,
				}
			}, "test", c).CheckNodeByProm(global.V.GetString("check_ql_map.file_system_read_only"))).JudgeDeadNode()
			if recovery == nil {
				utils.DefaultLogger.Info("无法获取自愈器，或没有需要处理的 node 异常")
			} else {
				recovery()
			}

			//  创建一个轮询链路，妈的写的潦草
			recovery2 := decision2.NewAbnormalNodeDecision(check2.NewAbnormalNode(ctx, func() *kubernetes.DefaultK8sOperator {
				return &kubernetes.DefaultK8sOperator{
					Clusters: clusters,
				}
			}, "test", c).CheckNodeByProm(global.V.GetString("check_ql_map.arp_too_many"))).JudgeDeadNode()
			if recovery == nil {
				utils.DefaultLogger.Info("无法获取自愈器，或没有需要处理的 node 异常")
			} else {
				recovery2()
			}
		}

	}()

	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
