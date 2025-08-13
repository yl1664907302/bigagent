package main

import (
	"bigagent/inits"
	"bigagent/internal/check"
	"bigagent/internal/config/global"
	"bigagent/internal/decision"
	"bigagent/internal/kubernetes"
	"bigagent/internal/polling"
	utils "bigagent/internal/util"
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
	// 测试代码
	func() {
		// 生成参数
		ctx := context.Background()
		clusters := []kubernetes.KubeCluster{
			{
				Name:       "test",
				Kubeconfig: "kubeconfig.yaml",
			},
		}

		result, _ := polling.NewKubePolling(check.NewAbnormalPod(ctx, func() *kubernetes.DefaultK8sOperator {
			return &kubernetes.DefaultK8sOperator{
				Clusters: clusters,
			}
		}, "test", "")).P.Check()

		num, err := decision.NewAbnormalPodDecision(result).Judge()
		if err != nil {
			utils.DefaultLogger.Errorf("判断异常 Pod 时发生错误: %v", err)
			return
		}
		utils.DefaultLogger.Info("pod诊断分数为：", num)

	}()

	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
