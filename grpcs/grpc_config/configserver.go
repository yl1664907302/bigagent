package grpc_config

import (
	"bigagent/config/global"
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
			return &ResponseMessage{
				Code:    "500",
				Message: "config update failed",
			}, err
		}
		global.ACTION_DETAIL = "当前配置[" + req.Id + "]"
	case "stand2":
	default:
	}
	return &ResponseMessage{
		Code:    "200",
		Message: "config update success",
	}, nil
}
