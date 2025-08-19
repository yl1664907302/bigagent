package main

import (
	"bigagent/inits"
	"bigagent/internal/check"
	"bigagent/internal/config/global"
	"bigagent/internal/decision"
	"bigagent/internal/kubernetes"
	"bigagent/internal/utils"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"time"
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
				Kubeconfig: "kubeconfig1.yaml",
			},
		}

		// 生成一个 k8s 操作对象mock数据
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "test-pod",
				Namespace:         "default",
				UID:               "12345",                        // 随便填一个
				DeletionTimestamp: &metav1.Time{Time: time.Now()}, // 模拟正在 Terminating
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		}

		//  创建一个轮询链路，妈的写的潦草
		recovery := decision.NewAbnormalPodDecision(check.NewAbnormalPod(ctx, func() *kubernetes.DefaultK8sOperator {
			return &kubernetes.DefaultK8sOperator{
				Clusters: clusters,
			}
		}, "test", "").CheckPodTerminating(pod)).JudgeTem()
		if recovery == nil {
			utils.DefaultLogger.Warnf("pod：%s 不存在，或者不处于Terminating状态", pod.Name)
		} else {
			recovery()
			utils.DefaultLogger.Warn("Pod 终止失败，可能需要手动干预")
		}
	}()

	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
