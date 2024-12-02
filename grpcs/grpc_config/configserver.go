package grpc_config

import (
	utils "bigagent/util"
	"context"
)

type GrpcConfigServer struct {
	UnimplementedAgentConfigServiceServer
}

func (g *GrpcConfigServer) PushAgentConfig(ctx context.Context, req *AgentConfig) (*ResponseMessage, error) {
	utils.DefaultLogger.Infof("收到配置信息：%v", req)
	switch req.DataName {
	case "grpc_cmdb1_stand1":
		err := utils.ModifyYAML("config.yml", "system.grpc_cmdb1_stand1", req.NetworkInfo.Host)
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
	case "grpc_cmdb2_stand1":

	default:

	}
	return &ResponseMessage{
		Code:    "200",
		Message: "config update success",
	}, nil
	return nil, nil
}
