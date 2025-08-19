package result

type NodeResult struct {
	Base
	CheckName string
	Severity  string
	Count     int
	Extra     map[string]interface{}
}

func NewResultNode(resultBase Base, extra map[string]interface{}, checkName string, severity string, count int, items interface{}) *NodeResult {
	return &NodeResult{Base: resultBase, Extra: extra, CheckName: checkName, Severity: severity, Count: count}
}

func (n *NodeResult) SetScore() (int, error) {
	if n.Count > 0 {
		switch n.Severity {
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

func (n *NodeResult) GetItem() (interface{}, error) {
	return n.Items, nil
}
