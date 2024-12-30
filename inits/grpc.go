package inits

import (
	"bigagent/internal/config/global"
	grpc_config2 "bigagent/internal/web/grpcs/grpc_config"
	"net"

	"google.golang.org/grpc"
)

func RunG() {
	go func() {
		s := grpc.NewServer()

		// 注册服务端
		server := grpc_config2.GrpcConfigServer{}
		grpc_config2.RegisterAgentConfigServiceServer(s, &server)

		// 启动服务
		lis, err := net.Listen("tcp", global.V.GetString("system.grpc")+":"+global.V.GetString("system.grpc_port"))
		if err != nil {
			panic(err)
		}
		err = s.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
}
