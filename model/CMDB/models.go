package CMDB

import (
	"cmdb_mini_agent/pkg/utils/CMDBClient"
	"cmdb_mini_agent/pkg/utils/loggers"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
		loggers.DefaultLogger.Error("json转换失败")
		return nil, err
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		loggers.DefaultLogger.Error("解析json出错：", err)
		return nil, err
	}
	return jsonData, nil
}

// ResultResponseDataProcess CMDBci 增、改、删接口返回的数据进行json格式化
func ResultResponseDataProcess(body []byte) (*ResultResponse, error) {
	var jsonData *ResultResponse

	_, err := json.Marshal(body)
	if err != nil {
		loggers.DefaultLogger.Error("json转换失败")
		return nil, err
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		loggers.DefaultLogger.Error("解析json出错：", err)
		return nil, err
	}
	return jsonData, nil
}

// QueryCMDBciDataProcess 针对查询到的数据进行json格式化
func QueryCMDBciDataProcess(body []byte) (*QueryCMDBciResponse, error) {
	var jsonData *QueryCMDBciResponse

	_, err := json.Marshal(body)
	if err != nil {
		loggers.DefaultLogger.Error("json转换失败")
		return nil, err
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		loggers.DefaultLogger.Error("解析json出错：", err)
		return nil, err
	}
	return jsonData, nil
}

// ResponseErrorData	用来转换接口错误时返回的各种信息
func ResponseErrorData(respStatusCode int, body []byte) (message string, err error) {
	jsonData, err := ErrorMessageProcess(body)
	if err != nil {
		loggers.DefaultLogger.Error("json转换失败")
		return "", err
	}

	switch respStatusCode {
	default:
		loggers.DefaultLogger.Error("错误码:", respStatusCode, " message:", jsonData.Message)
		return jsonData.Message, nil
	}
}

func QueryCMDBci(params map[string]string) (*QueryCMDBciResponse, error) {

	defaultParams := map[string]string{
		"sort":    "",     //	属性的排序，降序字段前面加负号-
		"page":    "1",    // 页数
		"count":   "9999", //	一页返回的CI数
		"ret_key": "name", //	返回字段类型,这里规定只能使用name
	}
	//定义CMDB的CI查询接口
	urlPath := "/api/v0.1/ci/s"

	//	这里将默认的模型信息和变动的ci数据结合
	allParams := make(map[string]string, len(params)+len(defaultParams))
	for k, v := range params {
		allParams[k] = v
	}
	for k, v := range defaultParams {
		allParams[k] = v
	}

	//拼接完整的CMDB连接串
	fullURL, err := CMDBClient.Flurl(urlPath, allParams)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	//	调试
	//fmt.Println("完整的url：", fullURL)

	// 发送HTTP GET请求
	resp, err := http.Get(fullURL)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接检测请求发送失败：", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			loggers.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接检测失败：")
		return nil, err
	}
	if resp.StatusCode != 200 {
		errorData, err := ResponseErrorData(resp.StatusCode, body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errorData)
	}
	jsonData, err := QueryCMDBciDataProcess(body)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func PostCMDBci(params map[string]string) (*ResultResponse, error) {

	defaultParams := map[string]string{
		"no_attribute_policy": "ignore", //	当添加不存在的attribute时的策略, 可选: reject、ignore, 默认ignore
		"exist_policy":        "reject", //	CI已经存在的处理策略, 可选: need、reject、replace 默认reject
	}
	//定义CMDB的CI新增数据接口
	urlPath := "/api/v0.1/ci"

	//	这里将默认的模型信息和变动的ci数据结合
	allParams := make(map[string]string, len(params)+len(defaultParams))
	for k, v := range params {
		allParams[k] = v
	}
	for k, v := range defaultParams {
		allParams[k] = v
	}

	//拼接完整的CMDB连接串
	fullURL, err := CMDBClient.Flurl(urlPath, allParams)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	//	调试
	//fmt.Println("完整的url：", fullURL)

	// 发送HTTP Post请求
	resp, err := http.Post(fullURL, "application/json", nil)
	//resp, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接检测请求发送失败：", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			loggers.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端响应体读取异常：", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		errorData, err := ResponseErrorData(resp.StatusCode, body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errorData)
	}
	jsonData, err := ResultResponseDataProcess(body)
	if err != nil {
		return nil, err
	}

	loggers.DefaultLogger.Info("新增数据返回结果为：", jsonData)
	return jsonData, nil
}

func PutCMDBci(params map[string]string, ciId string) (*ResultResponse, error) {

	defaultParams := map[string]string{
		"no_attribute_policy": "ignore", //	当添加不存在的attribute时的策略, 可选: reject、ignore, 默认ignore
	}
	//定义CMDB的CI修改数据接口
	urlPath := "/api/v0.1/ci"
	if ciId != "" {
		urlPath = "/api/v0.1/ci/" + ciId
	}

	//	这里将默认的模型信息和变动的ci数据结合
	allParams := make(map[string]string, len(params)+len(defaultParams))
	for k, v := range params {
		allParams[k] = v
	}
	for k, v := range defaultParams {
		allParams[k] = v
	}

	//拼接完整的CMDB连接串
	fullURL, err := CMDBClient.Flurl(urlPath, allParams)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	//	调试
	//fmt.Println("完整的url：", fullURL)

	// 初始化一个Client客户端
	client := &http.Client{}

	// 发送HTTP PUT请求
	req, err := http.NewRequest("PUT", fullURL, nil)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接请求发送失败：", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接异常：", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			loggers.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端响应体读取异常：", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		errorData, err := ResponseErrorData(resp.StatusCode, body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errorData)
	}

	jsonData, err := ResultResponseDataProcess(body)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func DeleteCMDBci(ciId string) (*ResultResponse, error) {

	//定义CMDB的CI删除数据接口
	urlPath := "/api/v0.1/ci/" + ciId

	//拼接完整的CMDB连接串
	fullURL, err := CMDBClient.Flurl(urlPath, nil)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	//	调试
	//fmt.Println("完整的url：", fullURL)

	// 初始化一个Client客户端
	client := &http.Client{}

	// 发送HTTP PUT请求
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接请求发送失败：", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接异常：", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			loggers.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端响应体读取异常：", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		errorData, err := ResponseErrorData(resp.StatusCode, body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errorData)
	}
	jsonData, err := ResultResponseDataProcess(body)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
