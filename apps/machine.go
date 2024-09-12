package apps

import (
	"cmdb_mini_agent/model/CMDB"
	"cmdb_mini_agent/pkg/utils/base"
	"cmdb_mini_agent/pkg/utils/loggers"
	"github.com/spf13/viper"
	"github.com/super-l/machine-code/machine"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type machineConfig struct {
	UniqueKey string //	唯一键
	Oid       string //	oid
}

var machineConfigData = machineConfig{
	UniqueKey: "cmdb_auto_product_uuid",
	Oid:       "oid",
}

// GetUUID 取出主机的product_uuid(/sys/class/dmi/id/product_uuid,实际执行逻辑为，执行 dmidecode -s system-uuid 命令)
func GetUUID() string {
	jsonData := machine.GetMachineData()
	uuid := jsonData.PlatformUUID

	// 处理数据，去除前后空白并提取有效UUID
	uuid = strings.TrimSpace(uuid)
	re := regexp.MustCompile(`[0-9A-Fa-f-]{36}`)
	matches := re.FindStringSubmatch(uuid)
	if len(matches) > 0 {
		uuid = matches[0]
	} else {
		uuid = ""
	}

	//fmt.Println("路径为：", *initialize.FullDir)
	return uuid
}

// GetIp 取出主机的IP地址。（执行逻辑为，向8.8.8.8的DNS发送请求，验证是从哪个ip发出去的流量，就使用该ip）
func GetIp() (string, error) {
	ipAddr, err := machine.GetLocalIpAddr()
	if err != nil {
		loggers.DefaultLogger.Error("获取主机IP失败：", err)
		return "", err
	}
	//fmt.Println("出口IP为：", ipAddr)
	return ipAddr, nil
}

// CreateMachineUUID	往CMDB中创建uuid,创建完成后返回的ci_id再写入到CMDb中的oid字段
func CreateMachineUUID() {
	var requestParams = map[string]string{
		"ci_type":                   viper.GetString("cmdb.ciType"),
		machineConfigData.UniqueKey: GetUUID(),
		"cmdb_auto_update_time":     time.Now().Format("2006-01-02 15:04:05"),
	}

	//	调试
	//fmt.Println("requestParams：", requestParams)

	response, err := CMDB.PostCMDBci(requestParams)
	if err != nil {
		loggers.DefaultLogger.Error("唯一键数据更新失败：", requestParams)
		return
	}
	loggers.DefaultLogger.Info("唯一键数据创建成功：", requestParams)

	oIdStr := strconv.FormatFloat((*response)["ci_id"].(float64), 'f', -1, 64)
	requestParams["oid"] = oIdStr
	requestParams["cmdb_auto_update_time"] = time.Now().Format("2006-01-02 15:04:05")
	_, err = CMDB.PutCMDBci(requestParams, "")
	if err != nil {
		loggers.DefaultLogger.Error("oid数据更新失败：", requestParams)
		return
	}
	loggers.DefaultLogger.Info("oid更新成功：", requestParams)
}

func UpdateMachineData() {
	ipAddr, err := GetIp()
	if err != nil {
		loggers.DefaultLogger.Error("主机数据更新失败：", err)
		return
	}

	var requestParams = map[string]string{
		"ci_type":                   viper.GetString("cmdb.ciType"),
		machineConfigData.UniqueKey: GetUUID(),
		"cmdb_auto_ipaddr":          ipAddr,
		"cmdb_auto_env":             viper.GetString("global.env"),
		"cmdb_auto_machine_type":    viper.GetString("global.machineType"),
		"cmdb_auto_os_type":         *base.OsType,
		"cmdb_auto_update_time":     time.Now().Format("2006-01-02 15:04:05"),
	}

	_, err = CMDB.PutCMDBci(requestParams, "")
	if err != nil {
		loggers.DefaultLogger.Error("主机数据更新失败：", requestParams)
		return
	}
	loggers.DefaultLogger.Info("数据更新成功：", requestParams)
}
