package check

import (
	"bigagent/internal/kubernetes"
	result2 "bigagent/internal/result"
	"context"
)

type AbnormalNode struct {
	k         func() *kubernetes.DefaultK8sOperator `json:"-"`
	ctx       context.Context                       `json:"-"`
	Cluster   string                                `json:"cluster"`
	Namespace string                                `json:"namespace"`
	Node      string                                `json:"node"`
	Message   string                                `json:"message"`
}

func NewAbnormalNode(ctx context.Context, k func() *kubernetes.DefaultK8sOperator, cluster string) *AbnormalNode {
	return &AbnormalNode{k: k, ctx: ctx, Cluster: cluster}
}

func (n *AbnormalNode) Check() result2.Result {
	return result2.NewResultNode(result2.Base{Cluster: n.Cluster, Items: nil}, nil, "NodeDown", "info", 0, nil)
}
