package result

// Result 定义检查统一返回结构，供 decision 层消费
// Items 可承载具体检查项的明细切片，例如 []AbnormalPod
// Severity 建议值："info" | "warn" | "critical"

type Result interface {
	SetScore() (int, error)
	GetItem() (interface{}, error)
}

type Base struct {
	Cluster string
	Items   interface{}
}

func (r *Base) SetScore() (int, error) {
	return 0, nil
}

func (r *Base) GetItem() (interface{}, error) {
	return r.Items, nil
}
