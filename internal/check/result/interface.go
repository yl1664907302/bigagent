package result

type Result interface {
	SetScore() (int, error)
	GetItem() (interface{}, error)
	GetMetric() (map[string]interface{}, error)
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

func (r *Base) GetMetric() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
