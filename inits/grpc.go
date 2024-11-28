package inits

import (
	"bigagent/config/global"
	"bigagent/grpcs/grpc_config"
	"google.golang.org/grpc"
	"net"
)

func RunG() {
	go func() {
		s := grpc.NewServer()
		server := grpc_config.GrpcConfigServer{}
		grpc_config.RegisterAgentConfigServiceServer(s, &server)
		// 启动服务
		lis, err := net.Listen("tcp", global.CONF.System.Grpc)
		if err != nil {
			panic(err)
		}
		err = s.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
}
