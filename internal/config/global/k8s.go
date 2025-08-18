package global

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 先定义一个空的 config
type Config struct {
	HttpAddr              string                `yaml:"http_addr"`
	MysqlC                *MySQLConf            `yaml:"mysql_info"`
	ImDingDingC           *ImDingDingConf       `yaml:"im_ding_ding"`
	K8sConfigS            map[string]string     `yaml:"k8s_configs"` // k8s多集群的kubeconfig
	K8sTimeoutSeconds     int                   `yaml:"k8s_timeout_seconds"`
	K8sComputePromAddrMap map[string]string     `yaml:"k8s_compute_prom_addr_map"`
	PluginNtpC            *PluginNtpConf        `yaml:"plugin_ntp"`
	RecoveryC             *RecoveryConf         `yaml:"recovery_conf"`
	ModuleCommonC         *ModuleCommonConf     `yaml:"common_module"`
	AbnormalPodCleanC     *AbnormalPodCleanConf `yaml:"abnormal_pod_clean"`
	NodeDownC             *NodeDownConf         `yaml:"node_down"`
}

// mount-pod clean
type AbnormalPodCleanConf struct {
	Enable                bool              `yaml:"enable"`
	CheckIntervalSeconds  int               `yaml:"check_interval_seconds"`
	DoubleCheckSecSeconds int               `yaml:"double_check_sec_seconds"`
	LabelSelector         string            `yaml:"label_selector"`
	FieldSelector         string            `yaml:"field_selector"`
	EnabledClusters       map[string]string `yaml:"enabled_clusters"`
	ImDingDingC           *ImDingDingConf   `yaml:"im_ding_ding"`
}

// node down
type NodeDownConf struct {
	Enable                  bool              `yaml:"enable"`
	CheckIntervalSeconds    int               `yaml:"check_interval_seconds"`
	QueryPromTimeOutSeconds int               `yaml:"query_prom_time_out_seconds"`
	CheckQls                []string          `yaml:"check_qls"`
	EnabledClusters         map[string]string `yaml:"enabled_clusters"`
	NodeNameToIpQl          string            `yaml:"node_name_to_ip_ql"` // name 转ip的 ql
	ImDingDingC             *ImDingDingConf   `yaml:"im_ding_ding"`
}

/*
ntp_offset:

	check_interval_seconds: 60   # 多久触发检查
	query_prom_time_out_seconds: 5 # 查询prometheus 超时
	cordon_daily_limit: 1      # 每日cordon 保护措施节点数量
	alert_title: "guard服务发现gpu坏卡 禁止节点并驱逐pod" # 通知的title
	check_ql: |-
	  gpu_bad_card_info >0
	enabled_clusters:  # 在哪些集群上开启这个模块
	  cpu-compute-01: yes
*/

type RecoveryConf struct {
	Enable                  bool            `yaml:"enable"`
	CheckIntervalSeconds    int             `yaml:"check_interval_seconds"`
	QueryPromTimeOutSeconds int             `yaml:"query_prom_time_out_seconds"`
	CheckDayNum             int             `yaml:"check_day_num"`
	ImDingDingC             *ImDingDingConf `yaml:"im_ding_ding"`
}

type ModuleCommonConf struct {
	Enable                  bool              `yaml:"enable"`
	CheckIntervalSeconds    int               `yaml:"check_interval_seconds"`
	QueryPromTimeOutSeconds int               `yaml:"query_prom_time_out_seconds"`
	DeletePodDryRun         bool              `yaml:"delete_pod_dry_run"`
	CheckQlMap              map[string]string `yaml:"check_ql_map"`
	RecoveryQlMap           map[string]string `yaml:"recovery_ql_map"`
	CordonDailyLimitMap     map[string]int    `yaml:"cordon_daily_limit_map"`
	EnabledClusters         map[string]string `yaml:"enabled_clusters"`
	ImDingDingC             *ImDingDingConf   `yaml:"im_ding_ding"`
}

type PluginNtpConf struct {
	Enable                  bool              `yaml:"enable"`
	CheckIntervalSeconds    int               `yaml:"check_interval_seconds"`
	QueryPromTimeOutSeconds int               `yaml:"query_prom_time_out_seconds"`
	DeletePodDryRun         bool              `yaml:"delete_pod_dry_run"`
	CheckQl                 string            `yaml:"check_ql"`
	RecoveryQl              string            `yaml:"recovery_ql"`
	CordonDailyLimit        int               `yaml:"cordon_daily_limit"`
	EnabledClusters         map[string]string `yaml:"enabled_clusters"`
	ImDingDingC             *ImDingDingConf   `yaml:"im_ding_ding"`
}

// 定义mysql配置结构体
type MySQLConf struct {
	Name  string `yaml:"name"`
	Addr  string `yaml:"addr"` // dsn
	Max   int    `yaml:"max"`
	Idle  int    `yaml:"idle"`
	Debug bool   `yaml:"debug"`
}

// 定义配置结构体
type ImDingDingConf struct {
	BotApiAddr string   `yaml:"bot_api_addr"`
	Title      string   `yaml:"title"`
	AtMobiles  []string `yaml:"atMobiles"`
}

// Load 把读取配置文件内容的string 用 yaml读成Config
func Load(s string) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadFile 根据配置文件的路径读取 掉Load 返回
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg, err := Load(string(content))
	if err != nil {
		fmt.Printf("[parsing YAML file errr...][error:%v]", err)
		return nil, err
	}
	return cfg, nil
}
