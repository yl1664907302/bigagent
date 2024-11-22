package strategy

import (
	"bigagent/web/request"
)

type VeopsStrategy struct {
	H string
}

// Push 当前的处理方式不当，结果被抛弃
func (s *VeopsStrategy) Push() error {
	_, err := request.NewPostVeops(s.H).Do()
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
