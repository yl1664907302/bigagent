package check

import (
	"bigagent/internal/kubernetes"
	"context"
)

type AbnormalNode struct {
	k         func() *kubernetes.DefaultK8sOperator `json:"-"`
	ctx       context.Context                       `json:"-"`
	Cluster   string                                `json:"cluster"`
	Namespace string                                `json:"namespace"`
	Pod       string                                `json:"pod"`
	Message   string                                `json:"message"`
}
