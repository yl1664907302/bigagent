package k8s

// DefaultService 是一个占位实现，后续可替换为使用 client-go 的真实实现

type DefaultService struct {
	Address string
}

func NewDefaultService(address string) Service {
	return &DefaultService{Address: address}
}

func (s *DefaultService) ClusterInfo() (interface{}, error) {
	return map[string]interface{}{
		"address": s.Address,
		"status":  "ok",
	}, nil
}
