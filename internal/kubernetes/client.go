package kubernetes

import (
	"fmt"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeCluster struct {
	Name       string
	Kubeconfig string
	Namespace  string
}
type ClientCache struct {
	mu sync.Mutex
}

type DefaultK8sOperator struct {
	Clusters []KubeCluster
	cache    *ClientCache
}

func NewDefaultK8sOperator(clusters []KubeCluster) *DefaultK8sOperator {
	m := []KubeCluster{}
	return &DefaultK8sOperator{
		Clusters: m,
	}
}

func (op *DefaultK8sOperator) getClient(clusterName string) (*kubernetes.Clientset, *rest.Config, error) {
	var cluster *KubeCluster

	if len(op.Clusters) == 0 {
		return nil, nil, fmt.Errorf("cluster_config not set")
	}

	// 查找集群配置与clusterName匹配
	for _, c := range op.Clusters {
		if c.Name == clusterName {
			cluster = &c
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", cluster.Kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return clientset, config, nil
}

// Client 对外暴露，返回 clientset 与 rest.Config
func (op *DefaultK8sOperator) Client(clusterName string) (*kubernetes.Clientset, *rest.Config, error) {
	return op.getClient(clusterName)
}
