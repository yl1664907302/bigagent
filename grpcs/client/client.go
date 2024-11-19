package grpc_client

import (
	"bigagent/config/global"
	grpc_server "bigagent/grpcs/server"
	model "bigagent/model/machine"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var (
	client grpc_server.PushAgantDataClient
)

func InitClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond) //连接超时设置为1000毫秒
	defer cancel()
	//连接到服务端
	conn, err := grpc.DialContext(
		ctx,
		global.CONF.System.Grpc,
		grpc.WithTransportCredentials(insecure.NewCredentials()), //Credential即使为空，也必须设置
		grpc.WithBlock(), //grpc.WithBlock()直到连接真正建立才会返回，否则连接是异步建立的。因此grpc.WithBlock()和Timeout结合使用才有意义。server端正常的情况下使用grpc.WithBlock()得到的connection.GetState()为READY，不使用grpc.WithBlock()得到的connection.GetState()为IDEL
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(10<<20), grpc.MaxCallRecvMsgSize(10<<20)), //默认情况下SendMsg上限是MaxInt32，RecvMsg上限是4M，这里都修改为10M
	)
	if err != nil {
		fmt.Printf("dial failed: %s", err)
		return
	}
	//创建client
	client = grpc_server.NewPushAgantDataClient(conn)
	if client == nil {
		fmt.Println("gRPC client is not initialized.")
		return
	}
}

func General() {
	//准备好请求参数
	data := model.NewStandData()
	request := grpc_server.StandData{
		Serct:    data.Serct,
		Uuid:     data.Uuid,
		Hostname: data.Hostname,
		Ipv4:     data.IPv4,
		Time:     uint64(time.Now().Unix()),
		Info:     nil,
		Cpu:      nil,
		Disk:     nil,
		Memory:   nil,
		Net:      nil,
	}
	//发送请求，取得响应
	response, err := client.SendData(context.Background(), &request)
	if err != nil {
		fmt.Printf("推送数据失败: %s", err)
	} else {
		fmt.Printf("消息推送成功：%s", response)
	}
	fmt.Println()
}
