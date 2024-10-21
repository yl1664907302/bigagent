package model

type VeopsData struct {
	Ci_type                string `json:"ci_type,omitempty"`
	Cmdb_auto_product_uuid string `json:"cmdb_auto_product_uuid,omitempty"`
	Cmdb_auto_update_time  string `json:"cmdb_auto_update_time,omitempty"`
	Cmdb_auto_ipaddr       string `json:"cmdb_auto_ipaddr,omitempty"`
	Cmdb_auto_env          string `json:"cmdb_auto_env,omitempty"`
	Cmdb_auto_machine_type string `json:"cmdb_auto_machine_type,omitempty"`
	Cmdb_auto_os_type      string `json:"cmdb_auto_os_type,omitempty"`
	No_attribute_policy    string `json:"no_attribute_policy,omitempty"`
	Exist_policy           string `json:"exist_policy,omitempty"`
	Secret                 string `json:"_secret,omitempty"`
	Key                    string `json:"_key,omitempty"`
	Sort                   string `json:"sort,omitempty"`
	Page                   string `json:"page,omitempty"`
	Count                  string `json:"count,omitempty"`
	Ret_key                string `json:"ret_key,omitempty"`
}

func NewVeopsdata() *VeopsData {
	return &VeopsData{}
}
