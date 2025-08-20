package decision

import (
	"bigagent/internal/check/node"
	"bigagent/internal/result"
	"bigagent/internal/utils"
)

type AbnormalNodeDecision struct {
	Result result.Result
}

func (p *AbnormalNodeDecision) Recovery(key bool) error {
	if key {
		utils.DefaultLogger.Warnf("正在处理节点异常")
		item, err := p.Result.GetItem()
		if err != nil {
			return err
		}
		nodes, ok := item.([]node.NodeDownItem)
		if !ok {
			return nil
		}

		for _, nd := range nodes {
			utils.DefaultLogger.Warnf("处理节点异常: 集群 %s, 节点 %s, IP %s", nd.Cluster, nd.Node, nd.IP)
		}
	}
	return nil
}

func (p *AbnormalNodeDecision) Judge() func() {
	score, err := p.Result.SetScore()
	if err != nil {
		utils.DefaultLogger.Errorf("设置分数失败node：%s", err.Error())
		return nil
	}
	if score > 2 {
		return func() {
			utils.DefaultLogger.Warnf("开始处理节点异常")
			if err := p.Recovery(true); err != nil {
				utils.DefaultLogger.Errorf("失败处理节点异常：%s", err.Error())
			} else {
				utils.DefaultLogger.Warnf("完成处理节点异常")
			}
		}
	}
	return nil
}
