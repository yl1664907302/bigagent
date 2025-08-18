package service

import (
	"bigagent/internal/config/global"
	"bigagent/internal/model"
	"bigagent/internal/scrape/machine"
	utils "bigagent/internal/utils"
	"bigagent/internal/web/response"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func CreateMachineUUID() {
	var requestParams = map[string]string{
		"ci_type":               global.V.GetString("veops.ciType"),
		"uuid":                  machine.SmpMa.Uuid,
		"cmdb_auto_update_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	resp, err := PostCMDBci(requestParams)
	if err != nil {
		utils.DefaultLogger.Error("唯一键数据更新失败：", err)
		return
	}
	utils.DefaultLogger.Info("唯一键数据创建成功：", requestParams)

	oIdStr := strconv.FormatFloat((*resp)["ci_id"].(float64), 'f', -1, 64)
	requestParams["oid"] = oIdStr
	requestParams["cmdb_auto_update_time"] = time.Now().Format("2006-01-02 15:04:05")
	_, err = PutCMDBci(requestParams, "")
	if err != nil {
		utils.DefaultLogger.Error("oid数据更新失败：", requestParams)
	}
	utils.DefaultLogger.Info("oid更新成功：", requestParams)
}

func UpdateMachineData() {
	var wg sync.WaitGroup

	var dataMap = map[string]string{}
	//	这里是每个数据所必须的信息，模型名称、主机IP、唯一键
	dataMap["ci_type"] = global.V.GetString("veops.ciType")
	dataMap["oid"] = ""
	os := runtime.GOOS
	switch os {
	case "windows":
		utils.DefaultLogger.Infof("主机是windows类型主机")
		wg.Add(1)
		go updateData(dataMap, &wg)
	case "linux":
		utils.DefaultLogger.Infof("主机是linux类型主机")
		wg.Add(1)
		go updateData(dataMap, &wg)
	default:
		utils.DefaultLogger.Infof("当前机器类型未支持采集...")
	}
	wg.Wait()
}

func updateData(dataMap map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	requestParams := map[string]string{}
	// 导入部分采集层数据
	data := model.NewVeoGuoChuData()
	// 推送数据
	requestParams["ci_type"] = dataMap["ci_type"]
	requestParams["oid"] = dataMap["oid"]
	requestParams["uuid"] = data.UUID
	requestParams["cmdb_auto_product_uuid"] = data.UUID
	requestParams["cmdb_auto_cpu"] = data.CmdbAutoCpu
	requestParams["cmdb_auto_memory"] = data.CmdbAutoMemory
	requestParams["cmdb_auto_disk"] = data.CmdbAutoDisk
	requestParams["cmdb_auto_machine_type"] = data.CmdbAutoMachineType
	requestParams["cmdb_auto_env"] = data.CmdbAutoEnv
	requestParams["cmdb_auto_update_time"] = data.CmdbAutoUpdateTime
	requestParams["cmdb_auto_os_type"] = data.CmdbAutoOstype
	requestParams["cmdb_auto_nodename"] = data.CmdbAutoNodename
	requestParams["cmdb_auto_architecture"] = data.CmdbAutoArchitecture
	requestParams["cmdb_auto_system_release"] = data.CmdbAutoSystemRelease
	requestParams["cmdb_auto_system_vendor"] = data.CmdbAutoSystemVendor
	requestParams["cmdb_auto_bios_date"] = data.CmdbAutoBiosDate
	requestParams["cmdb_auto_ipaddr"] = data.IP
	_, err := PutCMDBci(requestParams, "")
	if err != nil {
		utils.DefaultLogger.Error("主机数据更新失败：", err)
		return
	}
	utils.DefaultLogger.Info("数据更新成功")
}

func Query(q string, result any) (*response.QueryResp, error) {
	//定义CMDB的CI查询接口
	path := "/api/v0.1/ci/s"

	//	这里将默认的模型信息和变动的ci数据结合
	param := map[string]string{
		"sort":    "",       //	属性的排序，降序字段前面加负号-
		"page":    "1",      // 页数
		"count":   "100000", //	一页返回的CI数
		"ret_key": "name",   //	返回字段类型,这里规定只能使用name
		"q":       q,
	}

	//拼接完整的CMDB连接串
	fullURL, err := Flurl(path, param)
	if err != nil {
		//loggers.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	// 发送HTTP GET请求
	resp, err := http.Get(fullURL)
	if err != nil {
		//loggers.DefaultLogger.Error("CMDB客户端连接检测请求发送失败：", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//loggers.DefaultLogger.Error("CMDB客户端连接检测失败：")
		return nil, err
	}
	if resp.StatusCode != 200 {
		errorData, err := response.ResponseErrorData(resp.StatusCode, body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errorData)
	}

	queryResp := response.QueryResp{}
	err = json.Unmarshal(body, &queryResp)
	if err != nil {
		//loggers.DefaultLogger.Error("CMDB解析json失败：", err)
		return nil, err
	}

	err = json.Unmarshal(queryResp.Result, result)
	if err != nil {
		//loggers.DefaultLogger.Error("CMDB解析json失败：", err)
		return nil, err
	}

	return &queryResp, nil
}

// QueryCMDBciDataProcess 针对查询到的数据进行json格式化
func QueryCMDBciDataProcess(body []byte) (*response.QueryCMDBciResponse, error) {
	var jsonData *response.QueryCMDBciResponse

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

func QueryCMDBci(params map[string]string) (*response.QueryCMDBciResponse, error) {

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
	fullURL, err := Flurl(urlPath, allParams)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	//	调试
	//fmt.Println("完整的url：", fullURL)

	// 发送HTTP GET请求
	resp, err := http.Get(fullURL)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接检测请求发送失败：", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接检测失败：")
		return nil, err
	}
	if resp.StatusCode != 200 {
		errorData, err := response.ResponseErrorData(resp.StatusCode, body)
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

func PostCMDBci(params map[string]string) (*response.ResultResponse, error) {

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
	fullURL, err := Flurl(urlPath, allParams)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	//	调试
	//fmt.Println("完整的url：", fullURL)

	// 发送HTTP Post请求
	resp, err := http.Post(fullURL, "application/json", nil)
	//resp, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接检测请求发送失败：", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端响应体读取异常：", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		errorData, err := response.ResponseErrorData(resp.StatusCode, body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errorData)
	}
	jsonData, err := response.ResultResponseDataProcess(body)
	if err != nil {
		return nil, err
	}

	utils.DefaultLogger.Info("新增数据返回结果为：", jsonData)
	return jsonData, nil
}

func PutCMDBci(params map[string]string, ciId string) (*response.ResultResponse, error) {

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
	fullURL, err := Flurl(urlPath, allParams)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	//	调试
	//fmt.Println("完整的url：", fullURL)

	// 初始化一个Client客户端
	client := &http.Client{}

	//// 发送HTTP PUT请求
	//fmt.Println(fullURL)
	req, err := http.NewRequest("PUT", fullURL, nil)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接请求发送失败：", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接异常：", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("body:", string(body))
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端响应体读取异常：", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		errorData, err := response.ResponseErrorData(resp.StatusCode, body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errorData)
	}

	jsonData, err := response.ResultResponseDataProcess(body)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func DeleteCMDBci(ciId string) (*response.ResultResponse, error) {

	//定义CMDB的CI删除数据接口
	urlPath := "/api/v0.1/ci/" + ciId

	//拼接完整的CMDB连接串
	fullURL, err := Flurl(urlPath, nil)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	//	调试
	//fmt.Println("完整的url：", fullURL)

	// 初始化一个Client客户端
	client := &http.Client{}

	// 发送HTTP PUT请求
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接请求发送失败：", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接异常：", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端响应体读取异常：", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		errorData, err := response.ResponseErrorData(resp.StatusCode, body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(errorData)
	}
	jsonData, err := response.ResultResponseDataProcess(body)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func CheckClient(host string) {
	// 构建请求参数
	params := map[string]string{}
	urlPath := "/api/v0.1/ci/s"

	fullURL, err := CmdbClient(urlPath, params)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		panic("CMDB客户端连接串配置错误")
	}
	//fmt.Println("完整的url：", fullURL)

	// 发送HTTP GET请求
	resp, err := http.Get(fullURL)
	if err != nil {
		//fmt.Println("发送请求时出错:", err)
		utils.DefaultLogger.Error("CMDB客户端连接检测请求发送失败：", err)
		panic("CMDB客户端连接检测请求发送失败")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		utils.DefaultLogger.Error("CMDB客户端连接检测失败：", string(bodyBytes))
		panic("CMDB客户端连接检测失败,错误码为：" + strconv.Itoa(resp.StatusCode))

	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("Error reading response body:", err)
		utils.DefaultLogger.Error("读取响应结果时异常：", err)
		panic("读取CMDB客户端响应结果失败")
	}
}

func CmdbClient(apiPath string, params map[string]string) (string, error) {

	//params := map[string]string{}
	var keys []string

	for k := range params {
		keys = append(keys, k)
	}

	//	对keys进行可排序格式的转换，转换后进行排序
	keysSorted := sort.StringSlice(keys)
	keysSorted.Sort()
	keys = []string(keysSorted)

	//	使用一个空的字符串切片来过滤掉key和secret
	var values []string
	for _, k := range keys {
		if k == "_key" {
			continue
		}
		if k == "_secret" {
			continue
		}
		values = append(values, params[k])
	}
	valuesStr := strings.Join(values, "")

	//	调试信息
	//fmt.Println("valuesStr的值为：", valuesStr)

	//	构造用于计算SHA1签名的字符串
	_secretPre := apiPath + global.V.GetString("veops.secret") + valuesStr
	//	计算SHA1签名
	_secretH := sha1.Sum([]byte(_secretPre))
	_secret := hex.EncodeToString(_secretH[:])

	// 将计算得到的签名和Key添加到params
	params["_secret"] = _secret
	params["_key"] = global.V.GetString("veops.key")

	// 构造完整URL
	valuesMap := url.Values{}
	for k, v := range params {
		if v != "" {
			valuesMap.Add(k, v)
		}
	}
	valuesEncoded := valuesMap.Encode()
	fullURL := global.V.GetString("veops.address") + apiPath + "?" + valuesEncoded

	//	调试信息
	//fmt.Printf("%+v", params)
	//fmt.Println("完整的url为：", fullURL)

	return fullURL, nil
}

// Flurl	将加密后的密钥和完整URL进行拼接
func Flurl(apiPath string, params map[string]string) (string, error) {
	hashCMDBSecret, err := HashCMDBSecret(apiPath, params)
	if err != nil {
		utils.DefaultLogger.Error("密钥对hash加密失败：", err)
	}
	fullURL := global.V.GetString("veops.address") + apiPath + "?" + hashCMDBSecret

	return fullURL, nil
}

// HashCMDBSecret 对密钥对进行一次hash
func HashCMDBSecret(apiPath string, params map[string]string) (string, error) {
	//params := map[string]string{}
	var keys []string
	if params == nil {
		params = make(map[string]string)
	}
	for k := range params {
		keys = append(keys, k)
	}

	//	对keys进行可排序格式的转换，转换后进行排序
	keysSorted := sort.StringSlice(keys)
	keysSorted.Sort()
	keys = []string(keysSorted)

	//	使用一个空的字符串切片来过滤掉key和secret
	var values []string
	for _, k := range keys {
		if k == "_key" {
			continue
		}
		if k == "_secret" {
			continue
		}
		values = append(values, params[k])
	}
	valuesStr := strings.Join(values, "")

	//	调试信息
	//fmt.Println("valuesStr的值为：", valuesStr)

	//	构造用于计算SHA1签名的字符串
	_secretPre := apiPath + global.V.GetString("veops.secret") + valuesStr
	//	计算SHA1签名
	_secretH := sha1.Sum([]byte(_secretPre))
	_secret := hex.EncodeToString(_secretH[:])

	// 将计算得到的签名和Key添加到params
	params["_secret"] = _secret
	params["_key"] = global.V.GetString("veops.key")

	// 构造完整URL
	valuesMap := url.Values{}
	for k, v := range params {
		if v != "" {
			valuesMap.Add(k, v)
		}
	}
	valuesEncoded := valuesMap.Encode()

	return valuesEncoded, nil
}
