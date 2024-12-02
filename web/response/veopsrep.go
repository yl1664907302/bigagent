package response

import (
	utils "bigagent/util"
	"encoding/json"
)

// QueryCMDBciResponse	CMDBci查询接口返回的数据类型
type QueryCMDBciResponse struct {
	Counter  map[string]int         `json:"counter"`  //	当前页按模型的分类统计
	Facet    map[string]interface{} `json:"facet"`    //	返回的CI列表
	Numfound int                    `json:"numfound"` //	CI总数
	Page     int                    `json:"page"`     //	分页
	Result   []interface{}          `json:"result"`   //	返回的CI列表
	Total    int                    `json:"total"`    //	当前页的CI数
}

// ErrorMessage	查询报错结果返回的数据类型
type ErrorMessage struct {
	Message string `json:"message"`
}

// ResultResponse	CMDBci 增、改、删接口返回的数据类型
type ResultResponse map[string]interface{}

// ErrorMessageProcess CMDBci 错误信息的格式化
func ErrorMessageProcess(body []byte) (*ErrorMessage, error) {
	var jsonData *ErrorMessage

	_, err := json.Marshal(body)
	if err != nil {
		utils.DefaultLogger.Error("json转换失败")
		return nil, err
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		utils.DefaultLogger.Error("解析json出错：", err)
		return nil, err
	}
	return jsonData, nil
}

// ResultResponseDataProcess CMDBci 增、改、删接口返回的数据进行json格式化
func ResultResponseDataProcess(body []byte) (*ResultResponse, error) {
	var jsonData *ResultResponse

	_, err := json.Marshal(body)
	if err != nil {
		utils.DefaultLogger.Error("json转换失败")
		return nil, err
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		utils.DefaultLogger.Error("解析json出错：", err)
		return nil, err
	}
	return jsonData, nil
}

// QueryCMDBciDataProcess 针对查询到的数据进行json格式化
func QueryCMDBciDataProcess(body []byte) (*QueryCMDBciResponse, error) {
	var jsonData *QueryCMDBciResponse

	_, err := json.Marshal(body)
	if err != nil {
		utils.DefaultLogger.Error("json转换失败")
		return nil, err
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		utils.DefaultLogger.Error("解析json出错：", err)
		return nil, err
	}
	return jsonData, nil
}

// ResponseErrorData	用来转换接口错误时返回的各种信息
func ResponseErrorData(respStatusCode int, body []byte) (message string, err error) {
	jsonData, err := ErrorMessageProcess(body)
	if err != nil {
		utils.DefaultLogger.Error("json转换失败")
		return "", err
	}

	switch respStatusCode {
	default:
		utils.DefaultLogger.Error("错误码:", respStatusCode, " message:", jsonData.Message)
		return jsonData.Message, nil
	}
}
