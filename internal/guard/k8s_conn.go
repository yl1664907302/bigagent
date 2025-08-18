package guard

import (
	"context"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (gr *Guarder) GetK8sTwContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(gr.Cg.K8sTimeoutSeconds)*time.Second)
}
