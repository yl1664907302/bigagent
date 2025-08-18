package strategy

import (
	model2 "bigagent/internal/model/machine"
	"bigagent/internal/scrape/machine"
	utils "bigagent/internal/utils"
	"bigagent/internal/web/grpcs"
)

type StandardStrategy struct {
	K      string
	G      string
	KeyUse bool // 是否使用key
}

func (s *StandardStrategy) Push() error {
	if !s.KeyUse {
		return nil
	}
	//_, err := request.NewPostStand(s.H).Do()
	conn, err := grpcs.InitClient(s.G, s.K)
	if err != nil {
		utils.DefaultLogger.Errorf("grpc服务端连接异常: %v", err)
		return err
	}
	go grpcs.GrpcStandPush(conn)
	return err
}

func (s *StandardStrategy) Api(key string) (interface{}, error) {
	switch key {
	case "showdata":
		return map[string]interface{}{
			"code": 0,
			"data": machine.SmpMa,
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
