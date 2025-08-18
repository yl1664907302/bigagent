package decision

import (
	"bigagent/internal/check"
	"bigagent/internal/check/result"
	"bigagent/internal/utils"
)

type AbnormalNodeDecision struct {
	Result result.Result
}

func (p *AbnormalNodeDecision) Recovery(key bool) error {
	if key {
		utils.DefaultLogger.Warnf("正在处理pod异常")
		item, err := p.Result.GetItem()
		if err != nil {
			return err
		}
		pods, ok := item.([]check.AbnormalPod)
		if !ok {
			return nil
		}

		for _, pod := range pods {
			utils.DefaultLogger.Warnf("处理pod异常: 集群 %s, 命名空间 %s, Pod %s, 容器 %s, 重启次数 %d, 原因 %s, 信息 %s",
				pod.Cluster, pod.Namespace, pod.Pod, pod.Container, pod.RestartCount, pod.Reason, pod.Message)
		}
	}
	return nil
}

func (p *AbnormalNodeDecision) Judge() func() {
	score, err := p.Result.SetScore()
	if err != nil {
		utils.DefaultLogger.Errorf("设置分数失败pod：%s", err.Error())
		return nil
	}
	if score > 2 {
		return func() {
			utils.DefaultLogger.Warnf("开始处理pod异常")
			if err := p.Recovery(true); err != nil {
				utils.DefaultLogger.Errorf("失败处理pod异常：%s", err.Error())
			} else {
				utils.DefaultLogger.Warnf("完成处理pod异常")
			}
		}
	}
	return nil
}
