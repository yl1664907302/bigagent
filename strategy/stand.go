package strategy

import (
	"bigagent/grpcs"
	model "bigagent/model/machine"
	"fmt"
)

type StandardStrategy struct {
	K string
	G string
}

func (s *StandardStrategy) Push() error {
	//_, err := request.NewPostStand(s.H).Do()
	conn, err := grpcs.InitClient(s.G, s.K)
	if err == nil {
		go grpcs.GrpcStandPush(conn)
	}
	if s == nil {
		return fmt.Errorf("strategy is nil")
	}
	return err
}

func (s *StandardStrategy) Api(key string) (interface{}, error) {
	switch key {
	case "showdata":
		return map[string]interface{}{
			"code": 0,
			"data": model.NewSmpData(),
		}, nil
	case "patroldata":
		return map[string]interface{}{
			"code": 0,
			"data": model.NewPatrolData(),
		}, nil
	default:
		return map[string]interface{}{
			"code": 500,
			"data": "",
		}, nil
	}
}
