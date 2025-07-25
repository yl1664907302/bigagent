package strategy

import (
	"bigagent/internal/web/request"
)

type VeopsStrategy struct {
	H        string
	KeyFirst bool
}

// Push 当前的处理方式不当，结果被抛弃
func (s *VeopsStrategy) Push() error {
	veops := request.NewPostVeops(s.H, true)
	_, err := veops.Do()
	//super
	s.KeyFirst = veops.KeyFirst
	return err
}

func (s *VeopsStrategy) Api(key string) (interface{}, error) {
	switch key {
	case "xxx":
		return "", nil
	default:
		return nil, nil
	}
}
