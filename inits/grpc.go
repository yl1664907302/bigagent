package inits

import grpc_client "bigagent/grpcs/client"

func RunGC() {
	grpc_client.InitClient()
	grpc_client.General()
}
