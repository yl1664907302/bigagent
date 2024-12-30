package grpc_config

import (
	"bigagent/internal/config/global"
	utils2 "bigagent/internal/util"
	"context"
	"google.golang.org/grpc/metadata"
)

type GrpcConfigServer struct {
	UnimplementedAgentConfigServiceServer
}

func (g *GrpcConfigServer) PushAgentConfig(ctx context.Context, req *AgentConfig) (*ResponseMessage, error) {
	//密钥验证
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := utils2.ModifyYAML("config.yml", "action_detail", "config update failed["+req.Id+"]")
		if err != nil {
			utils2.DefaultLogger.Error(err)
		}
		return &ResponseMessage{
			Code:    "500",
			Message: "agent_config serct is error ！",
		}, nil
	}
	tokens := md["authorization"]
	validTokenFound := false
	for _, token := range tokens {
		// 检查每个token是否是期望的Token
		if token == "Bearer "+global.V.GetString("system.serct") {
			validTokenFound = true
			break
		}
	}
	if !validTokenFound {
		err := utils2.ModifyYAML("config.yml", "action_detail", "config update failed["+req.Id+"]")
		if err != nil {
			utils2.DefaultLogger.Error(err)
			return &ResponseMessage{
				Code:    "500",
				Message: "agent_config Authorization token is missing！",
			}, nil
		}
	}
	switch req.DataName {
	case "stand1":
		err := utils2.ModifyYAML("config.yml", "grpc_cmdb"+req.SlotName+"_stand1", req.NetworkInfo.Host)
		err = utils2.ModifyYAML("config.yml", "grpc_cmdb"+req.SlotName+"_stand1"+"_token", req.Token)
		err = utils2.ModifyYAML("config.yml", "action_detail", req.Id)
		if err != nil {
			utils2.DefaultLogger.Error(err)
			err = utils2.ModifyYAML("config.yml", "action_detail", "config update failed["+req.Id+"]")
			return &ResponseMessage{
				Code:    "500",
				Message: "config update failed",
			}, err
		}
	case "stand2":
	default:
	}
	err := utils2.ModifyYAML("config.yml", "collection_frequency", req.CollectionFrequency)
	if err != nil {
		utils2.DefaultLogger.Error(err)
		err = utils2.ModifyYAML("config.yml", "action_detail", "config update failed["+req.Id+"]")
		return &ResponseMessage{
			Code:    "500",
			Message: "config update failed",
		}, err
	}
	return &ResponseMessage{
		Code:    "200",
		Message: "config update success",
	}, nil
}
