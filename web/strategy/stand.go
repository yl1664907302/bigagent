package strategy

import model "bigagent/model/machine"

type StandardStrategy struct {
}

func (s *StandardStrategy) Push() (interface{}, error) {
	return nil, nil
}

func (s *StandardStrategy) Api() (interface{}, error) {
	return model.NewStandData().ToString(), nil
}
