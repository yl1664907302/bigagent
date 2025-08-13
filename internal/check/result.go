package check

// Result 定义检查统一返回结构，供 decision 层消费
// Items 可承载具体检查项的明细切片，例如 []AbnormalPod
// Severity 建议值："info" | "warn" | "critical"
type Result struct {
	CheckName string                 // 检查名称，如 PodAbnormalRestarts
	Cluster   string                 // 集群
	Namespace string                 // 命名空间（可空代表全量）
	Severity  string                 // 严重级别
	Count     int                    // 命中数量
	Items     interface{}            // 明细
	Extra     map[string]interface{} // 额外元数据
}
