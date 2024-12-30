package config

type Server struct {
	System System `json:"system"`
}

type System struct {
	Addr              string `json:"addr" yaml:"addr"`
	Grpc              string `json:"grpc" yaml:"grpc"`
	Serct             string `json:"serct" yaml:"serct"`
	Logfile           string `json:"logfile" yaml:"logfile"`
	Grpc_server       string `json:"grpc_server" yaml:"grpc_server"`
	Grpc_cmdb1_stand1 string `json:"grpc_cmdb_1_stand_1" yaml:"grpc_cmdb1_stand1"`
}
