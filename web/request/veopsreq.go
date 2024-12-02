package request

import (
	model "bigagent/model/machine"
	utils "bigagent/util"
	"bigagent/web/response"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PostVeops struct {
	k bool
	h string
	c *http.Client
	d *model.VeopsData
}

func NewPostVeops(host string) *PostVeops {
	return &PostVeops{h: host, c: &http.Client{}, d: model.NewVeopsdata()}
}

func (p *PostVeops) Do() (interface{}, error) {
	p.CreateMachineUUID()
	p.UpdateMachineData()
	return nil, nil
}

func (p *PostVeops) CreateMachineUUID() {
	if !p.k {
		veopsdata := model.NewVeopsdata()

		responses, err := p.PostCMDBci(veopsdata)
		if err != nil {
			utils.DefaultLogger.Error("唯一键数据更新失败：", veopsdata)
			return
		}
		utils.DefaultLogger.Info("唯一键数据创建成功：", veopsdata)
		p.k = true
		oIdStr := strconv.FormatFloat((*responses)["ci_id"].(float64), 'f', -1, 64)
		utils.DefaultLogger.Info("创建实例成功ci_id为：", oIdStr)
	}
}

func (p *PostVeops) UpdateMachineData() {
	ipAddr := "192.0.0.1"
	//if err != nil {
	//	logger.DefaultLogger.Error("主机数据更新失败：", err)
	//	return
	//}
	veopsdata := model.NewVeopsdata()
	veopsdata.Cmdb_auto_product_uuid = "uuid"
	veopsdata.Cmdb_auto_ipaddr = ipAddr
	veopsdata.Cmdb_auto_env = viper.GetString("global.env")
	veopsdata.Cmdb_auto_machine_type = viper.GetString("global.machineType")
	veopsdata.Cmdb_auto_os_type = "linux"
	veopsdata.Cmdb_auto_update_time = time.Now().Format("2006-01-02 15:04:05")

	_, err := p.PutCMDBci(veopsdata, "")
	if err != nil {
		utils.DefaultLogger.Error("主机数据更新失败：", veopsdata)
		return
	}
	utils.DefaultLogger.Info("数据更新成功：", veopsdata)
}

func (p *PostVeops) CheckClient() {
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
		utils.DefaultLogger.Error("CMDB客户端连接检测失败：")
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
	_secretPre := apiPath + viper.GetString("cmdb.Secret") + valuesStr
	//	计算SHA1签名
	_secretH := sha1.Sum([]byte(_secretPre))
	_secret := hex.EncodeToString(_secretH[:])

	// 将计算得到的签名和Key添加到params
	params["_secret"] = _secret
	params["_key"] = viper.GetString("cmdb.Key")

	// 构造完整URL
	valuesMap := url.Values{}
	for k, v := range params {
		if v != "" {
			valuesMap.Add(k, v)
		}
	}
	valuesEncoded := valuesMap.Encode()
	fullURL := viper.GetString("cmdb.address") + apiPath + "?" + valuesEncoded

	//	调试信息
	//fmt.Printf("%+v", params)
	//fmt.Println("完整的url为：", fullURL)

	return fullURL, nil
}

// Flurl	将加密后的密钥和完整URL进行拼接
func (p *PostVeops) Flurl(apiPath string, params model.VeopsData) (string, error) {
	err := HashCMDBSecret(apiPath, params)
	if err != nil {
		utils.DefaultLogger.Error("密钥对hash加密失败：", err)
	}
	hashCMDBSecret, err := utils.JSONToFormData(params)
	fullURL := p.h + apiPath + "?" + hashCMDBSecret

	return fullURL, nil
}

// HashCMDBSecret 对密钥对进行一次hash
func HashCMDBSecret(apiPath string, params model.VeopsData) error {
	if params.Cmdb_auto_product_uuid == "" {
		params = model.VeopsData{}
	}
	ks, vs := utils.GetNonEmptyFields(params)

	//	对keys进行可排序格式的转换，转换后进行排序
	keysSorted := sort.StringSlice(ks)
	keysSorted.Sort()
	ks = []string(keysSorted)

	//  移除指定key值
	utils.RemoveString(vs, "_secret")
	utils.RemoveString(vs, "_key")
	valuesStr := strings.Join(vs, "")

	//	构造用于计算SHA1签名的字符串
	_secretPre := apiPath + viper.GetString("cmdb.Secret") + valuesStr
	//	计算SHA1签名
	_secretH := sha1.Sum([]byte(_secretPre))
	_secret := hex.EncodeToString(_secretH[:])

	// 将计算得到的签名和Key添加到params
	params.Secret = _secret
	params.Key = viper.GetString("cmdb.Key")

	return nil
}

func (p *PostVeops) QueryCMDBci(params *model.VeopsData) (*response.QueryCMDBciResponse, error) {
	params.Sort = ""
	params.Page = "1"
	params.Count = "9999"
	params.Ret_key = "name"
	//定义CMDB的CI查询接口
	urlPath := "/api/v0.1/ci/s"

	//拼接完整的CMDB连接串
	fullURL, err := p.Flurl(urlPath, *params)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

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
	jsonData, err := response.QueryCMDBciDataProcess(body)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (p *PostVeops) PostCMDBci(params *model.VeopsData) (*response.ResultResponse, error) {
	//定义CMDB的CI新增数据接口
	urlPath := "/api/v0.1/ci"

	params.No_attribute_policy = "ignore" //	当添加不存在的attribute时的策略, 可选: reject、ignore, 默认ignore
	params.Exist_policy = "reject"        //	CI已经存在的处理策略, 可选: need、reject、replace 默认reject

	//拼接完整的CMDB连接串
	fullURL, err := p.Flurl(urlPath, *params)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

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

func (p *PostVeops) PutCMDBci(params *model.VeopsData, ciId string) (*response.ResultResponse, error) {

	//定义CMDB的CI修改数据接口
	urlPath := "/api/v0.1/ci"
	if ciId != "" {
		urlPath = "/api/v0.1/ci/" + ciId
	}

	params.No_attribute_policy = "ignore" //	当添加不存在的attribute时的策略, 可选: reject、ignore, 默认ignore
	params.Exist_policy = "reject"        //	CI已经存在的处理策略, 可选: need、reject、replace 默认reject

	//拼接完整的CMDB连接串
	fullURL, err := p.Flurl(urlPath, *params)
	if err != nil {
		utils.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		return nil, err
	}

	//	调试
	//fmt.Println("完整的url：", fullURL)

	// 初始化一个Client客户端
	client := &http.Client{}

	// 发送HTTP PUT请求
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

//func DeleteCMDBci(ciId string) (*response.ResultResponse, error) {
//
//	//定义CMDB的CI删除数据接口
//	urlPath := "/api/v0.1/ci/" + ciId
//
//	//拼接完整的CMDB连接串
//	fullURL, err := Flurl(urlPath, nil)
//	if err != nil {
//		logger.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
//		return nil, err
//	}
//
//	//	调试
//	//fmt.Println("完整的url：", fullURL)
//
//	// 初始化一个Client客户端
//	client := &http.Client{}
//
//	// 发送HTTP PUT请求
//	req, err := http.NewRequest("DELETE", fullURL, nil)
//	if err != nil {
//		logger.DefaultLogger.Error("CMDB客户端连接请求发送失败：", err)
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		logger.DefaultLogger.Error("CMDB客户端连接异常：", err)
//	}
//
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			logger.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
//		}
//	}(resp.Body)
//
//	// 读取并打印响应体
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		logger.DefaultLogger.Error("CMDB客户端响应体读取异常：", err)
//		return nil, err
//	}
//	if resp.StatusCode != 200 {
//		errorData, err := response.ResponseErrorData(resp.StatusCode, body)
//		if err != nil {
//			return nil, err
//		}
//		return nil, fmt.Errorf(errorData)
//	}
//	jsonData, err := response.ResultResponseDataProcess(body)
//	if err != nil {
//		return nil, err
//	}
//	return jsonData, nil
//}
