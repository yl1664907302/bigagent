package diskinfo

import (
	grpc_server "bigagent/grpcs/server"
	"bigagent/scrape/machine/formatsize"
	"encoding/json"
	"log"

	"github.com/shirou/gopsutil/v4/disk"
)

type smpInfo struct {
	Path        string `json:"path"`
	Total       string `json:"total"`
	Free        string `json:"free"`
	Used        string `json:"used"`
	UsedPercent string `json:"usedpercent"`
	Device      string `json:"device"`
	Fstype      string `json:"fstype"`
	MountPoint  string `json:"mountpoint"`
}

type SmpDisk map[string]smpInfo

func NewSmpDisk() *SmpDisk {
	// 初始化smpdisk
	smpdisk := make(SmpDisk)

	// 获取disk信息
	diskInfo, err := disk.Partitions(true)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range diskInfo {
		usage, err := disk.Usage(i.Mountpoint)
		if err != nil {
			continue
		}

		// smpdisk[i.Device] = make(map[string]smpInfo)

		smpinfo := smpInfo{
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
	// 初始化smpdisk
	smpdisk := make(map[string]*grpc_server.SmpDisk)
	// 获取disk信息
	diskInfo, err := disk.Partitions(true)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range diskInfo {
		usage, err := disk.Usage(i.Mountpoint)
		if err != nil {
			continue
		}

		smpinfo := smpInfo{
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
		log.Fatal(err)
	}
	return string(b)
}
