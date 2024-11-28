package model

type AuthDetails interface {
	ApplyAuth(args ...interface{}) error
}

type FieldMapping struct {
	StructField string `json:"struct_field"`
	Type        string `json:"type"`
}

type NetworkInfo struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Path     string `json:"path"`
}

type AgentConfig struct {
	AuthName     string                  `json:"auth_name"`
	AuthDetails  AuthDetails             `json:"auth_details"`
	FieldMapping map[string]FieldMapping `json:"field_mapping"`
	NetworkInfo  NetworkInfo             `json:"network_info"`
}

// 接收案例
//{
//"auth_mode": "token",
//"auth_details": {
//"token": "my-secret-token"
//},
//"field_mapping": {
//"cmdb_field1": { "struct_field": "Name", "type": "string" },
//"cmdb_field2": { "struct_field": "Age", "type": "integer" }
//},
//"data_rules": {
//"generate_data": true,
//"fields": {
//"Name": { "default": "Default Name" },
//"Age": { "default": 30 }
//}
//},
//"network_info": {
//"protocol": "http",
//"host": "cmdb.example.com",
//"port": 8080,
//"path": "/api/v1/data"
//}
//}
