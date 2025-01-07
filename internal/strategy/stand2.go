package strategy

import (
	"bigagent/internal/model/machine"
	utils "bigagent/internal/util"
	"bigagent/internal/web/grpcs"
)

type StandardStrategy2 struct {
	K string
	G string
}

func (s *StandardStrategy2) Push() error {
	//_, err := request.NewPostStand(s.H).Do()
	conn, err := grpcs.InitClient(s.G, s.K)
	if err != nil {
		utils.DefaultLogger.Errorf("grpc服务端连接异常: %v", err)
		return err
	}
	go grpcs.GrpcStandPush(conn)
	return err
}

func (s *StandardStrategy2) Api(key string) (interface{}, error) {
	switch key {
	case "showdata":
		return map[string]interface{}{
			"code": 0,
			"data": model.NewStandDataApi(),
		}, nil
	default:
		return map[string]interface{}{
			"code": 500,
			"data": "",
		}, nil
	}
}
