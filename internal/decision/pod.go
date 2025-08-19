package decision

import (
	"bigagent/internal/check"
	"bigagent/internal/check/result"
	utils "bigagent/internal/utils"
)

type AbnormalPodDecision struct {
	result result.Result
}

func NewAbnormalPodDecision(res result.Result) *AbnormalPodDecision {
	return &AbnormalPodDecision{result: res}
}

func (p *AbnormalPodDecision) recovery(key bool) error {
	if key {
		utils.DefaultLogger.Warnf("正在处理pod异常")
		item, err := p.result.GetItem()
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

func (p *AbnormalPodDecision) recoveryTem(key bool) error {
	if key {
		utils.DefaultLogger.Warnf("正在处理pod异常")
		item, err := p.result.GetItem()
		if err != nil {
			return err
		}
		pod, ok := item.(check.AbnormalPod)
		if !ok {
			return nil
		}
		// 检查 Pod 是否仍处于终止状态，默认无操作
		if pod.Terminating {
			utils.DefaultLogger.Warnf("正在删除pod%s", pod.Pod)
		} else {
			utils.DefaultLogger.Warnf("Pod %s 失效，无效处理", pod.Pod)
		}
	}
	return nil
}

func (p *AbnormalPodDecision) Judge() func() {
	score, err := p.result.SetScore()
	if err != nil {
		utils.DefaultLogger.Errorf("设置分数失败pod：%s", err.Error())
		return nil
	}
	if score > 2 {
		return func() {
			utils.DefaultLogger.Warnf("开始处理pod异常")
			if err := p.recovery(true); err != nil {
				utils.DefaultLogger.Errorf("失败处理pod异常：%s", err.Error())
			} else {
				utils.DefaultLogger.Warnf("完成处理pod异常")
			}
		}
	}
	return nil
}

func (p *AbnormalPodDecision) JudgeTem() func() {
	score, err := p.result.SetScore()
	if err != nil {
		utils.DefaultLogger.Errorf("设置分数失败pod：%s", err.Error())
		return nil
	}
	if score > 2 {
		return func() {
			utils.DefaultLogger.Warnf("开始处理pod异常")
			if err := p.recoveryTem(true); err != nil {
				utils.DefaultLogger.Errorf("失败处理pod异常：%s", err.Error())
			} else {
				utils.DefaultLogger.Warnf("完成处理pod异常")
			}
		}
	}
	return nil
}
