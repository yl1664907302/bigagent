package result

import "fmt"

type Results struct {
	Base
	CheckName string
	Severity  string
	Count     int
	Extra     map[string]interface{}
}

func NewResults(resultBase Base, extra map[string]interface{}, checkName string, severity string, count int, items interface{}) *Results {
	return &Results{Base: resultBase, Extra: extra, CheckName: checkName, Severity: severity, Count: count}
}

func (p *Results) SetScore() (int, error) {
	// 开始诊断,后面会优化复杂的打分算法
	if p.Count > 0 {
		switch p.Severity {
		case "info":
			return 1, nil
		case "warn":
			return 2, nil
		case "critical":
			return 3, nil
		default:
			return 0, nil
		}
	}
	return 0, nil
}

func (p *Results) GetItem() (interface{}, error) {
	return p.Items, nil
}

func (p *Results) GetCheckName() (string, error) {
	if p.CheckName == "" {
		return p.CheckName, fmt.Errorf("checkName is empty")
	}
	return p.CheckName, nil
}
