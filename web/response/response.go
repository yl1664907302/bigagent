package response

import (
	"encoding/json"
	"net/http"
)

// 定义响应结构体
type ResponseData struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// 成功响应
func SuccessWithDetailed(w http.ResponseWriter, msg string, data interface{}) {
	response := ResponseData{
		Code: 0,
		Data: data,
		Msg:  msg,
	}
	json.NewEncoder(w).Encode(response)
}

// 失败响应
func FailWithDetailed(w http.ResponseWriter, msg string, data interface{}) {
	response := ResponseData{
		Code: 1, // 或者其他表示错误状态的代码
		Data: data,
		Msg:  msg,
	}
	json.NewEncoder(w).Encode(response)
}
