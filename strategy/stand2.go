package strategy

import (
	"bigagent/grpcs"
	model "bigagent/model/machine"
	"fmt"
)

type StandardStrategy2 struct {
	H string
	G string
}

func (s *StandardStrategy2) Push() error {
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

func (s *StandardStrategy2) Api(key string) (interface{}, error) {
	switch key {
	case "showdata":
		return map[string]interface{}{
			"code": 0,
			"data": model.NewStandData(),
		}, nil
	default:
		return map[string]interface{}{
			"code": 500,
			"data": "",
		}, nil
	}
}
