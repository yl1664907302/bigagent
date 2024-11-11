package response

import (
	"encoding/json"
	"net/http"
)

//// 定义响应结构体
//type ResponseData struct {
//	//Data interface{} `json:"data"`
//	model.StandData
//}

// 成功响应
func SuccessWithDetailed(w http.ResponseWriter, data interface{}) {
	response := data

	json.NewEncoder(w).Encode(response)
}

//// 失败响应
//func FailWithDetailed(w http.ResponseWriter, data interface{}) {
//	response := ResponseData{
//		Data: data,
//	}
//	json.NewEncoder(w).Encode(response)
//}
