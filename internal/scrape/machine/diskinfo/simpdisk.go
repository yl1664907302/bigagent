package diskinfo

import (
	"bigagent/internal/scrape/machine/formatsize"
	utils "bigagent/internal/util"
	"bigagent/internal/web/grpcs/server"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v4/disk"
)

type SmpInfo struct {
	Path        string `json:"path"`
	Total       string `json:"total"`
	Free        string `json:"free"`
	Used        string `json:"used"`
	UsedPercent string `json:"usedpercent"`
	Device      string `json:"device"`
	Fstype      string `json:"fstype"`
	MountPoint  string `json:"mountpoint"`
}

type SmpDisk map[string]SmpInfo

func NewSmpDisk() *SmpDisk {
	// 初始化smpdisk
	smpdisk := make(SmpDisk)

	// 获取disk信息
	diskInfo, err := disk.Partitions(true)
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	for _, i := range diskInfo {
		usage, err := disk.Usage(i.Mountpoint)
		if err != nil {
			continue
		}
		smpinfo := SmpInfo{
			Path:        usage.Path,
			Total:       formatsize.FormatSize(usage.Total),
			Free:        formatsize.FormatSize(usage.Free),
			Used:        formatsize.FormatSize(usage.Used),
			UsedPercent: formatsize.FormatPercent(usage.UsedPercent),
			Device:      i.Device,
			Fstype:      usage.Fstype,
			MountPoint:  i.Mountpoint,
		}
		smpdisk[i.Mountpoint] = smpinfo

	}

	return &smpdisk
}

func NewSmpDiskGrpc() *map[string]*grpc_server.SmpDisk {
	defer func() {
		if r := recover(); r != nil {
			utils.DefaultLogger.Error(fmt.Sprintf("捕获到 panic: %v", r))
		}
	}()
	// 初始化smpdisk
	smpdisk := make(map[string]*grpc_server.SmpDisk)
	// 获取disk信息
	diskInfo, err := disk.Partitions(true)
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	for _, i := range diskInfo {
		usage, err := disk.Usage(i.Mountpoint)
		if err != nil {
			continue
		}

		smpinfo := SmpInfo{
			Path:        usage.Path,
			Total:       formatsize.FormatSize(usage.Total),
			Free:        formatsize.FormatSize(usage.Free),
			Used:        formatsize.FormatSize(usage.Used),
			UsedPercent: formatsize.FormatPercent(usage.UsedPercent),
			Device:      i.Device,
			Fstype:      usage.Fstype,
			MountPoint:  i.Mountpoint,
		}
		smpdisk[i.Mountpoint] = &grpc_server.SmpDisk{
			Path:        smpinfo.Path,
			Total:       smpinfo.Total,
			Free:        smpinfo.Free,
			Used:        smpinfo.Used,
			UsedPercent: smpinfo.UsedPercent,
			Device:      smpinfo.Device,
			Fstype:      smpinfo.Fstype,
			MountPoint:  smpinfo.MountPoint,
		}

	}

	return &smpdisk
}

func (d *SmpDisk) ToString() string {
	b, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	return string(b)
}
