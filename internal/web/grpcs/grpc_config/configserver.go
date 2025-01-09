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
		err := utils2.ModifyYAML("config.yml", "action_detail", global.V.GetString("action_detail")+" failed["+req.Id+"]")
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
		if req.NetworkInfo.Host == "" {
			previous := utils2.RemoveStringAndPrevious(global.V.GetString("action_detail"), req.Id)
			err = utils2.ModifyYAML("config.yml", "action_detail", previous)
			//	如果该配置从才来没发过
		} else if !utils2.ContainsReqId(global.V.GetString("action_detail"), req.Id) {
			//  也没失败过，直接追加
			if !utils2.ContainsFailedReqId(global.V.GetString("action_detail"), req.Id) {
				err = utils2.ModifyYAML("config.yml", "action_detail", global.V.GetString("action_detail")+" "+req.Id)
			} else {
				//  如果失败过，删除失败记录
				previous := utils2.RemoveStringAndPrevious(global.V.GetString("action_detail"), req.Id)
				err = utils2.ModifyYAML("config.yml", "action_detail", previous)
			}
		}
		if err != nil {
			utils2.DefaultLogger.Error(err)
			err = utils2.ModifyYAML("config.yml", "action_detail", global.V.GetString("action_detail")+" failed["+req.Id+"]")
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
		err = utils2.ModifyYAML("config.yml", "action_detail", global.V.GetString("action_detail")+" failed["+req.Id+"]")
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
