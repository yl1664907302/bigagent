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
		err = utils.ModifyYAML("config.yml", "action_detail", "当前配置["+req.Id+"]")
		if err != nil {
			utils.DefaultLogger.Error(err)
			return &ResponseMessage{
				Code:    "500",
				Message: "config update failed",
			}, err
		}
	case "stand2":
	default:
	}
	return &ResponseMessage{
		Code:    "200",
		Message: "config update success",
	}, nil
}
