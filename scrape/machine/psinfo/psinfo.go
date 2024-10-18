package psinfo

import (
	"github.com/shirou/gopsutil/v4/process"
)

var (
	openFile   process.OpenFilesStat
	memoryInfo process.MemoryInfoStat
	//signalInfo     process.SignalInfoStat
	rlimit         process.RlimitStat
	ioCountersInfo process.IOCountersStat
	numCtxSwitches int
	//pageFaultsInfo process.PageFaultsStat
)

// 定义进程类
type PsInfo struct {
	InfoPorgress   []*process.Process `json:"info_porgress"`
	InfoOpenFiles  string             `json:"info_open_files"`
	InfoMemoryInfo string             `json:"info_memory_info"`
	//InfoSignal         string
	InfoRlimit         string `json:"info_rlimit"`
	InfoIOCounters     string `json:"info_iocounters"`
	InfoNumCtxSwitches int    `json:"info_num_ctx_switches"`
	InfoPageFaults     string `json:"info_page_faults"`
}

func NewProgressInfo() *PsInfo {
	processes, _ := process.Processes()
	return &PsInfo{
		InfoPorgress:   processes,
		InfoOpenFiles:  openFile.String(),
		InfoMemoryInfo: memoryInfo.String(),
		//InfoSignal: ,
		InfoRlimit:         rlimit.String(),
		InfoIOCounters:     ioCountersInfo.String(),
		InfoNumCtxSwitches: numCtxSwitches,
		//InfoPageFaults:     pageFaultsInfo,
	}
}

func (pgf *PsInfo) GetPorgressInfo() []*process.Process {
	return pgf.InfoPorgress
}
func (pgf *PsInfo) GetOpenFilesInfo() string {
	return pgf.InfoOpenFiles
}
func (pgf *PsInfo) GetMemoryInfo() string {
	return pgf.InfoMemoryInfo
}

//	func (pgf *PsInfo) GetSignalInfo() process.Signal {
//		return pgf.InfoSignal
//	}
func (pgf *PsInfo) GetRlimitInfo() string {
	return pgf.InfoRlimit
}
func (pgf *PsInfo) GetIOCountersInfo() string {
	return pgf.InfoIOCounters
}
func (pgf *PsInfo) GetNumCtxSwitchesInfo() int {
	return pgf.InfoNumCtxSwitches
}

//func (pgf *PsInfo) GetPageFaultsInfo() process.PageFaultsStat {
//	return pgf.InfoPageFaults
//}
