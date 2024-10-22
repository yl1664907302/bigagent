package meminfo

import (
	"context"
	"encoding/json"
	"github.com/shirou/gopsutil/v4/mem"
	"log"
	"strings"
)

// 内存信息
type MemInfo struct {
	VirMem     *mem.VirtualMemoryStat
	SwapMem    *mem.SwapMemoryStat
	SwapDevice *mem.SwapDevice
}

// GetSwapDev 获取交换内存设备信息
func (m *MemInfo) GetSwapDev() ([]*mem.SwapDevice, error) {
	swapDevices, err := mem.SwapDevicesWithContext(context.Background())
	return swapDevices, err
}

// GetVirMem 获取虚拟内存信息
func (m *MemInfo) GetVirMem() (*mem.VirtualMemoryStat, error) {
	virtualMemoryStat, err := mem.VirtualMemoryWithContext(context.Background())
	return virtualMemoryStat, err
}

// GetSwapMem 获取交换内存信息
func (m *MemInfo) GetSwapMem() (*mem.SwapMemoryStat, error) {
	swapMemoryStat, err := mem.SwapMemoryWithContext(context.Background())
	return swapMemoryStat, err
}

func (m *MemInfo) ToString() (string, error) {
	var v []string
	virMem, err := m.GetVirMem()
	if err != nil {
		log.Println("get virmem err", err)
		return "", err
	}
	v = append(v, virMem.String())

	swapMem, err := m.GetSwapMem()
	if err != nil {
		log.Println("get swap mem err", err)
		return "", err
	}
	v = append(v, swapMem.String())

	devs, err := m.GetSwapDev()
	if err != nil {
		log.Println("get swap devices err", err)
		return "", err
	}

	for _, dev := range devs {
		v = append(v, dev.String())
	}
	mems := strings.Join(v, "")
	marshal, err := json.Marshal(mems)
	return string(marshal), err
}

// NewMemInfo 对外接口
func NewMemInfo() *MemInfo {
	return &MemInfo{}
}
