package strategy

import (
	grpc_client "bigagent/grpcs/client"
	model "bigagent/model/machine"
	"fmt"
)

type StandardStrategy struct {
	H string
	G string
}

func (s *StandardStrategy) Push() error {
	//_, err := request.NewPostStand(s.H).Do()
	conn, err := grpc_client.InitClient(s.G)
	if err == nil {
		grpc_client.GrpcStandPush(conn)
	}
	if s == nil {
		return fmt.Errorf("strategy is nil")
	}
	return err
}

func (s *StandardStrategy) Api(key string) (interface{}, error) {
	switch key {
	case "bigagent":
		return model.NewStandDataApi(), nil
	default:
		return nil, nil
	}
}
