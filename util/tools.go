package utils

import (
	config2 "bigagent/config"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
)

// 获取结构体实例中绑定了json标签的key值
func GetJSONKeysFromInstance(v interface{}) []string {
	var keys []string

	// 获取传入对象的值和类型
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// 如果传入的是指针，获取指针指向的元素值和类型
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// 确保是结构体
	if typ.Kind() != reflect.Struct {
		return keys
	}

	// 遍历结构体的字段
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		valueField := val.Field(i)

		// 获取 json 标签
		jsonTag := field.Tag.Get("json")

		// 处理标签中的 "omitempty" 等情况，只取第一个逗号前的部分
		if jsonTag != "" && jsonTag != "-" {
			jsonKey := strings.Split(jsonTag, ",")[0]

			// 如果字段有值（非零值）才取出key
			if !valueField.IsZero() {
				keys = append(keys, jsonKey)
			}
		}
	}

	return keys
}

// 获取值不为""的字段名和字段值
func GetNonEmptyFields(v interface{}) ([]string, []string) {
	var fieldNames []string
	var fieldValues []string

	val := reflect.ValueOf(v)

	// 确保是结构体
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // 解引用
	}
	if val.Kind() != reflect.Struct {
		return fieldNames, fieldValues
	}

	// 遍历结构体的字段
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()

		// 检查字段值是否非空字符串
		if strValue, ok := value.(string); ok && strValue != "" {
			fieldNames = append(fieldNames, field.Name)
			fieldValues = append(fieldValues, strValue)
		}
	}

	return fieldNames, fieldValues
}

// RemoveString 从切片中移除指定的字符串
func RemoveString(slice []string, str string) []string {
	var result []string
	for _, item := range slice {
		if item != str {
			result = append(result, item)
		}
	}
	return result
}

// JSON 转换为表单数据
func JSONToFormData(jsonData interface{}) (string, error) {

	ks, vs := GetNonEmptyFields(jsonData)

	// 创建 url.Values 对象
	formData := url.Values{}

	// 遍历 JSON 数据，将其转换为表单数据
	for i := 0; i < len(ks); i++ {
		formData.Set(ks[i], fmt.Sprintf("%v", vs[i]))
	}

	// 编码为表单格式
	return formData.Encode(), nil
}

// ModifyYAML 修改 YAML 文件中的字段值
func ModifyYAML(filePath, fieldPath, newValue string) error {
	// 读取 YAML 文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("无法读取文件: %w", err)
	}

	// 解析 YAML 数据
	var config config2.Server
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("解析 YAML 失败: %w", err)
	}

	// 修改指定字段
	switch fieldPath {
	case "system.serct":
		config.System.Serct = newValue
	case "system.addr":
		config.System.Addr = newValue
	case "system.grpc":
		config.System.Grpc = newValue
	case "system.logfile":
		config.System.Logfile = newValue
	case "system.grpc_server":
		config.System.Grpc_server = newValue
	case "system.grpc_cmdb1_stand1":
		config.System.Grpc_cmdb1_stand1 = newValue
	default:
		return fmt.Errorf("未找到字段: %s", fieldPath)
	}

	// 使用两格缩进序列化 YAML 数据，并保留双引号
	modifiedData, err := formatYAML(config)
	if err != nil {
		return fmt.Errorf("格式化 YAML 失败: %w", err)
	}

	// 写回修改后的 YAML 文件
	err = ioutil.WriteFile(filePath, modifiedData, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// formatYAML 格式化 YAML 输出，保留两格缩进
func formatYAML(config interface{}) ([]byte, error) {
	var buf strings.Builder
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2) // 设置两格缩进
	defer encoder.Close()

	err := encoder.Encode(config)
	if err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}
