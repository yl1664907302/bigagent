package strategy

import (
	model "bigagent/model/machine"
	"bigagent/web/request"
)

type StandardStrategy struct {
	H string
}

func (s *StandardStrategy) Push() (interface{}, error) {
	do, err := request.NewPostStand(s.H).Do()
	return do, err
}

func (s *StandardStrategy) Api(key string) (interface{}, error) {
	switch key {
	case "bigagent":
		return model.NewStandData().ToString(), nil
	default:
		return nil, nil
	}
}
