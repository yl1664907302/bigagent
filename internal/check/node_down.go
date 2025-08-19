package check

import (
	"bigagent/internal/check/result"
	"bigagent/internal/utils"
	"context"
	"fmt"
)

type NodeDownCheck struct {
	ctx           context.Context
	Cluster       string
	PromAddr      string
	CheckQlList   []string
	NodeNameToIpQ string
}

type NodeDownItem struct {
	Cluster string `json:"cluster"`
	Node    string `json:"node"`
	IP      string `json:"ip"`
}

func NewNodeDownCheck(ctx context.Context, cluster, promAddr string, checkQls []string, nodeNameToIpQl string) *NodeDownCheck {
	return &NodeDownCheck{ctx: ctx, Cluster: cluster, PromAddr: promAddr, CheckQlList: checkQls, NodeNameToIpQ: nodeNameToIpQl}
}

func (c *NodeDownCheck) Check() result.Result {
	// 统计每个节点命中规则的次数
	nodeHit := map[string]int{}
	for _, ql := range c.CheckQlList {
		vec, err := utils.PromInstantQuery(c.ctx, ql, c.PromAddr, 10)
		if err != nil {
			continue
		}
		for _, s := range vec {
			n := string(s.Metric["node"])
			if n == "" {
				continue
			}
			nodeHit[n]++
		}
	}
	threshold := len(c.CheckQlList)
	var items []NodeDownItem
	for node, num := range nodeHit {
		if num < threshold {
			continue
		}
		// 查询IP
		ip := ""
		if c.NodeNameToIpQ != "" {
			vec, err := utils.PromInstantQuery(c.ctx, formatNodeNameQl(c.NodeNameToIpQ, node), c.PromAddr, 10)
			if err == nil {
				for _, v := range vec {
					ip = string(v.Metric["instance"])
					break
				}
			}
		}
		items = append(items, NodeDownItem{Cluster: c.Cluster, Node: node, IP: ip})
	}
	severity := "info"
	switch len(items) {
	case 0:
		severity = "info"
	case 1:
		severity = "warn"
	default:
		severity = "critical"
	}
	return result.NewResultNode(result.Base{Cluster: c.Cluster, Items: items}, nil, "NodeDown", severity, len(items), items)
}

func formatNodeNameQl(tpl, node string) string {
	// 简单替换占位符 %s
	return fmt.Sprintf(tpl, node)
}
