package response

import (
	utils "bigagent/internal/util"
	"encoding/json"
)

type QueryCMDBciResponse struct {
	Counter  map[string]int         `json:"counter"`  //	当前页按模型的分类统计
	Facet    map[string]interface{} `json:"facet"`    //	返回的CI列表
	Numfound int                    `json:"numfound"` //	CI总数
	Page     int                    `json:"page"`     //	分页
	Result   []interface{}          `json:"result"`   //	返回的CI列表
	Total    int                    `json:"total"`    //	当前页的CI数
}

// QueryResp	Query查询接口返回的数据类型
type QueryResp struct {
	Counter  map[string]int             `json:"counter"`  //	当前页按模型的分类统计
	Facet    map[string]json.RawMessage `json:"facet"`    //	返回的CI列表
	Numfound int                        `json:"numfound"` //	CI总数
	Page     int                        `json:"page"`     //	分页
	Result   json.RawMessage            `json:"result"`   //	返回的CI列表
	Total    int                        `json:"total"`    //	当前页的CI数
}

type ErrorMessage struct {
	Message string `json:"message"`
}

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
