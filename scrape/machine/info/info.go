package info

import (
	"github.com/super-l/machine-code/machine"
)

// 定义系统类型
type Info struct {
	Uuid string `json:"uuid"`
}

// GetUuid 获得系统uuid
func (i *Info) GetUuid() string {

	data := machine.GetMachineData()
	i.Uuid = data.PlatformUUID
	return i.Uuid
}

// 对外接口
func NewInfo() *Info {
	return &Info{}
}
