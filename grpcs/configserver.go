package grpcs

import (
	grpc_server "bigagent/grpcs/config"
	"context"
)

type GrpcConfigServer struct {
	grpc_server.UnimplementedAgentConfigServiceServer
}

func (g *GrpcConfigServer) PushAgentConfig(ctx context.Context, req *grpc_server.AgentConfig) (*grpc_server.ResponseMessage, error) {
	// 1. 从请求中获取配置信息
	// 2. 将配置信息存储到数据库
	// 3. 返回响应
	return &grpc_server.ResponseMessage{
		Code:    "200",
		Message: "配置信息推送成功",
	}, nil

}
