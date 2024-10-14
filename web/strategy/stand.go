package strategy

import model "bigagent/model/machine"

type StandardStrategy struct {
}

func (s *StandardStrategy) Push() (interface{}, error) {

	return nil, nil
}

func (s *StandardStrategy) Api(key string) (interface{}, error) {
	switch key {
	case "bigagent":
		return model.NewStandData().ToString(), nil
	default:
		return nil, nil
	}
}
