package guard

import (
	"bigagent/internal/config/global"
	clientset "k8s.io/client-go/kubernetes"

	"k8s.io/klog/v2"
)

type Guarder struct {

	// hold 所有被管理k8s集群的 map ，item是k8s 的client-set
	ClientSetMap map[string]*clientset.Clientset
	Cg           *global.Config
}

func NewGuarder(cg *global.Config) *Guarder {

	gr := &Guarder{
		ClientSetMap: make(map[string]*clientset.Clientset),
		Cg:           cg,
	}
	num := len(cg.K8sConfigS)
	index := 1
	for clusterName, kConfig := range cg.K8sConfigS {

		kConfig, err := buildConfig(kConfig)
		if err != nil {
			klog.Fatalf("[NewGuarder.buildConfig.err][%d/%d][cluster:%v][err:%v]", index, num, clusterName, err)

		}
		kc := clientset.NewForConfigOrDie(kConfig)
		gr.ClientSetMap[clusterName] = kc
		klog.Infof("[NewGuarder.buildConfig.success][%d/%d][cluster:%v]", index, num, clusterName)
		index++
	}
	return gr

}
