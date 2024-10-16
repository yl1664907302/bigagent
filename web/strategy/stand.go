package strategy

import (
	model "bigagent/model/machine"
	"bigagent/web/request"
)

type StandardStrategy struct {
	H string
}

func (s *StandardStrategy) Push() error {
	_, err := request.NewPostStand(s.H).Do()
	return err
}

func (s *StandardStrategy) Api(key string) (interface{}, error) {
	switch key {
	case "bigagent":
		return model.NewStandDataApi().ToString(), nil
	default:
		return nil, nil
	}
}
