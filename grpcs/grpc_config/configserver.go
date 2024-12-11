package grpc_config

import (
	utils "bigagent/util"
	"context"
)

type GrpcConfigServer struct {
	UnimplementedAgentConfigServiceServer
}

func (g *GrpcConfigServer) PushAgentConfig(ctx context.Context, req *AgentConfig) (*ResponseMessage, error) {

	switch req.DataName {
	case "stand1":
		err := utils.ModifyYAML("config.yml", "grpc_cmdb"+req.SlotName+"_stand1", req.NetworkInfo.Host)
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
	case "stand2":

	default:

	}
	return &ResponseMessage{
		Code:    "200",
		Message: "config update success",
	}, nil
	return nil, nil
}
