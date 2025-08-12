package main

import (
	"bigagent/inits"
	"bigagent/internal/config/global"
	"bigagent/internal/decision"
	"bigagent/internal/kubernetes"
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

		k := decision.NewKubeDecision(decision.NewAbnormalPod(func() *kubernetes.DefaultK8sOperator {
			return &kubernetes.DefaultK8sOperator{
				Clusters: clusters,
			}
		}, ctx, "test", "nantong-20"))

		err := k.C.StartCheck()
		if err != nil {
			utils.DefaultLogger.Error("Kubernetes集群检查失败:", err)
			return
		}
	}()

	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
