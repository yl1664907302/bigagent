package check

import (
	"bigagent/internal/check/result"
	"bigagent/internal/kubernetes"
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

func (n *AbnormalNode) Check() result.Result {
	return result.NewResultNode(result.Base{Cluster: n.Cluster, Items: nil}, nil, "NodeDown", "info", 0, nil)
}
