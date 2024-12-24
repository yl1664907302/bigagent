package strategy

import (
	"bigagent/grpcs"
	model "bigagent/model/machine"
	"fmt"
)

type StandardStrategy struct {
	H string
	G string
}

func (s *StandardStrategy) Push() error {
	//_, err := request.NewPostStand(s.H).Do()
	conn, err := grpcs.InitClient(s.G)
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
	default:
		return map[string]interface{}{
			"code": 500,
			"data": "",
		}, nil
	}
}
