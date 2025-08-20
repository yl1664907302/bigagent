package main

import (
	"bigagent/inits"
	"bigagent/internal/check/pod"
	"bigagent/internal/config/global"
	"bigagent/internal/decision/pod"
	"bigagent/internal/kubernetes"
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
	inits.InitOsqueryClient()
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
			utils.DefaultLogger.Info("没有需要处理的 Pod 异常")
		} else {
			recovery()
		}
	}()

	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
