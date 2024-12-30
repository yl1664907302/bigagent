package strategy

import (
	model2 "bigagent/internal/model/machine"
	"bigagent/internal/web/grpcs"
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
		return nil
	}
	return err
}

func (s *StandardStrategy) Api(key string) (interface{}, error) {
	switch key {
	case "showdata":
		return map[string]interface{}{
			"code": 0,
			"data": model2.NewSmpData(),
		}, nil
	case "patroldata":
		return map[string]interface{}{
			"code": 0,
			"data": model2.NewPatrolData(),
		}, nil
	default:
		return map[string]interface{}{
			"code": 500,
			"data": "",
		}, nil
	}
}
