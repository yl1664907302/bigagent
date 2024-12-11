package grpcs

import (
	"bigagent/config/global"
	grpc_server "bigagent/grpcs/server"
	model "bigagent/model/machine"
	utils "bigagent/util"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	client grpc_server.PushAgantDataClient
)

// InitClient 通用grpc客户端
func InitClient(host string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond) //连接超时设置为1000毫秒
	defer cancel()
	//连接到服务端
	conn, err := grpc.DialContext(
		ctx,
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()), //Credential即使为空，也必须设置
		grpc.WithBlock(), //grpc.WithBlock()直到连接真正建立才会返回，否则连接是异步建立的。因此grpc.WithBlock()和Timeout结合使用才有意义。server端正常的情况下使用grpc.WithBlock()得到的connection.GetState()为READY，不使用grpc.WithBlock()得到的connection.GetState()为IDEL
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(10<<20), grpc.MaxCallRecvMsgSize(10<<20)), //默认情况下SendMsg上限是MaxInt32，RecvMsg上限是4M，这里都修改为10M
	)
	if err != nil {
		return nil, fmt.Errorf("服务端'%s'连接超时!", host)
	}
	return conn, err
}

// GrpcStandPush 执行GRPC标准数据类型推送方法
func GrpcStandPush(conn *grpc.ClientConn) {
	client = grpc_server.NewPushAgantDataClient(conn)
	//准备好请求参数
	data := model.NewSmpDataGrpc()
	request := grpc_server.SmpData{
		Serct:    global.V.GetString("system.serct"),
		Uuid:     data.Uuid,
		Hostname: data.Hostname,
		Ipv4:     data.IPv4,
		GrpcPort: global.V.GetString("system.grpc_port"),
		Time:     timestamppb.New(data.Time),
		Cpu: &grpc_server.SmpCpu{
			Name:  data.Cpu.Name,
			Core:  data.Cpu.Core,
			Usage: data.Cpu.Usage,
		},
		Disk: data.Disk,
		Memory: &grpc_server.SmpMemory{
			VirtualMemory: &grpc_server.VirtualMemory{
				Total:       data.Memory.Vmem.Total,
				Used:        data.Memory.Vmem.Used,
				Free:        data.Memory.Vmem.Free,
				UsedPercent: data.Memory.Vmem.UsedPercent,
			},
			SwapMemory: &grpc_server.SwapMemory{
				Total:       data.Memory.Swap.Total,
				Used:        data.Memory.Swap.Used,
				Free:        data.Memory.Swap.Free,
				UsedPercent: data.Memory.Swap.UsedPercent,
			},
		},
		Kmodules: data.Kmodules,
		Smpnet:   data.Net,
		Smpps:    data.Process,
	}
	//发送请求，取得响应
	response, err := client.SendData(context.Background(), &request)
	if err != nil {
		utils.DefaultLogger.Error("推送数据失败:", err)
	} else {
		utils.DefaultLogger.Info("消息推送成功:", response)
	}
	fmt.Println()
}
