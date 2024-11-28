package grpc_config

import (
	"bigagent/util/logger"
	"context"
	"fmt"
)

type GrpcConfigServer struct {
	UnimplementedAgentConfigServiceServer
}

func (g *GrpcConfigServer) PushAgentConfig(ctx context.Context, req *AgentConfig) (*ResponseMessage, error) {
	logger.DefaultLogger.Infof("收到配置信息：%v", req)
	//处理字段映射
	// 这里可以将配置信息写入到配置文件中
	fmt.Println(req.FieldMapping)
	return &ResponseMessage{
		Code:    "200",
		Message: "config update success",
	}, nil

}
