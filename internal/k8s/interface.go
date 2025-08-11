package k8s

type Service interface {
	ClusterInfo() (interface{}, error)
}
