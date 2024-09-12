package CMDBClient

import (
	"cmdb_mini_agent/pkg/utils/loggers"
	"crypto/sha1"
	"encoding/hex"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func CheckClient() {
	// 构建请求参数
	params := map[string]string{}
	urlPath := "/api/v0.1/ci/s"

	fullURL, err := CmdbClient(urlPath, params)
	if err != nil {
		loggers.DefaultLogger.Error("CMDB客户端连接串配置错误", err)
		panic("CMDB客户端连接串配置错误")
	}
	//fmt.Println("完整的url：", fullURL)

	// 发送HTTP GET请求
	resp, err := http.Get(fullURL)
	if err != nil {
		//fmt.Println("发送请求时出错:", err)
		loggers.DefaultLogger.Error("CMDB客户端连接检测请求发送失败：", err)
		panic("CMDB客户端连接检测请求发送失败")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			loggers.DefaultLogger.Error("CMDB客户端连接关闭失败：", err)
		}
	}(resp.Body)

	// 读取并打印响应体
	if resp.StatusCode != 200 {
		loggers.DefaultLogger.Error("CMDB客户端连接检测失败：")
		panic("CMDB客户端连接检测失败,错误码为：" + strconv.Itoa(resp.StatusCode))
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("Error reading response body:", err)
		loggers.DefaultLogger.Error("读取响应结果时异常：", err)
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
func Flurl(apiPath string, params map[string]string) (string, error) {
	hashCMDBSecret, err := HashCMDBSecret(apiPath, params)
	if err != nil {
		loggers.DefaultLogger.Error("密钥对hash加密失败：", err)
	}
	fullURL := viper.GetString("cmdb.address") + apiPath + "?" + hashCMDBSecret

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

	return valuesEncoded, nil
}
