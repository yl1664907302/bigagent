package machine

import (
	"bigagent/scrape/machine/info"
	"bigagent/scrape/machine/memory"
)

// Machine 存放所有的采集层数据，被懒汉式创建
type Machine struct {
	I *info.Info     `json:"i"`
	M *memory.Memory `json:"m"`
}

func NewMachine() *Machine {
	return &Machine{I: info.NewInfo(), M: memory.NewMemory()}
}

func NotifyMachineAddressChange() {
	// 非阻塞写入 MachineCh 通道，通知监听者地址变化
	select {
	case MachineCh <- struct{}{}:
		// 通知成功
	default:
		// 通道已满，跳过通知（避免阻塞）
	}
}

var (
	Ma        *Machine
	MachineCh = make(chan struct{}, 1)
)
