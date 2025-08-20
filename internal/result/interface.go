package result

type Result interface {
	SetScore() (int, error)
	GetItem() (interface{}, error)
	GetNode2Ip() (map[string]string, error)
	GetCheckName() (string, error)
}

type Base struct {
	Cluster string
	Items   interface{}
	Node2Ip map[string]string
}

func (r *Base) SetScore() (int, error) {
	return 0, nil
}

func (r *Base) GetItem() (interface{}, error) {
	return r.Items, nil
}

func (r *Base) GetNode2Ip() (map[string]string, error) {
	return r.Node2Ip, nil
}
func (r *Base) GetCheckName() (string, error) {
	return "", nil
}
