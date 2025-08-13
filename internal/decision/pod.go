package decision

import (
	"bigagent/internal/check"
)

type AbnormalPodDecision struct {
	Result check.Result
}

func NewAbnormalPodDecision(res check.Result) *AbnormalPodDecision {
	return &AbnormalPodDecision{Result: res}
}

func (p *AbnormalPodDecision) Recovery() error {
	// 这里可以根据 p.Result 执行自动化恢复动作（例如触发重启/告警），暂留空实现
	return nil
}

// Judge 根据统一 Result 输出判定分值（示例：Count>0 判为 1，否则 0）
func (p *AbnormalPodDecision) Judge() (interface{}, error) {

	if p.Result.Count > 0 {
		return 1, nil
	}
	return 0, nil
}
