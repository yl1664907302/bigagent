package model

import (
	"bigagent/internal/config/global"
	"bigagent/internal/scrape/machine"
	"time"
)

type VeopsMachineData_GUOCHU struct {
	UUID                     string `json:"cmdb_auto_product_uuid"`
	OID                      string `json:"oid"`
	CiID                     int    `json:"_id"`
	CiType                   string `json:"ci_type"`
	OSType                   string `json:"cmdb_auto_os_type"`
	IP                       string `json:"cmdb_auto_ipaddr"`
	CmdbGeographicalPosition string `json:"cmdb_geographical_position"`
	CmdbIdracIpaddr          string `json:"cmdb_idrac_ipaddr"`
	CmdbAutoCpu              string `json:"cmdb_auto_cpu"`
	CmdbAutoMemory           string `json:"cmdb_auto_memory"`
	CmdbAutoDisk             string `json:"cmdb_auto_disk"`
	CmdbAutoMachineType      string `json:"cmdb_auto_machine_type"`
	CmdbAutoEnv              string `json:"cmdb_auto_env"`
	CmdbAutoUpdateTime       string `json:"cmdb_auto_update_time"`
	CmdbAutoOstype           string `json:"cmdb_auto_os_type"`
	CmdbAutoNodename         string `json:"cmdb_auto_nodename"`
	CmdbAutoArchitecture     string `json:"cmdb_auto_architecture"`
	CmdbAutoSystemRelease    string `json:"cmdb_auto_system_release"`
	CmdbAutoSystemVendor     string `json:"cmdb_auto_system_vendor"`
	CmdbAutoBiosDate         string `json:"cmdb_auto_bios_date"`
}

type VeoCoreStruct struct {
	UUID   string `json:"cmdb_auto_product_uuid"`
	OID    string `json:"oid"`
	OSType string `json:"cmdb_auto_os_type"`
	IP     string `json:"cmdb_auto_ipaddr"`
	CiID   int    `json:"_id"`
}

func NewVeoGuoChuData() *VeopsMachineData_GUOCHU {
	data := new(VeopsMachineData_GUOCHU)
	data.UUID = machine.SmpMa.Uuid
	data.IP = machine.SmpMa.IPv4
	data.CmdbAutoArchitecture = machine.SmpMa.Arch
	data.CmdbAutoCpu = machine.SmpMa.Cpu.ToString()
	data.CmdbAutoDisk = machine.SmpMa.Disk.ToString()
	data.CmdbAutoEnv = global.V.GetString("global.env")
	data.CmdbAutoMachineType = global.V.GetString("global.machineType")
	data.CmdbAutoMemory = machine.SmpMa.Memory.Vmem.Total
	data.CmdbAutoNodename = machine.SmpMa.Hostname
	data.CmdbAutoOstype = machine.SmpMa.Os
	data.CmdbAutoSystemRelease = machine.SmpMa.Kernel
	data.CmdbAutoSystemVendor = machine.SmpMa.Platform
	data.CmdbAutoBiosDate = machine.SmpMa.BiosData
	data.CmdbAutoUpdateTime = time.Now().Format("2006-01-02 15:04:05")
	return data
}
