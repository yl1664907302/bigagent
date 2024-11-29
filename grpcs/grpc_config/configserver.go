package grpc_config

import (
	"bigagent/util/logger"
	"context"
)

type GrpcConfigServer struct {
	UnimplementedAgentConfigServiceServer
}

func (g *GrpcConfigServer) PushAgentConfig(ctx context.Context, req *AgentConfig) (*ResponseMessage, error) {
	logger.DefaultLogger.Infof("收到配置信息：%v", req)
	switch req.DataName {
	case "stand1":

	case "stand2":

	default:

	}
	return &ResponseMessage{
		Code:    "200",
		Message: "config update success",
	}, nil
	return nil, nil
}
